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

package user

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBUser_GetRelationByUserId(t *testing.T) {
	type testCase struct {
		name           string
		followerId     string
		followedId     string
		mockError      error
		mockRelation   *model.FollowRelation
		expectFound    bool
		expectRelation *model.FollowRelation
		expectingError bool
		errorMsg       string
	}

	now := time.Now()
	relationOK := &model.FollowRelation{
		Id:         1,
		FollowerId: "user1",
		FollowedId: "user2",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	testCases := []testCase{
		{
			name:           "success_found_relation",
			followerId:     "user1",
			followedId:     "user2",
			mockError:      nil,
			mockRelation:   relationOK,
			expectFound:    true,
			expectRelation: relationOK, // 期望返回实际的relation
			expectingError: false,
		},
		{
			name:           "record_not_found",
			followerId:     "user1",
			followedId:     "user3",
			mockError:      gorm.ErrRecordNotFound,
			mockRelation:   nil,
			expectFound:    false,
			expectRelation: nil,
			expectingError: false,
		},
		{
			name:           "database_error",
			followerId:     "user1",
			followedId:     "user2",
			mockError:      gorm.ErrInvalidDB,
			mockRelation:   nil,
			expectFound:    false,
			expectRelation: nil,
			expectingError: true,
			errorMsg:       "dal.GetRelationByUserId error",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			mockSnowflake := new(utils.Snowflake)
			mockDBUser := NewDBUser(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			whereCallCount := 0
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				whereCallCount++
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}

				if tc.mockRelation != nil {
					destValue := reflect.ValueOf(dest)
					if destValue.Kind() == reflect.Ptr {
						elem := destValue.Elem()
						relationValue := reflect.ValueOf(tc.mockRelation).Elem()
						elem.Set(relationValue)
					}
				}

				return mockGormDB
			}).Build()

			found, relation, err := mockDBUser.GetRelationByUserId(
				context.Background(),
				tc.followerId,
				tc.followedId,
			)

			if tc.expectingError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
				assert.False(t, found)
				assert.Nil(t, relation)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectFound, found)

				if tc.expectFound {
					assert.NotNil(t, relation)
					assert.Equal(t, tc.mockRelation.Id, relation.Id)
					assert.Equal(t, tc.mockRelation.FollowerId, relation.FollowerId)
					assert.Equal(t, tc.mockRelation.FollowedId, relation.FollowedId)
				} else {
					assert.Nil(t, relation)
				}
			}
		})
	}
}

func TestDBUser_GetUserFriendListLength(t *testing.T) {
	type testCase struct {
		name           string
		stuId          string
		mockError      error
		mockCount      int64
		expectLength   int64
		expectingError bool
		errorMsg       string
	}

	testCases := []testCase{
		{
			name:           "success_with_multiple_friends",
			stuId:          "102301001",
			mockError:      nil,
			mockCount:      5,
			expectLength:   5,
			expectingError: false,
		},
		{
			name:           "success_with_zero_friends",
			stuId:          "102301002",
			mockError:      nil,
			mockCount:      0,
			expectLength:   0,
			expectingError: false,
		},
		{
			name:           "success_with_one_friend",
			stuId:          "102301003",
			mockError:      nil,
			mockCount:      1,
			expectLength:   1,
			expectingError: false,
		},
		{
			name:           "database_connection_error",
			stuId:          "102301001",
			mockError:      gorm.ErrInvalidDB,
			mockCount:      0,
			expectLength:   -1,
			expectingError: true,
			errorMsg:       "dal.GetUserFriendListLength error",
		},
		{
			name:           "empty_user_id",
			stuId:          "",
			mockError:      nil,
			mockCount:      0,
			expectLength:   0,
			expectingError: false,
		},
		{
			name:           "record_not_found_error",
			stuId:          "102301001",
			mockError:      gorm.ErrRecordNotFound,
			mockCount:      0,
			expectLength:   -1,
			expectingError: true,
			errorMsg:       "dal.GetUserFriendListLength error",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := &gorm.DB{
				Config: &gorm.Config{},
				Error:  tc.mockError,
			}
			mockSnowflake := new(utils.Snowflake)
			mockDBUser := NewDBUser(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Count).To(func(count *int64) *gorm.DB {
				mockGormDB.Error = tc.mockError

				if tc.mockError == nil && count != nil {
					*count = tc.mockCount
				}

				return mockGormDB
			}).Build()

			length, err := mockDBUser.GetUserFriendListLength(context.Background(), tc.stuId)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Equal(t, tc.expectLength, length)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
				if tc.mockError != nil {
					assert.NotNil(t, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectLength, length)
			}
		})
	}
}

