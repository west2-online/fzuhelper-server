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

func TestDBCourse_CreateAutoAdjustCourse(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		input          *model.AutoAdjustCourse
		expectedResult *model.AutoAdjustCourse
		expectingError bool
	}

	toDate := "2025-10-08"
	expectedResult := &model.AutoAdjustCourse{
		Id:       1001,
		Term:     "202501",
		FromDate: "2025-10-01",
		ToDate:   &toDate,
	}

	testCases := []testCase{
		{
			name:           "CreateAutoAdjustCourse_Success",
			mockError:      nil,
			input:          expectedResult,
			expectedResult: expectedResult,
			expectingError: false,
		},
		{
			name:           "CreateAutoAdjustCourse_DBError",
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
				autoAdjustCourse, ok := value.(*model.AutoAdjustCourse)
				if ok {
					*autoAdjustCourse = *tc.input
				}
				return mockGormDB
			}).Build()

			result, err := mockDBCourse.CreateAutoAdjustCourse(context.Background(), tc.input)

			if tc.expectingError {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "dal.CreateAutoAdjustCourse error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
