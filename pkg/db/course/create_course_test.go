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

func TestDBCourse_CreateUserTermCourse(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		input          *model.UserCourse
		expectedResult *model.UserCourse
		expectingError bool
	}

	expectedResult := &model.UserCourse{
		Id:                1001,
		StuId:             "222200311",
		Term:              "202401",
		TermCourses:       `[{"courseId":"C123","courseName":"Math"}]`,
		TermCoursesSha256: "abc123def456",
	}

	testCases := []testCase{
		{
			name:           "CreateUserTermCourse_Success",
			mockError:      nil,
			input:          expectedResult,
			expectedResult: expectedResult,
			expectingError: false,
		},
		{
			name:           "CreateUserTermCourse_DBError",
			mockError:      fmt.Errorf("db error"),
			input:          expectedResult,
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
			mockey.Mock((*gorm.DB).Create).To(func(value interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				userCourse, ok := value.(*model.UserCourse)
				if ok {
					*userCourse = *tc.input
				}
				return mockGormDB
			}).Build()

			result, err := mockDBCourse.CreateUserTermCourse(context.Background(), tc.input)

			if tc.expectingError {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "dal.CreateUserTermCourse error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
