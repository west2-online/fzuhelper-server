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

package admin

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestDBAdmin_GetAdminByFeishuUserId(t *testing.T) {
	type testCase struct {
		name          string
		userId        string
		adminUser     *model.AdminUser
		dbError       error
		expectAllowed bool
		expectError   bool
		expectErrCode int64
	}

	testCases := []testCase{
		{
			name:   "success",
			userId: "feishu-user-id",
			adminUser: &model.AdminUser{
				Id:           42,
				FeishuUserId: "feishu-user-id",
				Name:         "admin",
				Enabled:      true,
			},
			expectAllowed: true,
		},
		{
			name:    "record not found",
			userId:  "unknown-user-id",
			dbError: gorm.ErrRecordNotFound,
		},
		{
			name:          "database error",
			userId:        "feishu-user-id",
			dbError:       gorm.ErrInvalidDB,
			expectError:   true,
			expectErrCode: errno.InternalDatabaseErrorCode,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			dbAdmin := NewDBAdmin(mockGormDB)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				assert.NotNil(t, ctx)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				assert.Equal(t, constants.AdminUserTableName, name)
				assert.Empty(t, args)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				assert.Equal(t, "feishu_user_id = ? AND enabled = ?", query)
				assert.Equal(t, []interface{}{tc.userId, true}, args)
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				assert.Empty(t, conds)
				mockGormDB.Error = tc.dbError
				if tc.dbError == nil && tc.adminUser != nil {
					adminUser, ok := dest.(*model.AdminUser)
					require.True(t, ok)
					*adminUser = *tc.adminUser
				}
				return mockGormDB
			}).Build()

			allowed, adminUser, err := dbAdmin.GetAdminByFeishuUserId(context.Background(), tc.userId)

			assert.Equal(t, tc.expectAllowed, allowed)
			if tc.expectError {
				require.Error(t, err)
				assert.Nil(t, adminUser)
				var gotErr errno.ErrNo
				require.ErrorAs(t, err, &gotErr)
				assert.Equal(t, tc.expectErrCode, gotErr.ErrorCode)
				assert.ErrorContains(t, err, "dal.GetAdminByFeishuUserId error")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.adminUser, adminUser)
		})
	}
}
