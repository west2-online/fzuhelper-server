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
	"net/url"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/admin"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	admincache "github.com/west2-online/fzuhelper-server/pkg/cache/admin"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestLogin(t *testing.T) {
	require.NoError(t, config.InitForTest(constants.AdminServiceName))

	type testCase struct {
		name          string
		returnTo      string
		cacheError    error
		expectError   bool
		expectCache   bool
		expectErrCode int64
	}

	allowedReturnTo := config.Admin.Feishu.AllowedReturnUrls[0]
	testCases := []testCase{
		{
			name:        "success",
			returnTo:    allowedReturnTo,
			expectCache: true,
		},
		{
			name:          "return URL is not allowed",
			returnTo:      "https://evil.example/auth/callback",
			expectError:   true,
			expectErrCode: errno.ParamErrorCode,
		},
		{
			name:          "cache error",
			returnTo:      allowedReturnTo,
			cacheError:    errors.New("redis unavailable"),
			expectError:   true,
			expectCache:   true,
			expectErrCode: errno.InternalRedisErrorCode,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			var cachedState, cachedReturnTo string
			cacheCalled := false
			mockey.Mock((*admincache.CacheAdmin).SetOAuthStateCache).
				To(func(_ *admincache.CacheAdmin, _ context.Context, state, returnTo string) error {
					cacheCalled = true
					cachedState = state
					cachedReturnTo = returnTo
					return tc.cacheError
				}).Build()

			clientSet := &base.ClientSet{
				CacheClient: &cache.Cache{Admin: admincache.NewCacheAdmin(nil)},
			}
			svc := NewAdminService(context.Background(), clientSet)
			authorizeUrl, err := svc.Login(&admin.LoginRequest{ReturnTo: tc.returnTo})

			assert.Equal(t, tc.expectCache, cacheCalled)
			if tc.expectError {
				require.Error(t, err)
				assert.Empty(t, authorizeUrl)
				var gotErr errno.ErrNo
				require.ErrorAs(t, err, &gotErr)
				assert.Equal(t, tc.expectErrCode, gotErr.ErrorCode)
				return
			}

			require.NoError(t, err)
			parsedUrl, err := url.Parse(authorizeUrl)
			require.NoError(t, err)
			assert.Equal(t, constants.FeishuOAuthAuthorizeUrl, parsedUrl.Scheme+"://"+parsedUrl.Host+parsedUrl.Path)
			assert.Equal(t, config.Admin.Feishu.AppID, parsedUrl.Query().Get("client_id"))
			assert.Equal(t, config.Admin.Feishu.RedirectURI, parsedUrl.Query().Get("redirect_uri"))
			assert.Equal(t, constants.FeishuAdminOAuthScope, parsedUrl.Query().Get("scope"))
			assert.Equal(t, cachedState, parsedUrl.Query().Get("state"))
			assert.Equal(t, tc.returnTo, cachedReturnTo)

			decodedState, err := base64.RawURLEncoding.DecodeString(cachedState)
			require.NoError(t, err)
			assert.Len(t, decodedState, oauthStateBytes)
		})
	}
}
