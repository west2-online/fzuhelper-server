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

	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/cache"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestLaunchScreenService_MobileGetImage(t *testing.T) {
	type testCase struct {
		name                  string
		mockIsCacheExist      bool // 用于控制返回值
		mockIsCacheExpire     bool
		mockExistReturn       bool // 当exist，此字段模拟依赖结果(dal,cache)应该返回的真实数据
		mockExpireReturn      bool
		mockCacheReturn       []int64
		mockDbReturn          *[]db.Picture
		mockCacheLastIdReturn int64
		mockDbLastIdReturn    int64
		expectedResult        *[]db.Picture // 期望的输出，指的是本方法调用后的输出
		expectingError        bool
	}
	expectedResult := []db.Picture{
		{
			ID:         1958,
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
		},
		{
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
		},
	}
	testCases := []testCase{
		{
			name:              "NoCache",
			mockIsCacheExist:  false,
			mockIsCacheExpire: true,
			mockExistReturn:   false,
			mockExpireReturn:  false,
			expectedResult:    &expectedResult,
			mockDbReturn:      &expectedResult,
		},
		{
			name:                  "CacheExist",
			mockIsCacheExist:      true,
			mockIsCacheExpire:     false,
			mockExistReturn:       true,
			mockExpireReturn:      true,
			expectedResult:        &expectedResult,
			mockCacheReturn:       []int64{expectedResult[0].ID, expectedResult[1].ID},
			mockDbReturn:          &expectedResult,
			mockDbLastIdReturn:    expectedResult[1].ID,
			mockCacheLastIdReturn: expectedResult[1].ID,
		},
		{
			name:                  "CacheExpired",
			mockIsCacheExist:      true,
			mockIsCacheExpire:     true,
			mockExistReturn:       true,
			mockExpireReturn:      false,
			expectedResult:        &expectedResult,
			mockCacheReturn:       []int64{expectedResult[0].ID},
			mockDbReturn:          &expectedResult,
			mockCacheLastIdReturn: expectedResult[0].ID,
		},
		{
			name:              "NoLaunchScreen",
			mockIsCacheExist:  false,
			mockIsCacheExpire: true,
			mockExistReturn:   false,
			mockExpireReturn:  false,
			mockDbReturn:      &[]db.Picture{},
			expectingError:    true,
		},
	}
	req := &launch_screen.MobileGetImageRequest{
		SType:     3, // 请确保该id对应picture存在
		StudentId: 102301517,
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			launchScreenService := NewLaunchScreenService(context.Background())

			mockey.Mock(cache.IsLaunchScreenCacheExist).Return(tc.mockIsCacheExist).Build()
			mockey.Mock(cache.IsLastLaunchScreenIdCacheExist).Return(tc.mockExpireReturn).Build()
			mockey.Mock(db.GetLastImageId).Return(tc.mockDbLastIdReturn, nil).Build()
			mockey.Mock(cache.GetLastLaunchScreenIdCache).Return(tc.mockCacheLastIdReturn, nil).Build()
			mockey.Mock(cache.GetLaunchScreenCache).Return(tc.mockCacheReturn, nil).Build()
			if tc.mockIsCacheExpire {
				mockey.Mock(db.GetImageBySType).Return(tc.mockDbReturn, len(*tc.mockDbReturn), nil).Build()
				mockey.Mock(cache.SetLaunchScreenCache).Return(nil).Build()
				mockey.Mock(cache.SetLastLaunchScreenIdCache).Return(nil).Build()
			} else {
				mockey.Mock(db.GetImageByIdList).Return(tc.mockDbReturn, len(*tc.mockDbReturn), nil).Build()
			}
			mockey.Mock(db.AddImageListShowTime).Return(nil).Build()

			result, _, err := launchScreenService.MobileGetImage(req)

			if tc.expectingError {
				assert.Nil(t, result)
				assert.Equal(t, err, errno.NoRunningPictureError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
