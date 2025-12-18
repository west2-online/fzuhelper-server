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

	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	maincontext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/cache/user"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestUserService_CancelInvitationCode(t *testing.T) {
	type testCase struct {
		name           string
		expectingError bool
		errorContains  string
		cacheExist     bool
		codeCacheError error
		mapKeyError    error
	}

	stuId := "102300217"
	testCode := "INVITE123"

	testCases := []testCase{
		{
			name:           "cache key does not exist",
			expectingError: true,
			errorContains:  "当前账号没有邀请码",
			cacheExist:     false,
		},
		{
			name:           "success",
			expectingError: false,
			cacheExist:     true,
		},
		{
			name:           "get code cache error",
			expectingError: true,
			errorContains:  "service.GetInvitationCodeCache",
			cacheExist:     true,
			codeCacheError: errors.New("code cache error"),
		},
		{
			name:           "get mapping cache error",
			expectingError: true,
			errorContains:  "service.GetCodeStuIdMappingCodeCache",
			cacheExist:     true,
			mapKeyError:    errors.New("map cache error"),
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    &db.Database{},
				CacheClient: &cache.Cache{},
			}
			mockClientSet.CacheClient.User = &user.CacheUser{}
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			mockey.Mock(maincontext.ExtractIDFromLoginData).Return(stuId).Build()

			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.cacheExist
			}).Build()

			mockey.Mock((*user.CacheUser).GetInvitationCodeCache).To(func(ctx context.Context, codeKey string) (string, int64, error) {
				if tc.codeCacheError != nil {
					return "", 0, tc.codeCacheError
				}
				return testCode, 102300200, nil
			}).Build()

			mockey.Mock((*user.CacheUser).GetCodeStuIdMappingCache).To(func(ctx context.Context, mapKey string) (string, error) {
				if tc.mapKeyError != nil {
					return "", tc.mapKeyError
				}
				return stuId, nil
			}).Build()

			mockey.Mock((*user.CacheUser).RemoveCodeStuIdMappingCache).Return(nil).Build()

			mockey.Mock((*user.CacheUser).RemoveInvitationCodeCache).Return(nil).Build()

			loginData := &kitexModel.LoginData{}
			err := userService.CancelInvitationCode(loginData)

			if tc.expectingError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
