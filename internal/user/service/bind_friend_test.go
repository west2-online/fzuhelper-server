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
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/cache/user"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	userDB "github.com/west2-online/fzuhelper-server/pkg/db/user"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestBindInvitation(t *testing.T) {
	type testCase struct {
		name              string
		expectError       string
		cacheGetError     error
		cacheFriendId     string
		dbRelationExist   bool
		dbRelationError   error
		dbCreateError     error
		userConfined      bool
		targetConfined    bool
		userConfinedError error
		targetConfinedErr error
	}

	stuId := "102300217"
	friendId := "102300218"
	code := "ABCDEF"

	testCases := []testCase{
		{
			name:          "cache get error",
			expectError:   "service.GetCodeStuIdMappingCode:",
			cacheGetError: assert.AnError,
		},
		{
			name:          "add self as friend",
			expectError:   "无法添加自己为好友",
			cacheFriendId: stuId,
		},
		{
			name:            "relation already exist",
			expectError:     "好友关系已存在",
			cacheFriendId:   friendId,
			dbRelationExist: true,
			dbRelationError: nil,
		},
		{
			name:            "db relation check error",
			expectError:     "service.GetRelationByUserId:",
			cacheFriendId:   friendId,
			dbRelationExist: false,
			dbRelationError: gorm.ErrInvalidData,
		},
		{
			name:              "user confined check error",
			expectError:       "assert.AnError",
			cacheFriendId:     friendId,
			dbRelationExist:   false,
			dbRelationError:   nil,
			userConfinedError: assert.AnError,
		},
		{
			name:            "db create error",
			expectError:     "service.CreateRelation:",
			cacheFriendId:   friendId,
			dbRelationExist: false,
			dbRelationError: nil,
			dbCreateError:   gorm.ErrInvalidData,
		},
		{
			name:            "success",
			cacheFriendId:   friendId,
			dbRelationExist: false,
			dbRelationError: nil,
			dbCreateError:   nil,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// 由于 BindInvitation 内部会启动一个 goroutine 来更新缓存和删除邀请码相关的缓存，所以在测试中我们需要等待这个 goroutine 执行完毕才能正确断言结果。
			shouldWait := tc.expectError == ""
			var wg sync.WaitGroup
			if shouldWait {
				wg.Add(4)
			}

			isKeyExistGuard := mockey.Mock((*cache.Cache).IsKeyExist).Return(true).Build()
			defer isKeyExistGuard.UnPatch()

			setUserFriendGuard := mockey.Mock((*user.CacheUser).SetUserFriendCache).To(func(ctx context.Context, stuId string, friend *dbmodel.UserFriend) error {
				if shouldWait {
					wg.Done()
				}
				return nil
			}).Build()
			defer setUserFriendGuard.UnPatch()

			removeCodeGuard := mockey.Mock((*user.CacheUser).RemoveCodeStuIdMappingCache).To(func(ctx context.Context, key string) error {
				if shouldWait {
					wg.Done()
				}
				return nil
			}).Build()
			defer removeCodeGuard.UnPatch()

			removeInvitationGuard := mockey.Mock((*user.CacheUser).RemoveInvitationCodeCache).To(func(ctx context.Context, key string) error {
				if shouldWait {
					wg.Done()
				}
				return nil
			}).Build()
			defer removeInvitationGuard.UnPatch()

			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			mockey.Mock((*user.CacheUser).GetCodeStuIdMappingCache).Return(tc.cacheFriendId, tc.cacheGetError).Build()
			mockey.Mock((*userDB.DBUser).GetRelationByUserId).Return(tc.dbRelationExist, nil, tc.dbRelationError).Build()

			// Mock 好友数量检查
			mockey.Mock((*UserService).IsFriendNumsConfined).To(func(s *UserService, stuId string) (bool, error) {
				if stuId == "102300217" {
					return tc.userConfined, tc.userConfinedError
				}
				return tc.targetConfined, tc.targetConfinedErr
			}).Build()

			mockey.Mock((*UserService).writeRelationToDB).Return(tc.dbCreateError).Build()

			err := userService.BindInvitation(stuId, code)
			if shouldWait && err == nil {
				done := make(chan struct{})
				go func() {
					wg.Wait()
					close(done)
				}()
				select {
				case <-done:
				case <-time.After(500 * time.Millisecond):
					t.Fatalf("async cache update did not finish in time")
				}
			}

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestWriteRelationToDB(t *testing.T) {
	type testCase struct {
		name         string
		followedId   string
		followerId   string
		snowflakeId1 int64
		snowflakeId2 int64
		snowflakeErr error
		dbError      error
		expectError  bool
	}

	followedId := "102300217"
	followerId := "102300218"

	testCases := []testCase{
		{
			name:         "successful write to database",
			followedId:   followedId,
			followerId:   followerId,
			snowflakeId1: 1001,
			snowflakeId2: 1002,
			snowflakeErr: nil,
			dbError:      nil,
			expectError:  false,
		},
		{
			name:         "first snowflake ID generation fails",
			followedId:   followedId,
			followerId:   followerId,
			snowflakeErr: fmt.Errorf("snowflake generation error"),
			expectError:  true,
		},
		{
			name:         "second snowflake ID generation fails",
			followedId:   followedId,
			followerId:   followerId,
			snowflakeId1: 1001,
			snowflakeErr: fmt.Errorf("snowflake generation error"),
			expectError:  true,
		},
		{
			name:         "database write fails",
			followedId:   followedId,
			followerId:   followerId,
			snowflakeId1: 1001,
			snowflakeId2: 1002,
			snowflakeErr: nil,
			dbError:      fmt.Errorf("database write error"),
			expectError:  true,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient: new(utils.Snowflake),
				DBClient: new(db.Database),
			}
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			snowflakeCallCount := 0
			mockey.Mock((*utils.Snowflake).NextVal).To(func() (int64, error) {
				snowflakeCallCount++
				if snowflakeCallCount == 1 {
					return tc.snowflakeId1, tc.snowflakeErr
				}
				return tc.snowflakeId2, tc.snowflakeErr
			}).Build()

			mockey.Mock((*userDB.DBUser).CreateRelation).To(func(ctx context.Context, relations []*dbmodel.FollowRelation) error {
				if snowflakeCallCount != 2 {
					return fmt.Errorf("snowflake generator should be called exactly twice, called %d times", snowflakeCallCount)
				}

				if len(relations) != 2 {
					return fmt.Errorf("expected 2 relations, got %d", len(relations))
				}

				if relations[0].Id != tc.snowflakeId1 {
					return fmt.Errorf("first relation ID mismatch: expected %d, got %d", tc.snowflakeId1, relations[0].Id)
				}
				if relations[0].FollowedId != tc.followedId {
					return fmt.Errorf("first relation FollowedId mismatch: expected %s, got %s", tc.followedId, relations[0].FollowedId)
				}
				if relations[0].FollowerId != tc.followerId {
					return fmt.Errorf("first relation FollowerId mismatch: expected %s, got %s", tc.followerId, relations[0].FollowerId)
				}

				if relations[1].Id != tc.snowflakeId2 {
					return fmt.Errorf("second relation ID mismatch: expected %d, got %d", tc.snowflakeId2, relations[1].Id)
				}
				if relations[1].FollowedId != tc.followerId {
					return fmt.Errorf("second relation FollowedId mismatch: expected %s, got %s", tc.followerId, relations[1].FollowedId)
				}
				if relations[1].FollowerId != tc.followedId {
					return fmt.Errorf("second relation FollowerId mismatch: expected %s, got %s", tc.followedId, relations[1].FollowerId)
				}

				if relations[0].UpdatedAt.IsZero() {
					return fmt.Errorf("first relation UpdatedAt is zero")
				}
				if relations[1].UpdatedAt.IsZero() {
					return fmt.Errorf("second relation UpdatedAt is zero")
				}

				return tc.dbError
			}).Build()

			err := userService.writeRelationToDB(tc.followedId, tc.followerId)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, "unexpected error: %v", err)
			}
		})
	}
}

