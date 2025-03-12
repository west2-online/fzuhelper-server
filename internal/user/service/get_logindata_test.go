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
	"github.com/west2-online/jwch"
)

func TestUserService_GetLoginData(t *testing.T) {
	type testCase struct {
		name           string
		expectedId     string
		expectedCookie []*http.Cookie
		mockError      error
		expectingError bool
	}
	testCases := []testCase{
		{
			name:       "success",
			expectedId: "2024102301000",
			expectedCookie: []*http.Cookie{
				{
					Name: "test",
				},
			},
		},
		{
			name:           "jwch error",
			mockError:      errno.InternalServiceError,
			expectingError: true,
		},
	}
	req := &user.GetLoginDataRequest{
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
			mockey.Mock((*jwch.Student).GetIdentifierAndCookies).To(func() (string, []*http.Cookie, error) {
				if tc.expectingError {
					return "", nil, tc.mockError
				}
				return tc.expectedId, tc.expectedCookie, nil
			}).Build()

			id, cookie, err := userService.GetLoginData(req)
			if tc.expectingError {
				assert.Equal(t, cookie, "")
				assert.Contains(t, err.Error(), errno.InternalServiceError.ErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedId, id)
				assert.Equal(t, utils.ParseCookiesToString(tc.expectedCookie), cookie)
			}
		})
	}
}
