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

	"github.com/west2-online/fzuhelper-server/internal/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestLaunchScreenService_CreateImage(t *testing.T) {
	type testCase struct {
		name            string
		mockIsExist     bool
		mockCloudReturn interface{}
		mockReturn      interface{}
		expectedResult  interface{}
		expectingError  bool
	}
	expectedResult := &db.Picture{
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
			name:            "CreateImage",
			mockIsExist:     true,
			mockReturn:      expectedResult,
			mockCloudReturn: nil,
			expectedResult:  expectedResult,
		},
		{
			name:            "cloudFail",
			mockIsExist:     true,
			mockReturn:      expectedResult,
			mockCloudReturn: errno.UpcloudError,
			expectedResult:  nil,
			expectingError:  true,
		},
	}
	req := &launch_screen.CreateImageRequest{
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
		Regex:     expectedResult.Regex,
	}
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			launchScreenService := NewLaunchScreenService(context.Background())
			db.SF = &utils.Snowflake{}
			mockey.Mock(db.SF.NextVal).To(func() (int64, error) { return expectedResult.ID, nil }).Build()
			mockey.Mock(upyun.GenerateImgName).Return(expectedResult.Url).Build()
			mockey.Mock(db.CreateImage).Return(tc.mockReturn, nil).Build()
			mockey.Mock(upyun.UploadImg).Return(tc.mockCloudReturn).Build()

			result, err := launchScreenService.CreateImage(req)

			if tc.expectingError {
				assert.Nil(t, result)
				assert.EqualError(t, err, "LaunchScreenService.CreateImage error:[40006] upload to upcloud error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
