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
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/api/mw"
	"github.com/west2-online/fzuhelper-server/kitex_gen/admin"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	admincache "github.com/west2-online/fzuhelper-server/pkg/cache/admin"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestExchangeTicket(t *testing.T) {
	type testCase struct {
		name          string
		adminId       string
		ticketExists  bool
		cacheError    error
		token         string
		tokenError    error
		expectError   bool
		expectErrCode int64
		expectSign    bool
	}

	testCases := []testCase{
		{
			name:         "success",
			adminId:      "42",
			ticketExists: true,
			token:        "admin-access-token",
			expectSign:   true,
		},
		{
			name:          "ticket does not exist",
			expectError:   true,
			expectErrCode: errno.AuthErrorCode,
		},
		{
			name:          "empty admin id",
			ticketExists:  true,
			expectError:   true,
			expectErrCode: errno.AuthErrorCode,
		},
		{
			name:          "cache error",
			cacheError:    errors.New("redis unavailable"),
			expectError:   true,
			expectErrCode: errno.InternalRedisErrorCode,
		},
		{
			name:          "sign token error",
			adminId:       "42",
			ticketExists:  true,
			tokenError:    errno.AuthError,
			expectError:   true,
			expectErrCode: errno.AuthErrorCode,
			expectSign:    true,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			var actualTicket, actualAdminId string
			signCalled := false
			mockey.Mock((*admincache.CacheAdmin).ConsumeLoginTicketCache).
				To(func(_ *admincache.CacheAdmin, _ context.Context, ticket string) (string, bool, error) {
					actualTicket = ticket
					return tc.adminId, tc.ticketExists, tc.cacheError
				}).Build()
			mockey.Mock(mw.CreateAdminToken).To(func(adminId string) (string, error) {
				signCalled = true
				actualAdminId = adminId
				return tc.token, tc.tokenError
			}).Build()

			clientSet := &base.ClientSet{
				CacheClient: &cache.Cache{Admin: admincache.NewCacheAdmin(nil)},
			}
			svc := NewAdminService(context.Background(), clientSet)
			accessToken, err := svc.ExchangeTicket(&admin.ExchangeTicketRequest{Ticket: "login-ticket"})

			assert.Equal(t, "login-ticket", actualTicket)
			assert.Equal(t, tc.expectSign, signCalled)
			if tc.expectSign {
				assert.Equal(t, tc.adminId, actualAdminId)
			}
			if tc.expectError {
				require.Error(t, err)
				assert.Empty(t, accessToken)
				var gotErr errno.ErrNo
				require.ErrorAs(t, err, &gotErr)
				assert.Equal(t, tc.expectErrCode, gotErr.ErrorCode)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.token, accessToken)
		})
	}
}
