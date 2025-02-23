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

func TestDBLaunchScreen_UpdateImage(t *testing.T) {
	type testCase struct {
		name           string
		mockErrorSave  error
		mockErrorTake  error
		inputPicture   *model.Picture
		expectedResult *model.Picture
		expectingError bool
	}
	origin := &model.Picture{
		ID:         2024,
		Url:        "Url",
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
	updated := &model.Picture{
		ID:         2024,
		Url:        "newUrl",
		Href:       "newHref",
		Text:       "newText",
		PicType:    3,
		ShowTimes:  1,
		PointTimes: 1,
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
			name:           "UpdateImage_Success",
			mockErrorSave:  nil,
			mockErrorTake:  nil,
			inputPicture:   origin,
			expectedResult: updated,
			expectingError: false,
		},
		{
			name:           "UpdateImage_SaveError",
			mockErrorSave:  errors.New("save error"),
			mockErrorTake:  nil,
			inputPicture:   origin,
			expectedResult: nil,
			expectingError: true,
		},
		{
			name:           "UpdateImage_TakeError",
			mockErrorSave:  nil,
			mockErrorTake:  errors.New("take error"),
			inputPicture:   origin,
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
			mockey.Mock((*gorm.DB).Save).To(func(value interface{}) *gorm.DB {
				if tc.mockErrorSave != nil {
					mockGormDB.Error = tc.mockErrorSave
					return mockGormDB
				}
				if picture, ok := value.(*model.Picture); ok && tc.expectedResult != nil {
					picture.Url = tc.expectedResult.Url
					picture.Href = tc.expectedResult.Href
					picture.Text = tc.expectedResult.Text
					picture.PointTimes = tc.expectedResult.PointTimes
					picture.ShowTimes = tc.expectedResult.ShowTimes
				}
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Take).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockErrorTake != nil {
					mockGormDB.Error = tc.mockErrorTake
					return mockGormDB
				}
				if picture, ok := dest.(*model.Picture); ok && tc.expectedResult != nil {
					*picture = *tc.expectedResult
				}
				return mockGormDB
			}).Build()

			result, err := mockDBLaunchScreen.UpdateImage(context.Background(), tc.inputPicture)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "dal.UpdateImage error")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
