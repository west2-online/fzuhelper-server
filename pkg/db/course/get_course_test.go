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
	"errors"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBCourse_GetUserTermCourseByStuIdAndTerm(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		mockRows       *gorm.DB
		stuId          string
		term           string
		expectedResult *model.UserCourse
		expectingError bool
	}
	testCases := []testCase{
		{
			name:      "GetUserTermCourseByStuIdAndTerm_Success",
			mockError: nil,
			mockRows:  &gorm.DB{Error: nil},
			stuId:     "222200311",
			term:      "202401",
			expectedResult: &model.UserCourse{
				Id:                1001,
				StuId:             "222200311",
				Term:              "202401",
				TermCourses:       `[{"courseId":"C123","courseName":"Math"}]`,
				TermCoursesSha256: "abc123def456",
			},
			expectingError: false,
		},
		{
			name:           "GetUserTermCourseByStuIdAndTerm_RecordNotFound",
			mockError:      nil,
			mockRows:       &gorm.DB{Error: gorm.ErrRecordNotFound},
			stuId:          "222200311",
			term:           "202401",
			expectedResult: nil,
			expectingError: false,
		},
		{
			name:           "GetUserTermCourseByStuIdAndTerm_DBError",
			mockError:      fmt.Errorf("db error"),
			mockRows:       &gorm.DB{Error: fmt.Errorf("db error")},
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
			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				switch {
				case tc.mockError != nil:
					mockGormDB.Error = tc.mockError
					return mockGormDB
				case tc.mockRows != nil && errors.Is(tc.mockRows.Error, gorm.ErrRecordNotFound):
					mockGormDB.Error = gorm.ErrRecordNotFound
					return mockGormDB
				default:
					userCourse, ok := dest.(*model.UserCourse)
					if ok {
						*userCourse = *tc.expectedResult
					}
					return mockGormDB
				}
			}).Build()

			result, err := mockDBCourse.GetUserTermCourseByStuIdAndTerm(context.Background(), tc.stuId, tc.term)

			switch {
			case tc.mockRows != nil && errors.Is(tc.mockRows.Error, gorm.ErrRecordNotFound):
				assert.NoError(t, err)
				assert.Nil(t, result)
			case tc.expectingError:
				assert.Error(t, err)
				assert.Nil(t, result)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestDBCourse_GetUserTermCourseSha256ByStuIdAndTerm(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		mockRows       *gorm.DB
		stuId          string
		term           string
		expectedResult *model.UserCourse
		expectingError bool
	}

	testCases := []testCase{
		{
			name:      "GetUserTermCourseSha256ByStuIdAndTerm_Success",
			mockError: nil,
			mockRows:  &gorm.DB{Error: nil}, // No error
			stuId:     "222200311",
			term:      "202401",
			expectedResult: &model.UserCourse{
				Id:                1001,
				StuId:             "222200311",
				Term:              "202401",
				TermCourses:       `[{"courseId":"C123","courseName":"Math"}]`,
				TermCoursesSha256: "abc123def456",
			},
			expectingError: false,
		},
		{
			name:           "GetUserTermCourseSha256ByStuIdAndTerm_RecordNotFound",
			mockError:      nil,
			mockRows:       &gorm.DB{Error: gorm.ErrRecordNotFound}, // Simulate not found
			stuId:          "222200311",
			term:           "202401",
			expectedResult: nil,
			expectingError: false,
		},
		{
			name:           "GetUserTermCourseSha256ByStuIdAndTerm_DBError",
			mockError:      fmt.Errorf("db error"),
			mockRows:       &gorm.DB{Error: fmt.Errorf("db error")}, // Simulate db error
			stuId:          "222200311",
			term:           "202401",
			expectedResult: nil,
			expectingError: true,
		},
	}

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
			mockey.Mock((*gorm.DB).Select).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				switch {
				case tc.mockError != nil:
					mockGormDB.Error = tc.mockError
					return mockGormDB
				case tc.mockRows != nil && errors.Is(tc.mockRows.Error, gorm.ErrRecordNotFound):
					mockGormDB.Error = gorm.ErrRecordNotFound
					return mockGormDB
				default:
					userCourse, ok := dest.(*model.UserCourse)
					if ok {
						*userCourse = *tc.expectedResult
					}
					return mockGormDB
				}
			}).Build()

			result, err := mockDBCourse.GetUserTermCourseSha256ByStuIdAndTerm(context.Background(), tc.stuId, tc.term)

			switch {
			case tc.mockRows != nil && errors.Is(tc.mockRows.Error, gorm.ErrRecordNotFound):
				assert.NoError(t, err)
				assert.Nil(t, result)
			case tc.expectingError:
				assert.Error(t, err)
				assert.Nil(t, result)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
