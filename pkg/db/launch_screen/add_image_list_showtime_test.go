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
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBLaunchScreen_AddImageListShowTime(t *testing.T) {
	type testCase struct {
		name                string
		mockError           error
		inputPictureList    *[]model.Picture
		expectedPictureList *[]model.Picture
		expectingError      bool
	}

	testCases := []testCase{
		{
			name:      "AddImageListShowTime_Success",
			mockError: nil,
			inputPictureList: &[]model.Picture{
				{ID: 1, Url: "https://example.com/image1.jpg", ShowTimes: 2},
				{ID: 2, Url: "https://example.com/image2.jpg", ShowTimes: 5},
			},
			expectedPictureList: &[]model.Picture{
				{ID: 1, Url: "https://example.com/image1.jpg", ShowTimes: 3},
				{ID: 2, Url: "https://example.com/image2.jpg", ShowTimes: 6},
			},
			expectingError: false,
		},
		{
			name:      "AddImageListShowTime_DBError",
			mockError: fmt.Errorf("db error"),
			inputPictureList: &[]model.Picture{
				{ID: 1, Url: "https://example.com/image1.jpg", ShowTimes: 2},
			},
			expectedPictureList: nil,
			expectingError:      true,
		},
		{
			name:                "AddImageListShowTime_EmptyList",
			mockError:           nil,
			inputPictureList:    &[]model.Picture{},
			expectedPictureList: &[]model.Picture{},
			expectingError:      false,
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
			mockey.Mock((*gorm.DB).Save).To(func(value interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				return mockGormDB
			}).Build()

			err := mockDBLaunchScreen.AddImageListShowTime(context.Background(), tc.inputPictureList)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "dal.AddImageListShowTime error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPictureList, tc.inputPictureList)
			}
		})
	}
}
