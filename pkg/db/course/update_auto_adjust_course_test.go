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

func TestDBCourse_UpdateAutoAdjustCourse(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		adjustCourse   *model.AutoAdjustCourse
		expectingError bool
	}

	testCases := []testCase{
		{
			name:           "UpdateAutoAdjustCourse_Success",
			mockError:      nil,
			adjustCourse:   &model.AutoAdjustCourse{Id: 1001, Enabled: true},
			expectingError: false,
		},
		{
			name:           "UpdateAutoAdjustCourse_DBError",
			mockError:      fmt.Errorf("db error"),
			adjustCourse:   &model.AutoAdjustCourse{Id: 1001, Enabled: true},
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
			mockey.Mock((*gorm.DB).Updates).To(func(values interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				return &gorm.DB{Error: nil}
			}).Build()

			err := mockDBCourse.UpdateAutoAdjustCourse(context.Background(), tc.adjustCourse)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "dal.UpdateAutoAdjustCourse update error")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
