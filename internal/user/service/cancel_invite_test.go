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

func TestCancelInvitationCode(t *testing.T) {
	type testCase struct {
		name           string
		expectError    string
		cacheExist     bool
		codeCacheError error
		mapKeyError    error
	}

	stuId := "102300217"
	testCode := "INVITE123"

	testCases := []testCase{
		{
			name:        "cache key does not exist",
			expectError: "当前账号暂无处于生效状态的邀请码",
			cacheExist:  false,
		},
		{
			name:       "success",
			cacheExist: true,
		},
		{
			name:        "get code cache error",
			expectError: "service.GetInvitationCodeCache",

			cacheExist:     true,
			codeCacheError: errors.New("code cache error"),
		},
		{
			name:        "get mapping cache error",
			expectError: "service.GetCodeStuIdMappingCodeCache",
			cacheExist:  true,
			mapKeyError: errors.New("map cache error"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			mockey.Mock(maincontext.ExtractIDFromLoginData).Return(stuId).Build()
			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()
			mockey.Mock((*user.CacheUser).GetInvitationCodeCache).Return(testCode, 0, tc.codeCacheError).Build()
			mockey.Mock((*user.CacheUser).GetCodeStuIdMappingCache).Return(stuId, tc.mapKeyError).Build()
			mockey.Mock((*user.CacheUser).RemoveCodeStuIdMappingCache).Return(nil).Build()
			mockey.Mock((*user.CacheUser).RemoveInvitationCodeCache).Return(nil).Build()

			loginData := &kitexModel.LoginData{}
			err := userService.CancelInvitationCode(loginData)
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
