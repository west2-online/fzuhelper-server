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
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	launchScreenDB "github.com/west2-online/fzuhelper-server/pkg/db/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/oss"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestUpdateImageProperty(t *testing.T) {
	type testCase struct {
		name             string
		mockIsExist      bool
		mockOriginReturn interface{}
		mockReturn       interface{}
		expectResult     interface{}
		expectError      bool
	}

	origin := &model.Picture{
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

	expectedResult := &model.Picture{
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
			expectResult:     expectedResult,
		},
		{
			name:             "LaunchScreenNotExist",
			mockIsExist:      false,
			mockOriginReturn: gorm.ErrRecordNotFound,
			expectResult:     nil,
			expectError:      true,
		},
		{
			name:             "UpdateImage error",
			mockIsExist:      true,
			mockOriginReturn: origin,
			mockReturn:       gorm.ErrInvalidData,
			expectResult:     nil,
			expectError:      true,
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
			mockClientSet := &base.ClientSet{
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
				SFClient:    new(utils.Snowflake),
				OssSet: &oss.OSSSet{
					Provider: oss.UpYunProvider,
					Upyun:    new(oss.UpYunConfig),
				},
			}
			launchScreenService := NewLaunchScreenService(context.Background(), mockClientSet)

			if tc.mockIsExist {
				mockey.Mock((*launchScreenDB.DBLaunchScreen).GetImageById).Return(tc.mockOriginReturn, nil).Build()
			} else {
				mockey.Mock((*launchScreenDB.DBLaunchScreen).GetImageById).Return(nil, tc.mockOriginReturn).Build()
			}
			if tc.expectError && tc.name == "UpdateImage error" {
				mockey.Mock((*launchScreenDB.DBLaunchScreen).UpdateImage).Return(nil, tc.mockReturn).Build()
			} else {
				mockey.Mock((*launchScreenDB.DBLaunchScreen).UpdateImage).Return(tc.mockReturn, nil).Build()
			}

			result, err := launchScreenService.UpdateImageProperty(req)
			if tc.expectError {
				assert.Nil(t, result)
				if tc.name == "UpdateImage error" {
					assert.Error(t, err)
					assert.ErrorContains(t, err, "LaunchScreenService.UpdateImageProperty error")
				} else {
					assert.EqualError(t, err, "LaunchScreenService.UpdateImageProperty error: record not found")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}
