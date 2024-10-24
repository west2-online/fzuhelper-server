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

func TestLaunchScreenService_GetImageById(t *testing.T) {
	type testCase struct {
		name           string
		mockIsExist    bool
		mockReturn     interface{}
		expectedResult interface{}
		expectingError bool
	}
	expectedResult := &db.Picture{
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
	testCases := []testCase{
		{
			name:           "GetImageById",
			mockIsExist:    true,
			mockReturn:     expectedResult,
			expectedResult: expectedResult,
		},
		{
			name:           "LaunchScreenNotExist",
			mockIsExist:    false,
			mockReturn:     gorm.ErrRecordNotFound,
			expectedResult: nil,
			expectingError: true,
		},
	}
	req := &launch_screen.GetImageRequest{
		PictureId: expectedResult.ID,
	}
	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			launchScreenService := NewLaunchScreenService(context.Background())
			if tc.mockIsExist {
				mockey.Mock(db.GetImageById).Return(tc.mockReturn, nil).Build()
			} else {
				mockey.Mock(db.GetImageById).Return(nil, tc.mockReturn).Build()
			}
			result, err := launchScreenService.GetImageById(req.PictureId)
			if tc.expectingError {
				assert.Nil(t, result)
				assert.EqualError(t, err, "LaunchScreenService.GetImageById error:record not found")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
