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
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBUser_GetStudentById(t *testing.T) {
	type testCase struct {
		name            string
		inPutId         string
		mockError       error
		expectingError  bool
		expectedStudent *model.Student
		ErrorMsg        string
	}
	stu := &model.Student{
		StuId:   "102301000",
		Sex:     "男",
		College: "计算机与大数据学院",
		Grade:   2023,
		Major:   "计算机科学与技术",
	}
	testCases := []testCase{
		{
			name:            "success",
			inPutId:         stu.StuId,
			mockError:       nil,
			expectingError:  false,
			expectedStudent: stu,
		},
		{
			name:            "record not found",
			inPutId:         stu.StuId,
			mockError:       gorm.ErrRecordNotFound,
			expectingError:  true,
			expectedStudent: nil,
		},
		{
			name:            "error",
			inPutId:         stu.StuId,
			mockError:       gorm.ErrInvalidValue,
			expectingError:  true,
			expectedStudent: nil,
			ErrorMsg:        "dal.GetStudentById error",
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
			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
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

				if res, ok := dest.(*model.Student); ok && tc.expectedStudent != nil {
					*res = *tc.expectedStudent
				}
				return mockGormDB
			}).Build()

			_, result, err := mockDBUser.GetStudentById(context.Background(), tc.inPutId)
			if tc.expectingError {
				if err == nil {
					return
				}
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tc.ErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedStudent, result)
			}
		})
	}
}
