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

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestLaunchScreenService_UpdateImagePath(t *testing.T) {
	type testCase struct {
		name             string
		mockIsExist      bool
		mockOriginReturn interface{}
		mockCloudReturn  interface{}
		mockReturn       interface{}
		expectedResult   interface{}
		expectingError   bool
	}
	origin := &model.Picture{
		ID:         2024,
		Url:        "oldUrl",
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
	expectedResult := &model.Picture{
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
			name:             "UpdateImagePath",
			mockIsExist:      true,
			mockOriginReturn: origin,
			mockCloudReturn:  nil,
			mockReturn:       expectedResult,
			expectedResult:   expectedResult,
		},
		{
			name:             "LaunchScreenNotExist",
			mockIsExist:      false,
			mockOriginReturn: gorm.ErrRecordNotFound,
			mockCloudReturn:  nil,
			mockReturn:       nil,
			expectedResult:   nil,
			expectingError:   true,
		},
		{
			name:             "cloudFail",
			mockIsExist:      true,
			mockCloudReturn:  errno.UpcloudError,
			mockOriginReturn: origin,
			mockReturn:       expectedResult,
			expectedResult:   nil,
			expectingError:   true,
		},
	}
	req := &launch_screen.ChangeImageRequest{}
	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			launchScreenService := NewLaunchScreenService(context.Background(), nil)
			launchScreenService.sf = &utils.Snowflake{}
			launchScreenService.db = &db.Database{}
			launchScreenService.cache = &cache.Cache{}

			if tc.mockIsExist {
				mockey.Mock(launchScreenService.db.LaunchScreen.GetImageById).Return(tc.mockOriginReturn, nil).Build()
			} else {
				mockey.Mock(launchScreenService.db.LaunchScreen.GetImageById).Return(nil, tc.mockOriginReturn).Build()
			}
			mockey.Mock(upyun.DeleteImg).Return(tc.mockCloudReturn).Build()
			mockey.Mock(upyun.UploadImg).Return(tc.mockCloudReturn).Build()

			mockey.Mock(launchScreenService.db.LaunchScreen.UpdateImage).Return(tc.mockReturn, nil).Build()
			result, err := launchScreenService.UpdateImagePath(req)

			if tc.expectingError {
				assert.Nil(t, result)
				if !tc.mockIsExist {
					assert.EqualError(t, err, "LaunchScreenService.UpdateImagePath db.GetImageById error: record not found")
				} else {
					assert.EqualError(t, err, "LaunchScreenService.UpdateImagePath error: [40006] upload to upcloud error")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
