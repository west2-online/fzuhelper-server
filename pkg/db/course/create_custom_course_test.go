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

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBCourse_CreateCustomCourse(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		input          *model.UserCustomCourse
		expectingError bool
	}

	inputCourse := &model.UserCustomCourse{
		StuId:      "222200311",
		Term:       "202401",
		CourseId:   "uuid-1",
		Name:       "自习",
		Location:   "图书馆",
		StartClass: 1,
		EndClass:   2,
		StartWeek:  1,
		EndWeek:    16,
		Weekday:    1,
	}

	testCases := []testCase{
		{
			name:           "CreateCustomCourse_Success",
			mockError:      nil,
			input:          inputCourse,
			expectingError: false,
		},
		{
			name:           "CreateCustomCourse_DBError",
			mockError:      fmt.Errorf("db error"),
			input:          inputCourse,
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
			mockey.Mock((*gorm.DB).Create).To(func(value interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				return mockGormDB
			}).Build()

			err := mockDBCourse.CreateCustomCourse(context.Background(), tc.input)

			if tc.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDBCourse_CheckDuplicateCustomCourse(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		mockCount      int64
		expectingError bool
		expectedResult bool
	}

	testCases := []testCase{
		{
			name:           "CheckDuplicateCustomCourse_Duplicate",
			mockError:      nil,
			mockCount:      1,
			expectingError: false,
			expectedResult: true,
		},
		{
			name:           "CheckDuplicateCustomCourse_NotDuplicate",
			mockError:      nil,
			mockCount:      0,
			expectingError: false,
			expectedResult: false,
		},
		{
			name:           "CheckDuplicateCustomCourse_DBError",
			mockError:      fmt.Errorf("db error"),
			mockCount:      0,
			expectingError: true,
			expectedResult: false,
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
			mockey.Mock((*gorm.DB).Model).To(func(value interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Count).To(func(count *int64) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				*count = tc.mockCount
				return mockGormDB
			}).Build()

			result, err := mockDBCourse.CheckDuplicateCustomCourse(
				context.Background(),
				"222200311", "202401",
				"自习", "图书馆",
				1, 2, 1, 16, 1,
			)

			if tc.expectingError {
				assert.Error(t, err)
				assert.False(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
