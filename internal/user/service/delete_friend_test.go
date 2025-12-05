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
	"gorm.io/gorm"

	loginmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	maincontext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/cache/user"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	userDB "github.com/west2-online/fzuhelper-server/pkg/db/user"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestUserService_DeleteUserFriend(t *testing.T) {
	type testCase struct {
		name string

		expectingError    bool
		expectingErrorMsg string

		dbRelationExist  bool
		dbRelationError  error
		dbDeleteError    error
		cacheDeleteError error
	}

	stuId := "102300217"
	targetStuId := "102300218"

	testCases := []testCase{
		{
			name:              "relation not exist",
			expectingError:    true,
			expectingErrorMsg: "service.DeleteUserFriend: RelationShip No Exist",
			dbRelationExist:   false,
			dbRelationError:   nil,
		},
		{
			name:              "db relation check error",
			expectingError:    true,
			expectingErrorMsg: "service.GetRelationByUserId:",
			dbRelationExist:   false,
			dbRelationError:   gorm.ErrInvalidData,
		},
		{
			name:              "db delete error",
			expectingError:    true,
			expectingErrorMsg: "service.DeleteRelation:",
			dbRelationExist:   true,
			dbRelationError:   nil,
			dbDeleteError:     gorm.ErrInvalidData,
		},
		{
			name:             "success",
			expectingError:   false,
			dbRelationExist:  true,
			dbRelationError:  nil,
			dbDeleteError:    nil,
			cacheDeleteError: nil,
		},
	}

	defer mockey.UnPatchAll()
	mockey.Mock((*user.CacheUser).DeleteUserFriendCache).To(func(ctx context.Context, stuId, targetStuId string) error {
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

			// Mock context.ExtractIDFromLoginData
			mockey.Mock(maincontext.ExtractIDFromLoginData).Return(stuId).Build()

			// Mock GetRelationByUserId
			mockey.Mock((*userDB.DBUser).GetRelationByUserId).To(func(ctx context.Context, stuId, targetStuId string) (bool, *dbmodel.FollowRelation, error) {
				return tc.dbRelationExist, nil, tc.dbRelationError
			}).Build()

			// Mock DeleteRelation
			mockey.Mock((*userDB.DBUser).DeleteRelation).To(func(ctx context.Context, stuId, targetStuId string) error {
				return tc.dbDeleteError
			}).Build()

			loginData := &loginmodel.LoginData{}
			err := userService.DeleteUserFriend(loginData, targetStuId)

			if tc.expectingError {
				assert.Error(t, err)
				if tc.expectingErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.expectingErrorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
		time.Sleep(500 * time.Millisecond)
	}
}