func TestDBUser_GetUserFriends(t *testing.T) {
	type testCase struct {
		name           string
		stuId          string
		mockError      error
		mockFriends    []*model.UserFriend
		expectFriends  []*model.UserFriend
		expectingError bool
		errorMsg       string
	}

	now := time.Now()

	// 测试数据 - 使用指针类型
	mockFriends1 := []*model.UserFriend{
		{
			FriendId:  "user2",
			UpdatedAt: now.Add(-24 * time.Hour),
		},
		{
			FriendId:  "user3",
			UpdatedAt: now.Add(-12 * time.Hour),
		},
	}

	mockFriends2 := []*model.UserFriend{
		{
			FriendId:  "user4",
			UpdatedAt: now.Add(-6 * time.Hour),
		},
	}

	mockFriendsEmpty := []*model.UserFriend{}

	testCases := []testCase{
		{
			name:           "success_with_multiple_friends",
			stuId:          "user1",
			mockError:      nil,
			mockFriends:    mockFriends1,
			expectFriends:  mockFriends1,
			expectingError: false,
		},
		{
			name:           "success_with_single_friend",
			stuId:          "user5",
			mockError:      nil,
			mockFriends:    mockFriends2,
			expectFriends:  mockFriends2,
			expectingError: false,
		},
		{
			name:           "success_with_no_friends",
			stuId:          "user6",
			mockError:      nil,
			mockFriends:    mockFriendsEmpty,
			expectFriends:  mockFriendsEmpty,
			expectingError: false,
		},
		{
			name:           "database_error",
			stuId:          "user1",
			mockError:      gorm.ErrInvalidDB,
			mockFriends:    nil,
			expectFriends:  nil,
			expectingError: true,
			errorMsg:       "dal.GetUserFriends error",
		},
		{
			name:           "empty_user_id",
			stuId:          "",
			mockError:      nil,
			mockFriends:    mockFriendsEmpty,
			expectFriends:  mockFriendsEmpty,
			expectingError: false,
		},
		{
			name:           "connection_error",
			stuId:          "user1",
			mockError:      gorm.ErrInvalidTransaction,
			mockFriends:    nil,
			expectFriends:  nil,
			expectingError: true,
			errorMsg:       "dal.GetUserFriends error",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := &gorm.DB{
				Config: &gorm.Config{},
				Error:  tc.mockError,
			}
			mockSnowflake := new(utils.Snowflake)
			mockDBUser := NewDBUser(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Select).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Find).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				mockGormDB.Error = tc.mockError

				if tc.mockError == nil && dest != nil && tc.mockFriends != nil {
					destValue := reflect.ValueOf(dest)
					if destValue.Kind() == reflect.Ptr {
						slicePtr := destValue.Elem()

						sliceType := reflect.TypeOf(tc.mockFriends)
						newSlice := reflect.MakeSlice(sliceType, len(tc.mockFriends), len(tc.mockFriends))

						for i := range tc.mockFriends {
							newSlice.Index(i).Set(reflect.ValueOf(tc.mockFriends[i]))
						}

						slicePtr.Set(newSlice)
					}
				}

				return mockGormDB
			}).Build()

			friends, err := mockDBUser.GetUserFriends(context.Background(), tc.stuId)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, friends)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)

				if len(tc.expectFriends) == 0 {
					assert.NotNil(t, friends)
					assert.Equal(t, 0, len(friends))
				} else {
					assert.NotNil(t, friends)
					assert.Equal(t, len(tc.expectFriends), len(friends))

					for i := range tc.expectFriends {
						assert.NotNil(t, friends[i])
						assert.Equal(t, tc.expectFriends[i].FriendId, friends[i].FriendId)
						assert.Equal(t, tc.expectFriends[i].UpdatedAt.Unix(), friends[i].UpdatedAt.Unix())
					}
				}
			}
		})
	}
}
