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

func TestDBCourse_DeleteCustomCourse(t *testing.T) {
	type testCase struct {
		name             string
		mockError        error
		mockRowsAffected int64
		stuId            string
		id               int64
		expectingError   bool
		expectedRows     int64
	}

	testCases := []testCase{
		{
			name:             "DeleteCustomCourse_Success",
			mockError:        nil,
			mockRowsAffected: 1,
			stuId:            "222200311",
			id:               1,
			expectingError:   false,
			expectedRows:     1,
		},
		{
			name:             "DeleteCustomCourse_NotFound",
			mockError:        nil,
			mockRowsAffected: 0,
			stuId:            "222200311",
			id:               999,
			expectingError:   false,
			expectedRows:     0,
		},
		{
			name:             "DeleteCustomCourse_DBError",
			mockError:        fmt.Errorf("db error"),
			mockRowsAffected: 0,
			stuId:            "222200311",
			id:               1,
			expectingError:   true,
			expectedRows:     0,
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
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Delete).To(func(value interface{}, conds ...interface{}) *gorm.DB {
				mockGormDB.RowsAffected = tc.mockRowsAffected
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				return mockGormDB
			}).Build()

			rows, err := mockDBCourse.DeleteCustomCourse(context.Background(), tc.stuId, tc.id)

			if tc.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedRows, rows)
		})
	}
}
