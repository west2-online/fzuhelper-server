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
	"strings"
	"sync"
	"testing"
	"time"

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

func TestGetFriendList(t *testing.T) {
	type testCase struct {
		name                 string
		expectError          string
		cacheFriendListExist bool
		cacheFriendListError error
		cacheFriendIds       []*dbmodel.UserFriend
		dbFriendIds          []*dbmodel.UserFriend
		dbFriendIdsError     error
		cacheStuInfoExist    bool
		cacheStuInfoError    error
		cacheStuInfoMap      map[string]*dbmodel.Student
		dbStuInfoExist       bool
		dbStuInfoError       error
		dbStuInfoMap         map[string]*dbmodel.Student
	}

	stuId := "102300217"
	friendId1 := &dbmodel.UserFriend{FriendId: "102300218"}
	friendId2 := &dbmodel.UserFriend{FriendId: "102300219"}

	stuInfo1 := &dbmodel.Student{
		StuId:    friendId1.FriendId,
		Name:     "张三",
		Sex:      "男",
		Birthday: "2000-01-01",
		College:  "计算机与大数据学院",
		Grade:    2023,
		Major:    "计算机类",
	}

	stuInfo2 := &dbmodel.Student{
		StuId:    friendId2.FriendId,
		Name:     "李四",
		Sex:      "女",
		Birthday: "2000-02-02",
		College:  "电气学院",
		Grade:    2023,
		Major:    "电气工程",
	}

	testCases := []testCase{
		{
			name:                 "cache friend list get error",
			expectError:          "service.GetUserFriendCache:",
			cacheFriendListExist: true,
			cacheFriendListError: errno.InternalServiceError,
		},
		{
			name:                 "db friend ids error",
			expectError:          "service.GetUserFriendsIdDB:",
			cacheFriendListExist: false,
			dbFriendIdsError:     gorm.ErrInvalidData,
		},
		{
			name:                 "cache stu info get error",
			expectError:          "service.GetFriendList:",
			cacheFriendListExist: true,
			cacheFriendIds:       []*dbmodel.UserFriend{friendId1},
			cacheStuInfoExist:    true,
			cacheStuInfoError:    errno.InternalServiceError,
		},
		{
			name:                 "db stu info error",
			expectError:          "service.GetFriendList:",
			cacheFriendListExist: true,
			cacheFriendIds:       []*dbmodel.UserFriend{friendId1},
			cacheStuInfoExist:    false,
			dbStuInfoExist:       false,
			dbStuInfoError:       gorm.ErrInvalidData,
		},
		{
			name:                 "success with cache friend list and cache stu info",
			cacheFriendListExist: true,
			cacheFriendIds:       []*dbmodel.UserFriend{friendId1, friendId2},
			cacheStuInfoExist:    true,
			cacheStuInfoMap: map[string]*dbmodel.Student{
				friendId1.FriendId: stuInfo1,
				friendId2.FriendId: stuInfo2,
			},
			dbStuInfoExist: true,
			dbStuInfoMap: map[string]*dbmodel.Student{
				friendId1.FriendId: stuInfo1,
				friendId2.FriendId: stuInfo2,
			},
		},
		{
			name:                 "success with db friend list and db stu info",
			cacheFriendListExist: false,
			dbFriendIds:          []*dbmodel.UserFriend{friendId1},
			dbStuInfoExist:       true,
			dbStuInfoMap: map[string]*dbmodel.Student{
				friendId1.FriendId: stuInfo1,
			},
		},
		{
			name:                 "success with cache set error",
			cacheFriendListExist: false,
			dbFriendIds:          []*dbmodel.UserFriend{friendId1},
			dbStuInfoExist:       true,
			dbStuInfoMap: map[string]*dbmodel.Student{
				friendId1.FriendId: stuInfo1,
			},
		},
		{
			name:                 "empty friend list",
			cacheFriendListExist: true,
			cacheFriendIds:       nil,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			shouldWait := !tc.cacheFriendListExist && tc.dbFriendIdsError == nil
			var wg sync.WaitGroup
			if shouldWait {
				wg.Add(1)
			}

			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			setFriendListGuard := mockey.Mock((*user.CacheUser).SetUserFriendListCache).To(
				func(ctx context.Context, stuId string, friendList []*dbmodel.UserFriend) error {
					if shouldWait {
						wg.Done()
					}
					return nil
				}).Build()
			defer setFriendListGuard.UnPatch()

			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				if strings.Contains(key, "user_friends:") {
					return tc.cacheFriendListExist
				}
				return tc.cacheStuInfoExist
			}).Build()

			mockey.Mock((*user.CacheUser).GetUserFriendCache).Return(tc.cacheFriendIds, tc.cacheFriendListError).Build()

			mockey.Mock((*userDB.DBUser).GetUserFriends).Return(tc.dbFriendIds, tc.dbFriendIdsError).Build()

			mockey.Mock((*user.CacheUser).GetStuInfoCache).To(func(ctx context.Context, key string) (*dbmodel.Student, error) {
				if tc.cacheStuInfoError != nil {
					return nil, tc.cacheStuInfoError
				}
				if tc.cacheStuInfoMap != nil {
					if stuInfo, exists := tc.cacheStuInfoMap[key]; exists {
						return stuInfo, nil
					}
				}
				return nil, nil
			}).Build()

			mockey.Mock((*userDB.DBUser).GetStudentById).To(func(ctx context.Context, stuId string) (bool, *dbmodel.Student, error) {
				if tc.dbStuInfoError != nil {
					return false, nil, tc.dbStuInfoError
				}
				if tc.dbStuInfoMap != nil {
					if stuInfo, exists := tc.dbStuInfoMap[stuId]; exists {
						return tc.dbStuInfoExist, stuInfo, nil
					}
				}
				return tc.dbStuInfoExist, nil, nil
			}).Build()

			friendList, err := userService.GetFriendList(stuId)
			if shouldWait && err == nil {
				done := make(chan struct{})
				go func() {
					wg.Wait()
					close(done)
				}()
				select {
				case <-done:
				case <-time.After(500 * time.Millisecond):
					t.Fatalf("async cache set did not finish in time")
				}
			}
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, friendList)

				expectLength := 0
				if tc.cacheFriendListExist {
					expectLength = len(tc.cacheFriendIds)
				} else {
					expectLength = len(tc.dbFriendIds)
				}
				assert.Equal(t, expectLength, len(friendList))

				if expectLength > 0 {
					var expectedIds []*dbmodel.UserFriend
					if tc.cacheFriendListExist {
						expectedIds = tc.cacheFriendIds
					} else {
						expectedIds = tc.dbFriendIds
					}

					actualIds := make([]*dbmodel.UserFriend, len(friendList))
					for i, friend := range friendList {
						actualIds[i] = &dbmodel.UserFriend{FriendId: friend.StuId}
					}

					assert.Equal(t, expectedIds, actualIds, "friend list stuIds should match")
				}
			}
		})
	}
}
