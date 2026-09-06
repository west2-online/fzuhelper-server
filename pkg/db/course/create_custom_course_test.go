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
		Id:         1,
		StuId:      "222200311",
		Term:       "202401",
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

func TestDBCourse_GetCustomCourseIDByContent(t *testing.T) {
	const existingCourseId int64 = 123

	type testCase struct {
		name           string
		mockFirstError error
		expectingError bool
		expectedID     int64
	}

	testCases := []testCase{
		{
			name:       "GetCustomCourseIDByContent_Found",
			expectedID: existingCourseId,
		},
		{
			name:           "GetCustomCourseIDByContent_NotFound",
			mockFirstError: gorm.ErrRecordNotFound,
		},
		{
			name:           "GetCustomCourseIDByContent_DBError",
			mockFirstError: fmt.Errorf("db error"),
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
			mockey.Mock((*gorm.DB).Model).To(func(value interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Select).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockFirstError != nil {
					mockGormDB.Error = tc.mockFirstError
					return mockGormDB
				}
				if v, ok := dest.(*model.UserCustomCourse); ok {
					v.Id = existingCourseId
				}
				return mockGormDB
			}).Build()

			result, err := mockDBCourse.GetCustomCourseIDByContent(
				context.Background(),
				"222200311", "202401",
				"自习", "张老师", "图书馆",
				1, 2, 1, 16, 1,
				false, false,
			)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Zero(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, result)
			}
		})
	}
}
