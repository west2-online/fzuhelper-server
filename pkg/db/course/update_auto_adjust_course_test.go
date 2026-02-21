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

package course

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBCourse_UpdateAutoAdjustCourseEnabledByID(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		id             int64
		enabled        bool
		expectingError bool
	}

	testCases := []testCase{
		{
			name:           "UpdateAutoAdjustCourseEnabledByID_Success",
			mockError:      nil,
			id:             1001,
			enabled:        true,
			expectingError: false,
		},
		{
			name:           "UpdateAutoAdjustCourseEnabledByID_DBError",
			mockError:      fmt.Errorf("db error"),
			id:             1001,
			enabled:        true,
			expectingError: true,
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			mockSnowflake := new(utils.Snowflake)
			mockDBCourse := NewDBCourse(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Update).To(func(column string, value interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				return &gorm.DB{Error: nil}
			}).Build()

			err := mockDBCourse.UpdateAutoAdjustCourseEnabledByID(context.Background(), tc.id, tc.enabled)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "dal.UpdateAutoAdjustCourseEnabledByID error")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
