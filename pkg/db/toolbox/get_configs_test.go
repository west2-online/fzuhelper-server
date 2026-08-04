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

package toolbox

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBToolbox_ListToolboxConfigs(t *testing.T) {
	type whereExpectation struct {
		query string
		arg   interface{}
	}

	type testCase struct {
		name           string
		pageNum        int
		pageSize       int
		filter         ListToolboxConfigsFilter
		mockTotal      int64
		mockConfigs    []*model.ToolboxConfig
		mockCountError error
		mockFindError  error
		expectError    string
	}

	testCases := []testCase{
		{
			name:      "success_with_configs",
			pageNum:   2,
			pageSize:  2,
			mockTotal: 3,
			mockConfigs: []*model.ToolboxConfig{
				{Id: 3, ToolID: 2, Name: new("tool 2")},
				{Id: 2, ToolID: 1, Name: new("tool 1 android")},
			},
		},
		{
			name:     "success_with_filters",
			pageNum:  1,
			pageSize: 20,
			filter: ListToolboxConfigsFilter{
				ToolID:     new(int64(1)),
				StudentID:  new("102300217"),
				Platform:   new("android"),
				MinVersion: new(int64(2)),
			},
			mockTotal: 1,
			mockConfigs: []*model.ToolboxConfig{
				{Id: 3, ToolID: 1, StudentID: new("102300217"), Platform: new("android"), Version: new(int64(2))},
			},
		},
		{
			name:        "success_with_empty_configs",
			pageNum:     1,
			pageSize:    20,
			mockTotal:   0,
			mockConfigs: []*model.ToolboxConfig{},
		},
		{
			name:           "count_error",
			pageNum:        1,
			pageSize:       20,
			mockCountError: gorm.ErrInvalidDB,
			expectError:    "dal.ListToolboxConfigs count error",
		},
		{
			name:          "find_error",
			pageNum:       1,
			pageSize:      20,
			mockTotal:     3,
			mockFindError: gorm.ErrInvalidValue,
			expectError:   "dal.ListToolboxConfigs error",
		},
		{
			name:        "invalid_page_num",
			pageNum:     0,
			pageSize:    20,
			expectError: "page_num and page_size must be positive",
		},
		{
			name:        "invalid_page_size",
			pageNum:     1,
			pageSize:    0,
			expectError: "page_num and page_size must be positive",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			whereExpectations := make([]whereExpectation, 0, 4)
			if tc.filter.ToolID != nil {
				whereExpectations = append(whereExpectations, whereExpectation{"tool_id = ?", *tc.filter.ToolID})
			}
			if tc.filter.StudentID != nil && *tc.filter.StudentID != "" {
				whereExpectations = append(whereExpectations, whereExpectation{"student_id = ?", *tc.filter.StudentID})
			}
			if tc.filter.Platform != nil && *tc.filter.Platform != "" {
				whereExpectations = append(whereExpectations, whereExpectation{"platform = ?", *tc.filter.Platform})
			}
			if tc.filter.MinVersion != nil {
				whereExpectations = append(whereExpectations, whereExpectation{"version >= ?", *tc.filter.MinVersion})
			}
			whereCalls := 0

			mockGormDB := new(gorm.DB)
			mockSnowflake := new(utils.Snowflake)
			mockDBToolbox := NewDBToolbox(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Model).To(func(value interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				assert.NotEmpty(t, whereExpectations)
				expected := whereExpectations[whereCalls%len(whereExpectations)]
				assert.Equal(t, expected.query, query)
				assert.Equal(t, []interface{}{expected.arg}, args)
				whereCalls++
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Count).To(func(count *int64) *gorm.DB {
				if tc.mockCountError != nil {
					mockGormDB.Error = tc.mockCountError
					return mockGormDB
				}
				*count = tc.mockTotal
				mockGormDB.Error = nil
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Order).To(func(value interface{}) *gorm.DB {
				assert.Equal(t, "id DESC", value)
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Limit).To(func(limit int) *gorm.DB {
				assert.Equal(t, tc.pageSize, limit)
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Offset).To(func(offset int) *gorm.DB {
				assert.Equal(t, (tc.pageNum-1)*tc.pageSize, offset)
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Find).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockFindError != nil {
					mockGormDB.Error = tc.mockFindError
					return mockGormDB
				}

				configs, ok := dest.(*[]*model.ToolboxConfig)
				if ok {
					*configs = tc.mockConfigs
				}
				mockGormDB.Error = nil
				return mockGormDB
			}).Build()

			configs, total, err := mockDBToolbox.ListToolboxConfigs(context.Background(), tc.pageNum, tc.pageSize, tc.filter)

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, configs)
				assert.Zero(t, total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockTotal, total)
				assert.Equal(t, tc.mockConfigs, configs)
			}
			if len(whereExpectations) > 0 {
				assert.Equal(t, len(whereExpectations)*2, whereCalls)
			} else {
				assert.Zero(t, whereCalls)
			}
		})
	}
}
