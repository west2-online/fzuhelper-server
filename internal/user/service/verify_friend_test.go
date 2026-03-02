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

func TestVerifyUserFriend(t *testing.T) {
	type testCase struct {
		name           string
		stuId          string
		friendId       string
		cacheKeyExist  bool
		cacheIsFriend  bool
		cacheError     error
		dbRelation     bool
		dbRelationData *dbmodel.FollowRelation
		dbError        error
		expectResult   bool
		expectError    string
	}

	stuId := "102300217"
	friendId := "102300218"

	testCases := []testCase{
		{
			name:          "cache exists and is friend",
			stuId:         stuId,
			friendId:      friendId,
			cacheKeyExist: true,
			cacheIsFriend: true,
			cacheError:    nil,
			expectResult:  true,
		},
		{
			name:          "cache exists but not friend",
			stuId:         stuId,
			friendId:      friendId,
			cacheKeyExist: true,
			cacheIsFriend: false,
			cacheError:    nil,
			expectResult:  false,
		},
		{
			name:          "cache exists but cache error",
			stuId:         stuId,
			friendId:      friendId,
			cacheKeyExist: true,
			cacheIsFriend: false,
			cacheError:    errno.InternalServiceError,
			expectResult:  false,
			expectError:   "service.VerifyUserFriend: Get friend cache fail:",
		},
		{
			name:           "cache not exist and db relation exists",
			stuId:          stuId,
			friendId:       friendId,
			cacheKeyExist:  false,
			dbRelation:     true,
			dbRelationData: &dbmodel.FollowRelation{},
			dbError:        nil,
			expectResult:   true,
		},
		{
			name:           "cache not exist and db relation not exists",
			stuId:          stuId,
			friendId:       friendId,
			cacheKeyExist:  false,
			dbRelation:     false,
			dbRelationData: nil,
			dbError:        nil,
			expectResult:   false,
		},
		{
			name:           "cache not exist and db error",
			stuId:          stuId,
			friendId:       friendId,
			cacheKeyExist:  false,
			dbRelation:     false,
			dbRelationData: nil,
			dbError:        gorm.ErrInvalidData,
			expectResult:   false,
			expectError:    "service.VerifyUserFriend: Get friend db fail:",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			isKeyExistGuard := mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheKeyExist).Build()
			defer isKeyExistGuard.UnPatch()

			isFriendCacheGuard := mockey.Mock((*user.CacheUser).IsFriendCache).Return(tc.cacheIsFriend, tc.cacheError).Build()
			defer isFriendCacheGuard.UnPatch()

			getRelationGuard := mockey.Mock((*userDB.DBUser).GetRelationByUserId).Return(tc.dbRelation, tc.dbRelationData, tc.dbError).Build()
			defer getRelationGuard.UnPatch()

			result, err := userService.VerifyUserFriend(tc.stuId, tc.friendId)
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectResult, result)
		})
	}
}
