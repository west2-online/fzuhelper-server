/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/admin"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	admincache "github.com/west2-online/fzuhelper-server/pkg/cache/admin"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	dbadmin "github.com/west2-online/fzuhelper-server/pkg/db/admin"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestCallback(t *testing.T) {
	require.NoError(t, config.InitForTest(constants.AdminServiceName))
	originalAppSecret := config.Admin.Feishu.AppSecret
	config.Admin.Feishu.AppSecret = "test-secret"
	t.Cleanup(func() {
		config.Admin.Feishu.AppSecret = originalAppSecret
	})

	type testCase struct {
		name              string
		stateExists       bool
		stateError        error
		tokenResponse     string
		tokenStatus       int
		userInfoResponse  string
		allowed           bool
		dbError           error
		ticketError       error
		expectErrorCode   int64
		expectRedirectKey string
		expectHttpCalls   int
		expectTicket      bool
	}

	testCases := []testCase{
		{
			name:              "success",
			stateExists:       true,
			tokenResponse:     `{"access_token":"user-token","token_type":"Bearer"}`,
			userInfoResponse:  `{"code":0,"data":{"user_id":"fzu-admin","name":"Admin"}}`,
			allowed:           true,
			expectRedirectKey: "ticket",
			expectHttpCalls:   2,
			expectTicket:      true,
		},
		{
			name:            "invalid state",
			expectErrorCode: errno.AuthErrorCode,
		},
		{
			name:              "user is not in whitelist",
			stateExists:       true,
			tokenResponse:     `{"access_token":"user-token","token_type":"Bearer"}`,
			userInfoResponse:  `{"code":0,"data":{"user_id":"unknown-user"}}`,
			expectRedirectKey: "error",
			expectHttpCalls:   2,
		},
		{
			name:            "Feishu rejects authorization code",
			stateExists:     true,
			tokenResponse:   `{"error":"invalid_grant"}`,
			tokenStatus:     http.StatusBadRequest,
			expectErrorCode: errno.AuthErrorCode,
			expectHttpCalls: 1,
		},
		{
			name:             "ticket cache error",
			stateExists:      true,
			tokenResponse:    `{"access_token":"user-token","token_type":"Bearer"}`,
			userInfoResponse: `{"code":0,"data":{"user_id":"fzu-admin"}}`,
			allowed:          true,
			ticketError:      errors.New("redis unavailable"),
			expectErrorCode:  errno.InternalRedisErrorCode,
			expectHttpCalls:  2,
			expectTicket:     true,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			const returnTo = "http://localhost:5173/auth/callback?from=sso"
			mockey.Mock((*admincache.CacheAdmin).ConsumeOAuthStateCache).
				Return(returnTo, tc.stateExists, tc.stateError).Build()

			var cachedTicket, cachedAdminId string
			ticketCalled := false
			mockey.Mock((*admincache.CacheAdmin).SetLoginTicketCache).
				To(func(_ *admincache.CacheAdmin, _ context.Context, ticket, adminId string) error {
					ticketCalled = true
					cachedTicket = ticket
					cachedAdminId = adminId
					return tc.ticketError
				}).Build()

			mockey.Mock((*dbadmin.DBAdmin).GetAdminByFeishuUserId).
				Return(tc.allowed, &model.AdminUser{Id: 42}, tc.dbError).Build()

			httpCalls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				httpCalls++
				w.Header().Set("Content-Type", "application/json")
				switch req.URL.Path {
				case "/token":
					if err := req.ParseForm(); err != nil {
						t.Errorf("parse token request form: %v", err)
						return
					}
					assert.Equal(t, "authorization_code", req.Form.Get("grant_type"))
					assert.Equal(t, "test-code", req.Form.Get("code"))
					assert.Equal(t, config.Admin.Feishu.AppID, req.Form.Get("client_id"))
					assert.Equal(t, "test-secret", req.Form.Get("client_secret"))
					assert.Equal(t, config.Admin.Feishu.RedirectURI, req.Form.Get("redirect_uri"))
					if tc.tokenStatus != 0 {
						w.WriteHeader(tc.tokenStatus)
					}
					_, _ = w.Write([]byte(tc.tokenResponse))
				case "/user_info":
					assert.Equal(t, "Bearer user-token", req.Header.Get("Authorization"))
					_, _ = w.Write([]byte(tc.userInfoResponse))
				default:
					t.Errorf("unexpected Feishu request URL: %s", req.URL.String())
				}
			}))
			defer server.Close()

			clientSet := &base.ClientSet{
				CacheClient: &cache.Cache{Admin: admincache.NewCacheAdmin(nil)},
				DBClient:    &db.Database{Admin: dbadmin.NewDBAdmin(nil)},
			}
			svc := NewAdminService(context.Background(), clientSet)
			svc.oauthConfig.Endpoint = oauth2.Endpoint{
				AuthURL:   server.URL + "/authorize",
				TokenURL:  server.URL + "/token",
				AuthStyle: oauth2.AuthStyleInParams,
			}
			svc.userInfoUrl = server.URL + "/user_info"
			redirectUrl, callbackErr := svc.Callback(&admin.CallbackRequest{State: "test-state", Code: "test-code"})

			assert.Equal(t, tc.expectHttpCalls, httpCalls)
			assert.Equal(t, tc.expectTicket, ticketCalled)
			if tc.expectErrorCode != 0 {
				require.Error(t, callbackErr)
				var gotErr errno.ErrNo
				require.ErrorAs(t, callbackErr, &gotErr)
				assert.Equal(t, tc.expectErrorCode, gotErr.ErrorCode)
				assert.Empty(t, redirectUrl)
				return
			}

			require.NoError(t, callbackErr)
			parsedUrl, err := url.Parse(redirectUrl)
			require.NoError(t, err)
			assert.Equal(t, "sso", parsedUrl.Query().Get("from"))
			if tc.expectRedirectKey == "ticket" {
				assert.Equal(t, cachedTicket, parsedUrl.Query().Get("ticket"))
				assert.Equal(t, "42", cachedAdminId)
				decodedTicket, err := base64.RawURLEncoding.DecodeString(cachedTicket)
				require.NoError(t, err)
				assert.Len(t, decodedTicket, oauthStateBytes)
			} else {
				assert.Equal(t, "forbidden", parsedUrl.Query().Get("error"))
			}
		})
	}
}
