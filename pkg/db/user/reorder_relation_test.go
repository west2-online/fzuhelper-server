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
	"database/sql"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBUser_ReorderFriendList(t *testing.T) {
	type testCase struct {
		name           string
		stuId          string
		friendIds      []string
		mockTxError    error
		expectingError bool
		errorMsg       string
	}

	testCases := []testCase{
		{
			name:           "success_reorder_multiple_friends",
			stuId:          "102301001",
			friendIds:      []string{"102301002", "102301003", "102301004"},
			mockTxError:    nil,
			expectingError: false,
		},
		{
			name:           "success_reorder_single_friend",
			stuId:          "102301001",
			friendIds:      []string{"102301002"},
			mockTxError:    nil,
			expectingError: false,
		},
		{
			name:           "success_empty_friend_list",
			stuId:          "102301001",
			friendIds:      []string{},
			mockTxError:    nil,
			expectingError: false,
		},
		{
			name:           "transaction_error",
			stuId:          "102301001",
			friendIds:      []string{"102301002", "102301003"},
			mockTxError:    gorm.ErrInvalidDB,
			expectingError: true,
			errorMsg:       "dal.ReorderFriendList error",
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

			mockey.Mock((*gorm.DB).Transaction).To(func(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
				if tc.mockTxError != nil {
					return tc.mockTxError
				}
				return fc(mockGormDB)
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

			mockey.Mock((*gorm.DB).Order).To(func(value interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Find).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				// 模拟查询返回空列表（无需补集逻辑）
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Update).To(func(column string, value interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			err := mockDBUser.ReorderFriendList(context.Background(), tc.stuId, tc.friendIds)

			if tc.expectingError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
