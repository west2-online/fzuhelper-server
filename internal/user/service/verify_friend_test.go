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

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/cache/user"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	userDB "github.com/west2-online/fzuhelper-server/pkg/db/user"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestUserService_VerifyUserFriend(t *testing.T) {
	type testCase struct {
		name string

		stuId    string
		friendId string

		// Mock 返回值
		cacheKeyExist  bool
		cacheIsFriend  bool
		cacheError     error
		dbRelation     bool
		dbRelationData *dbmodel.FollowRelation
		dbError        error

		// 期望结果
		expectingResult bool
		expectingError  bool
		errorMsg        string
	}

	stuId := "102300217"
	friendId := "102300218"

	testCases := []testCase{
		{
			name:            "cache exists and is friend",
			stuId:           stuId,
			friendId:        friendId,
			cacheKeyExist:   true,
			cacheIsFriend:   true,
			cacheError:      nil,
			expectingResult: true,
			expectingError:  false,
		},
		{
			name:            "cache exists but not friend",
			stuId:           stuId,
			friendId:        friendId,
			cacheKeyExist:   true,
			cacheIsFriend:   false,
			cacheError:      nil,
			expectingResult: false,
			expectingError:  false,
		},
		{
			name:            "cache exists but cache error",
			stuId:           stuId,
			friendId:        friendId,
			cacheKeyExist:   true,
			cacheIsFriend:   false,
			cacheError:      errno.InternalServiceError,
			expectingResult: false,
			expectingError:  true,
			errorMsg:        "service.VerifyUserFriend: Get friend cache fail:",
		},
		{
			name:            "cache not exist and db relation exists",
			stuId:           stuId,
			friendId:        friendId,
			cacheKeyExist:   false,
			dbRelation:      true,
			dbRelationData:  &dbmodel.FollowRelation{},
			dbError:         nil,
			expectingResult: true,
			expectingError:  false,
		},
		{
			name:            "cache not exist and db relation not exists",
			stuId:           stuId,
			friendId:        friendId,
			cacheKeyExist:   false,
			dbRelation:      false,
			dbRelationData:  nil,
			dbError:         nil,
			expectingResult: false,
			expectingError:  false,
		},
		{
			name:            "cache not exist and db error",
			stuId:           stuId,
			friendId:        friendId,
			cacheKeyExist:   false,
			dbRelation:      false,
			dbRelationData:  nil,
			dbError:         gorm.ErrInvalidData,
			expectingResult: false,
			expectingError:  true,
			errorMsg:        "service.VerifyUserFriend: Get friend db fail:",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClientSet := &base.ClientSet{
				SFClient:    &utils.Snowflake{},
				DBClient:    &db.Database{},
				CacheClient: &cache.Cache{},
			}
			mockClientSet.CacheClient.User = &user.CacheUser{}

			ctx := context.Background()
			userService := NewUserService(ctx, "", nil, mockClientSet)

			isKeyExistGuard := mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.cacheKeyExist
			}).Build()
			defer isKeyExistGuard.UnPatch()

			isFriendCacheGuard := mockey.Mock((*user.CacheUser).IsFriendCache).To(func(ctx context.Context, stuId, friendId string) (bool, error) {
				return tc.cacheIsFriend, tc.cacheError
			}).Build()
			defer isFriendCacheGuard.UnPatch()

			getRelationGuard := mockey.Mock((*userDB.DBUser).GetRelationByUserId).To(func(ctx context.Context, stuId, targetStuId string) (
				bool, *dbmodel.FollowRelation, error,
			) {
				return tc.dbRelation, tc.dbRelationData, tc.dbError
			}).Build()
			defer getRelationGuard.UnPatch()

			result, err := userService.VerifyUserFriend(tc.stuId, tc.friendId)

			if tc.expectingError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectingResult, result)
		})
	}
}
