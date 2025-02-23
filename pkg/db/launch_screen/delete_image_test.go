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
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestDBLaunchScreen_DeleteImage(t *testing.T) {
	type testCase struct {
		name            string
		inputID         int64
		mockErrorTake   error
		mockErrorDelete error
		expectedResult  *model.Picture
		expectingError  bool
	}
	picture := &model.Picture{
		ID:         2024,
		Url:        "newUrl",
		Href:       "href",
		Text:       "text",
		PicType:    3,
		ShowTimes:  0,
		PointTimes: 0,
		Duration:   0,
		StartAt:    time.Now().Add(-24 * time.Hour),
		EndAt:      time.Now().Add(24 * time.Hour),
		StartTime:  0,
		EndTime:    24,
		SType:      3,
		Frequency:  4,
		Regex:      "{\"device\": \"android,ios\", \"student_id\": \"102301517,102301544\"}",
	}
	testCases := []testCase{
		{
			name:            "DeleteImage_Success",
			inputID:         1001,
			mockErrorTake:   nil,
			mockErrorDelete: nil,
			expectedResult:  picture,
			expectingError:  false,
		},
		{
			name:            "DeleteImage_TakeError",
			inputID:         1001,
			mockErrorTake:   errors.New("record not found"),
			mockErrorDelete: nil,
			expectedResult:  nil,
			expectingError:  true,
		},
		{
			name:            "DeleteImage_DeleteError",
			inputID:         1001,
			mockErrorTake:   nil,
			mockErrorDelete: errors.New("delete failed"),
			expectedResult:  picture,
			expectingError:  true,
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
			mockey.Mock((*gorm.DB).Take).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockErrorTake != nil {
					mockGormDB.Error = tc.mockErrorTake
					return mockGormDB
				}
				expectedPicture, ok := dest.(*model.Picture)
				if ok {
					*expectedPicture = *tc.expectedResult
				}
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Delete).To(func(value interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockErrorDelete != nil {
					mockGormDB.Error = tc.mockErrorDelete
					return mockGormDB
				}
				return mockGormDB
			}).Build()

			result, err := mockDBLaunchScreen.DeleteImage(context.Background(), tc.inputID)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "dal.DeleteImage error")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
