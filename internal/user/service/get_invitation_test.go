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
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/cache/user"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestUserService_GetInvitationCode(t *testing.T) {
	type testCase struct {
		name string

		expectedExist     bool
		mockError         error
		expectingError    bool
		expectingErrorMsg string

		IsRefresh bool

		cacheExist    bool
		cacheGetError error
		cacheCode     string
	}
	stuId := "102300217"
	testCases := []testCase{
		{
			name:              "cache get error",
			IsRefresh:         false,
			expectedExist:     true,
			expectingError:    true,
			expectingErrorMsg: "service.GetInvitationCode:",
			mockError:         errno.InternalServiceError,
			cacheExist:        true,
			cacheGetError:     errno.InternalServiceError,
			cacheCode:         "",
		},
		{
			name:           "IsRefresh true - force regenerate",
			IsRefresh:      true,
			expectedExist:  true,
			expectingError: false,
			cacheExist:     true,
			cacheGetError:  nil,
			cacheCode:      "123456",
		},
		{
			name:           "cache code exist and no refresh",
			expectingError: false,
			expectedExist:  true,
			IsRefresh:      false,
			cacheExist:     true,
			cacheGetError:  nil,
			cacheCode:      "123456",
		},
		{
			name:           "cache not exist and refresh true",
			IsRefresh:      true,
			expectedExist:  false,
			expectingError: false,
			cacheExist:     false,
		},
		{
			name:           "cache not exist and refresh false",
			IsRefresh:      false,
			expectedExist:  false,
			expectingError: false,
			cacheExist:     false,
		},
	}
	defer mockey.UnPatchAll()
	mockey.Mock((*user.CacheUser).SetInvitationCodeCache).Return(nil).Build()
	mockey.Mock((*user.CacheUser).SetCodeStuIdMappingCache).Return(nil).Build()
	mockey.Mock((*user.CacheUser).RemoveCodeStuIdMappingCache).Return(nil).Build()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient:    &db.Database{},
				CacheClient: &cache.Cache{},
			}
			mockClientSet.CacheClient.User = &user.CacheUser{}
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.cacheExist
			}).Build()

			mockey.Mock((*user.CacheUser).GetInvitationCodeCache).To(func(ctx context.Context, key string) (code string, err error) {
				if tc.cacheGetError != nil {
					return "", tc.cacheGetError
				}
				return tc.cacheCode, nil
			}).Build()

			if !tc.cacheExist || tc.IsRefresh {
				mockey.Mock(utils.GenerateRandomCode).Return("ABCDEF").Build()
			}

			code, err := userService.GetInvitationCode(stuId, tc.IsRefresh)

			if tc.expectingError {
				assert.Equal(t, "", code)
				assert.Error(t, err)
				if tc.expectingErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.expectingErrorMsg)
				}
			} else {
				assert.NoError(t, err)
				if tc.cacheExist && !tc.IsRefresh && tc.cacheGetError == nil {
					assert.Equal(t, tc.cacheCode, code)
				}
				if !tc.cacheExist || tc.IsRefresh {
					assert.Equal(t, 6, len(code))
					assert.Equal(t, "ABCDEF", code)
				}
			}
		})
		time.Sleep(500 * time.Millisecond)
	}
}
