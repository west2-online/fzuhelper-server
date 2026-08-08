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

func TestDBCourse_GetCustomCourses(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		stuId          string
		term           string
		expectedResult []*model.UserCustomCourse
		expectingError bool
	}

	expectedCourses := []*model.UserCustomCourse{
		{
			ID:         1,
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
		},
		{
			ID:         2,
			StuId:      "222200311",
			Term:       "202401",
			CourseId:   "uuid-2",
			Name:       "开会",
			Location:   "会议室",
			StartClass: 3,
			EndClass:   4,
			StartWeek:  1,
			EndWeek:    16,
			Weekday:    3,
		},
	}

	testCases := []testCase{
		{
			name:           "GetCustomCourses_Success",
			mockError:      nil,
			stuId:          "222200311",
			term:           "202401",
			expectedResult: expectedCourses,
			expectingError: false,
		},
		{
			name:           "GetCustomCourses_Empty",
			mockError:      nil,
			stuId:          "222200311",
			term:           "202402",
			expectedResult: []*model.UserCustomCourse{},
			expectingError: false,
		},
		{
			name:           "GetCustomCourses_DBError",
			mockError:      fmt.Errorf("db error"),
			stuId:          "222200311",
			term:           "202401",
			expectedResult: nil,
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
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Find).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				courses, ok := dest.(*[]*model.UserCustomCourse)
				if ok && tc.expectedResult != nil {
					*courses = tc.expectedResult
				}
				return mockGormDB
			}).Build()

			result, err := mockDBCourse.GetCustomCourses(context.Background(), tc.stuId, tc.term)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestDBCourse_GetCustomCourseByID(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		stuId          string
		term           string
		courseId       string
		expectedResult *model.UserCustomCourse
		expectingError bool
	}

	expectedCourse := &model.UserCustomCourse{
		ID:         1,
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
			name:           "GetCustomCourseByID_Success",
			mockError:      nil,
			stuId:          "222200311",
			term:           "202401",
			courseId:       "uuid-1",
			expectedResult: expectedCourse,
			expectingError: false,
		},
		{
			name:           "GetCustomCourseByID_NotFound",
			mockError:      gorm.ErrRecordNotFound,
			stuId:          "222200311",
			term:           "202401",
			courseId:       "not-exist",
			expectedResult: nil,
			expectingError: true,
		},
		{
			name:           "GetCustomCourseByID_DBError",
			mockError:      fmt.Errorf("db error"),
			stuId:          "222200311",
			term:           "202401",
			courseId:       "uuid-1",
			expectedResult: nil,
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
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				course, ok := dest.(*model.UserCustomCourse)
				if ok && tc.expectedResult != nil {
					*course = *tc.expectedResult
				}
				return mockGormDB
			}).Build()

			result, err := mockDBCourse.GetCustomCourseByID(context.Background(), tc.stuId, tc.term, tc.courseId)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
