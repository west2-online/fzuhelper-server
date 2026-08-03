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

package launch_screen

import (
	"context"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBLaunchScreen_ListImage(t *testing.T) {
	type testCase struct {
		name           string
		pageNum        int
		pageSize       int
		mockError      error
		expectedResult *[]model.Picture
		expectedCount  int64
		expectingError bool
	}

	picture := model.Picture{
		ID:        2024,
		Url:       "newUrl",
		Href:      "href",
		Text:      "text",
		PicType:   3,
		SType:     3,
		Frequency: 4,
	}
	pictures := &[]model.Picture{picture}

	testCases := []testCase{
		{
			name:           "ListImage_Success",
			pageNum:        1,
			pageSize:       20,
			mockError:      nil,
			expectedResult: pictures,
			expectedCount:  1,
			expectingError: false,
		},
		{
			name:           "ListImage_InvalidPageParam",
			pageNum:        0,
			pageSize:       0,
			expectedResult: nil,
			expectingError: true,
		},
		{
			name:           "ListImage_CountError",
			pageNum:        1,
			pageSize:       20,
			mockError:      errors.New("database error"),
			expectedResult: nil,
			expectingError: true,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			mockSnowflake := new(utils.Snowflake)
			mockDBLaunchScreen := NewDBLaunchScreen(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Count).To(func(count *int64) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				*count = tc.expectedCount
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Order).To(func(value interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Limit).To(func(limit int) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Offset).To(func(offset int) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Find).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				if pictures, ok := dest.(*[]model.Picture); ok && tc.expectedResult != nil {
					*pictures = *tc.expectedResult
				}
				return mockGormDB
			}).Build()

			result, total, err := mockDBLaunchScreen.ListImage(context.Background(), tc.pageNum, tc.pageSize)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedCount, total)
			}
		})
	}
}
