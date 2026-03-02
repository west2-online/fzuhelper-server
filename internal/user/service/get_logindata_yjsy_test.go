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
	"net/http"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/yjsy"
)

func TestGetLoginDataForYJSY(t *testing.T) {
	type testCase struct {
		name            string
		expectCookie    []*http.Cookie
		mockLoginError  error
		mockCookieError error
		expectError     string
	}

	testCases := []testCase{
		{
			name: "success",
			expectCookie: []*http.Cookie{
				{
					Name:  "YJSY_COOKIE",
					Value: "test_cookie_value",
				},
			},
			mockLoginError:  nil,
			mockCookieError: nil,
		},
		{
			name:           "yjsy login error",
			mockLoginError: errno.InternalServiceError,
			expectError:    errno.InternalServiceError.ErrorMsg,
		},
		{
			name:            "yjsy get cookies error",
			mockLoginError:  nil,
			mockCookieError: errno.InternalServiceError,
			expectError:     errno.InternalServiceError.ErrorMsg,
		},
	}

	req := &user.GetLoginDataForYJSYRequest{
		Id:       "102301000",
		Password: "102301000",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient: new(utils.Snowflake),
				DBClient: new(db.Database),
			}
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			// Mock YJSY Student methods
			mockey.Mock((*yjsy.Student).WithUser).Return(yjsy.NewStudent()).Build()

			mockey.Mock((*yjsy.Student).Login).Return(tc.mockLoginError).Build()

			mockey.Mock((*yjsy.Student).GetCookies).Return(tc.expectCookie, tc.mockCookieError).Build()

			cookieStr, err := userService.GetLoginDataForYJSY(req)
			if tc.expectError != "" {
				assert.Equal(t, "", cookieStr)
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, utils.ParseCookiesToString(tc.expectCookie), cookieStr)
			}
		})
	}
}
