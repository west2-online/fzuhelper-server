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

func TestUserService_GetLoginDataForYJSY(t *testing.T) {
	type testCase struct {
		name              string
		expectedCookie    []*http.Cookie
		mockLoginError    error
		mockCookieError   error
		expectingError    bool
		expectingErrorMsg string
	}
	testCases := []testCase{
		{
			name: "success",
			expectedCookie: []*http.Cookie{
				{
					Name:  "YJSY_COOKIE",
					Value: "test_cookie_value",
				},
			},
			mockLoginError:  nil,
			mockCookieError: nil,
			expectingError:  false,
		},
		{
			name:              "yjsy login error",
			mockLoginError:    errno.InternalServiceError,
			expectingError:    true,
			expectingErrorMsg: errno.InternalServiceError.ErrorMsg,
		},
		{
			name:              "yjsy get cookies error",
			mockLoginError:    nil,
			mockCookieError:   errno.InternalServiceError,
			expectingError:    true,
			expectingErrorMsg: errno.InternalServiceError.ErrorMsg,
		},
	}
	req := &user.GetLoginDataForYJSYRequest{
		Id:       "102301000",
		Password: "102301000",
	}
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := new(base.ClientSet)
			mockClientSet.SFClient = new(utils.Snowflake)
			mockClientSet.DBClient = new(db.Database)
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			// Mock YJSY Student methods
			mockey.Mock((*yjsy.Student).WithUser).To(func(id string, password string) *yjsy.Student {
				return yjsy.NewStudent()
			}).Build()

			mockey.Mock((*yjsy.Student).Login).To(func() error {
				return tc.mockLoginError
			}).Build()

			mockey.Mock((*yjsy.Student).GetCookies).To(func() ([]*http.Cookie, error) {
				if tc.mockCookieError != nil {
					return nil, tc.mockCookieError
				}
				return tc.expectedCookie, nil
			}).Build()

			cookieStr, err := userService.GetLoginDataForYJSY(req)
			if tc.expectingError {
				assert.Equal(t, "", cookieStr)
				assert.Error(t, err)
				if tc.expectingErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.expectingErrorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, utils.ParseCookiesToString(tc.expectedCookie), cookieStr)
			}
		})
	}
}
