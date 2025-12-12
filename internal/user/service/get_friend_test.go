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
	"sort"
	"strings"
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

func TestUserService_GetFriendList(t *testing.T) {
	type testCase struct {
		name string

		expectingError    bool
		expectingErrorMsg string

		cacheFriendListExist bool
		cacheFriendListError error
		cacheFriendIds       []string

		dbFriendIds      []string
		dbFriendIdsError error

		cacheStuInfoExist bool
		cacheStuInfoError error
		cacheStuInfoMap   map[string]*dbmodel.Student

		dbStuInfoExist bool
		dbStuInfoError error
		dbStuInfoMap   map[string]*dbmodel.Student
	}

	stuId := "102300217"
	friendId1 := "102300218"
	friendId2 := "102300219"

	stuInfo1 := &dbmodel.Student{
		StuId:    friendId1,
		Name:     "张三",
		Sex:      "男",
		Birthday: "2000-01-01",
		College:  "计算机与大数据学院",
		Grade:    2023,
		Major:    "计算机类",
	}

	stuInfo2 := &dbmodel.Student{
		StuId:    friendId2,
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
			expectingError:       true,
			expectingErrorMsg:    "service.GetUserFriendCache:",
			cacheFriendListExist: true,
			cacheFriendListError: errno.InternalServiceError,
		},
		{
			name:                 "db friend ids error",
			expectingError:       true,
			expectingErrorMsg:    "service.GetUserFriendsIdDB:",
			cacheFriendListExist: false,
			dbFriendIdsError:     gorm.ErrInvalidData,
		},
		{
			name:                 "cache stu info get error",
			expectingError:       true,
			expectingErrorMsg:    "service.GetFriendList:",
			cacheFriendListExist: true,
			cacheFriendIds:       []string{friendId1},
			cacheStuInfoExist:    true,
			cacheStuInfoError:    errno.InternalServiceError,
		},
		{
			name:                 "db stu info error",
			expectingError:       true,
			expectingErrorMsg:    "service.GetFriendList:",
			cacheFriendListExist: true,
			cacheFriendIds:       []string{friendId1},
			cacheStuInfoExist:    false,
			dbStuInfoExist:       false,
			dbStuInfoError:       gorm.ErrInvalidData,
		},
		{
			name:                 "success with cache friend list and cache stu info",
			expectingError:       false,
			cacheFriendListExist: true,
			cacheFriendIds:       []string{friendId1, friendId2},
			cacheStuInfoExist:    true,
			cacheStuInfoMap: map[string]*dbmodel.Student{
				friendId1: stuInfo1,
				friendId2: stuInfo2,
			},
			dbStuInfoExist: true,
			dbStuInfoMap: map[string]*dbmodel.Student{
				friendId1: stuInfo1,
				friendId2: stuInfo2,
			},
		},
		{
			name:                 "success with db friend list and db stu info",
			expectingError:       false,
			cacheFriendListExist: false,
			dbFriendIds:          []string{friendId1},
			dbStuInfoExist:       true,
			dbStuInfoMap: map[string]*dbmodel.Student{
				friendId1: stuInfo1,
			},
		},
		{
			name:                 "success with cache set error",
			expectingError:       false,
			cacheFriendListExist: false,
			dbFriendIds:          []string{friendId1},
			dbStuInfoExist:       true,
			dbStuInfoMap: map[string]*dbmodel.Student{
				friendId1: stuInfo1,
			},
		},
		{
			name:                 "empty friend list",
			expectingError:       false,
			cacheFriendListExist: true,
			cacheFriendIds:       []string{},
		},
	}

	defer mockey.UnPatchAll()
	mockey.Mock((*user.CacheUser).SetUserFriendListCache).To(func(ctx context.Context, stuId string, friendIds []string) error {
		return nil
	}).Build()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    &db.Database{},
				CacheClient: &cache.Cache{},
			}
			mockClientSet.CacheClient.User = &user.CacheUser{}
			userService := NewUserService(context.Background(), "", nil, mockClientSet)

			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				if strings.Contains(key, "user_friends:") {
					return tc.cacheFriendListExist
				}
				return tc.cacheStuInfoExist
			}).Build()

			mockey.Mock((*user.CacheUser).GetUserFriendCache).To(func(ctx context.Context, key string) ([]string, error) {
				if tc.cacheFriendListError != nil {
					return nil, tc.cacheFriendListError
				}
				return tc.cacheFriendIds, nil
			}).Build()

			mockey.Mock((*userDB.DBUser).GetUserFriendsId).To(func(ctx context.Context, stuId string) ([]string, error) {
				if tc.dbFriendIdsError != nil {
					return nil, tc.dbFriendIdsError
				}
				return tc.dbFriendIds, nil
			}).Build()

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

			if tc.expectingError {
				assert.Error(t, err)
				if tc.expectingErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.expectingErrorMsg)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, friendList)

			expectedLength := 0
			if tc.cacheFriendListExist {
				expectedLength = len(tc.cacheFriendIds)
			} else {
				expectedLength = len(tc.dbFriendIds)
			}
			assert.Equal(t, expectedLength, len(friendList))

			if expectedLength > 0 {
				var expectedIds []string
				if tc.cacheFriendListExist {
					expectedIds = tc.cacheFriendIds
				} else {
					expectedIds = tc.dbFriendIds
				}

				actualIds := make([]string, len(friendList))
				for i, friend := range friendList {
					actualIds[i] = friend.StuId
				}

				sort.Strings(expectedIds)
				sort.Strings(actualIds)
				assert.Equal(t, expectedIds, actualIds, "friend list stuIds should match")
			}
		})
		time.Sleep(500 * time.Millisecond)
	}
}
