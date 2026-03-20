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
	userDB "github.com/west2-online/fzuhelper-server/pkg/db/user"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestReorderFriendList(t *testing.T) {
	type testCase struct {
		name          string
		stuId         string
		friendIds     []string
		dbError       error
		cacheExist    bool
		cacheDelError error
		expectError   string
	}

	testCases := []testCase{
		{
			name:       "success_with_cache",
			stuId:      "102301001",
			friendIds:  []string{"102301002", "102301003"},
			dbError:    nil,
			cacheExist: true,
		},
		{
			name:       "success_without_cache",
			stuId:      "102301001",
			friendIds:  []string{"102301002"},
			dbError:    nil,
			cacheExist: false,
		},
		{
			name:        "db_error",
			stuId:       "102301001",
			friendIds:   []string{"102301002"},
			dbError:     gorm.ErrInvalidDB,
			expectError: "service.ReorderFriendList:",
		},
		{
			name:          "success_cache_delete_error_ignored",
			stuId:         "102301001",
			friendIds:     []string{"102301002"},
			dbError:       nil,
			cacheExist:    true,
			cacheDelError: gorm.ErrInvalidDB,
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
			userService := NewUserService(context.Background(), "", nil, mockClientSet, new(taskqueue.BaseTaskQueue))

			mockey.Mock((*userDB.DBUser).ReorderFriendList).Return(tc.dbError).Build()

			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()

			mockey.Mock((*user.CacheUser).InvalidateFriendListCache).Return(tc.cacheDelError).Build()

			err := userService.ReorderFriendList(tc.stuId, tc.friendIds)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
