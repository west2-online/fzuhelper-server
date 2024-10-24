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

package service

import (
	"context"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
)

func TestLaunchScreenService_UpdateImageProperty(t *testing.T) {
	type testCase struct {
		name             string
		mockIsExist      bool
		mockOriginReturn interface{}
		mockReturn       interface{}
		expectedResult   interface{}
		expectingError   bool
	}
	origin := &db.Picture{
		ID:         2024,
		Url:        "url",
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
	expectedResult := &db.Picture{
		ID:         2024,
		Url:        "url",
		Href:       "href",
		Text:       "text",
		PicType:    3,
		ShowTimes:  0,
		PointTimes: 0,
		Duration:   6,
		StartAt:    time.Now().Add(-24 * time.Hour),
		EndAt:      time.Now().Add(24 * time.Hour),
		StartTime:  6,
		EndTime:    24,
		SType:      2,
		Frequency:  6,
		Regex:      "{\"device\": \"android,ios\", \"student_id\": \"102301517,102301545\"}",
	}
	testCases := []testCase{
		{
			name:             "UpdateImageProperty",
			mockIsExist:      true,
			mockOriginReturn: origin,
			mockReturn:       expectedResult,
			expectedResult:   expectedResult,
		},
		{
			name:             "LaunchScreenNotExist",
			mockIsExist:      false,
			mockOriginReturn: gorm.ErrRecordNotFound,
			expectedResult:   nil,
			expectingError:   true,
		},
	}
	req := &launch_screen.ChangeImagePropertyRequest{
		PicType:   expectedResult.PicType,
		Duration:  &expectedResult.Duration,
		Href:      &expectedResult.Href,
		StartAt:   expectedResult.StartAt.Unix(),
		EndAt:     expectedResult.EndAt.Unix(),
		SType:     expectedResult.SType,
		Frequency: expectedResult.Frequency,
		StartTime: expectedResult.StartTime,
		EndTime:   expectedResult.EndTime,
		Text:      expectedResult.Text,
		PictureId: expectedResult.ID,
		Regex:     expectedResult.Regex,
	}
	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			launchScreenService := NewLaunchScreenService(context.Background())
			if tc.mockIsExist {
				mockey.Mock(db.GetImageById).Return(tc.mockOriginReturn, nil).Build()
			} else {
				mockey.Mock(db.GetImageById).Return(nil, tc.mockOriginReturn).Build()
			}
			mockey.Mock(db.UpdateImage).Return(tc.mockReturn, nil).Build()

			result, err := launchScreenService.UpdateImageProperty(req)
			if tc.expectingError {
				assert.Nil(t, result)
				assert.EqualError(t, err, "LaunchScreenService.UpdateImageProperty error: record not found")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