func TestIsFriendNumsConfined(t *testing.T) {
	type testCase struct {
		name           string
		cacheExist     bool
		cacheFriends   []*dbmodel.UserFriend
		cacheError     error
		dbLength       int64
		dbError        error
		expectConfined bool
		expectError    bool
	}

	// 初始化配置
	_ = config.InitForTest("friend")
	maxNum := int(config.Friend.MaxNum)
	stuId := "102300217"

	makeFriends := func(count int) []*dbmodel.UserFriend {
		friends := make([]*dbmodel.UserFriend, 0, count)
		for i := 0; i < count; i++ {
			friends = append(friends, &dbmodel.UserFriend{FriendId: fmt.Sprintf("%d", i)})
		}
		return friends
	}

	testCases := []testCase{
		{
			name:           "cache hit and confined",
			cacheExist:     true,
			cacheFriends:   makeFriends(maxNum),
			expectConfined: true,
		},
		{
			name:           "cache hit and not confined",
			cacheExist:     true,
			cacheFriends:   makeFriends(maxNum - 1),
			expectConfined: false,
		},
		{
			name:        "cache hit error",
			cacheExist:  true,
			cacheError:  assert.AnError,
			expectError: true,
		},
		{
			name:           "cache miss and confined",
			cacheExist:     false,
			dbLength:       int64(maxNum),
			expectConfined: true,
		},
		{
			name:           "cache miss and not confined",
			cacheExist:     false,
			dbLength:       int64(maxNum - 1),
			expectConfined: false,
		},
		{
			name:        "cache miss db error",
			cacheExist:  false,
			dbError:     assert.AnError,
			expectError: true,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()
			mockey.Mock((*user.CacheUser).GetUserFriendCache).Return(tc.cacheFriends, tc.cacheError).Build()
			mockey.Mock((*userDB.DBUser).GetUserFriendListLength).Return(tc.dbLength, tc.dbError).Build()

			userService := NewUserService(context.Background(), "", nil, mockClientSet)
			confined, err := userService.IsFriendNumsConfined(stuId)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectConfined, confined)
			}
		})
	}
}
