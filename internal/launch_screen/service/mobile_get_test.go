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

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	launchScreenCache "github.com/west2-online/fzuhelper-server/pkg/cache/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	launchScreenDB "github.com/west2-online/fzuhelper-server/pkg/db/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/oss"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestLaunchScreenService_MobileGetImage(t *testing.T) {
	type testCase struct {
		name string // 用例名
		// 控制返回值与mock函数行为
		mockIsCacheExist      bool
		mockIsCacheExpire     bool
		mockExistReturn       bool
		mockExpireReturn      bool
		mockCacheReturn       []int64
		mockCacheError        error
		mockDbReturn          *[]model.Picture
		mockDbError           error
		mockCacheLastIdReturn int64
		mockDbLastIdReturn    int64
		mockDbLastIdError     error
		mockAddShowTimeError  error
		// 期望输出
		expectedResult *[]model.Picture
		expectingError bool
		expectErrorMsg string
	}
	expectedResult := []model.Picture{
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
	// 创建测试用例，要做到覆盖大部分情况
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
			mockDbReturn:      &[]model.Picture{},
			expectingError:    true,
		},
		{
			name:              "shouldGetFromMySQL error",
			mockIsCacheExist:  true,
			mockExpireReturn:  true,
			mockDbLastIdError: assert.AnError,
			expectingError:    true,
			expectErrorMsg:    "GetLastImageId",
		},
		{
			name:              "GetLaunchScreenCache error",
			mockIsCacheExist:  true,
			mockIsCacheExpire: false,
			mockExpireReturn:  true,
			mockCacheError:    assert.AnError,
			expectingError:    true,
			expectErrorMsg:    "cache.GetLaunchScreenCache error",
		},
		{
			name:              "GetImageByIdList non-record-not-found error",
			mockIsCacheExist:  true,
			mockIsCacheExpire: false,
			mockExpireReturn:  true,
			mockCacheReturn:   []int64{1, 2},
			mockDbError:       assert.AnError,
			expectingError:    true,
			expectErrorMsg:    "db.GetImageByIdList error",
		},
		{
			name:              "GetImageByIdList returns empty result",
			mockIsCacheExist:  true,
			mockIsCacheExpire: false,
			mockExpireReturn:  true,
			mockCacheReturn:   []int64{1, 2},
			mockDbReturn:      &[]model.Picture{},
			expectingError:    true,
		},
		{
			name:                  "AddImageListShowTime error",
			mockIsCacheExist:      true,
			mockIsCacheExpire:     false,
			mockExpireReturn:      true,
			mockCacheReturn:       []int64{expectedResult[0].ID, expectedResult[1].ID},
			mockDbReturn:          &expectedResult,
			mockDbLastIdReturn:    expectedResult[1].ID,
			mockCacheLastIdReturn: expectedResult[1].ID,
			mockAddShowTimeError:  assert.AnError,
			expectingError:        true,
			expectErrorMsg:        "db.AddImageListShowTime error",
		},
	}
	// 通用请求
	req := &launch_screen.MobileGetImageRequest{
		SType:     3, // 请确保该id对应picture存在
		StudentId: "102301517",
	}

	// 用于在测试结束时确保Mock行为不会泄漏
	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		// PatchConvey封装了testCase，在其中组织testCase逻辑，同时匿名函数中的mock行为只会在函数作用域中生效
		mockey.PatchConvey(tc.name, t, func() {
			// 进行服务的初始化
			mockClientSet := new(base.ClientSet)
			mockClientSet.SFClient = new(utils.Snowflake)
			mockClientSet.DBClient = new(db.Database)
			mockClientSet.CacheClient = new(cache.Cache)
			mockClientSet.OssSet = &oss.OSSSet{Provider: oss.UpYunProvider, Upyun: new(oss.UpYunConfig)}

			launchScreenService := NewLaunchScreenService(context.Background(), mockClientSet)

			// 模拟外部依赖函数的行为，确保所以的外部函数不会影响到测试
			mockey.Mock(mockey.GetMethod(launchScreenService.cache, "IsKeyExist")).Return(tc.mockIsCacheExist).Build()
			mockey.Mock((*launchScreenCache.CacheLaunchScreen).IsLastLaunchScreenIdCacheExist).Return(tc.mockExpireReturn).Build()
			mockey.Mock((*launchScreenDB.DBLaunchScreen).GetLastImageId).Return(tc.mockDbLastIdReturn, tc.mockDbLastIdError).Build()
			mockey.Mock((*launchScreenCache.CacheLaunchScreen).GetLastLaunchScreenIdCache).Return(tc.mockCacheLastIdReturn, nil).Build()
			mockey.Mock((*launchScreenCache.CacheLaunchScreen).GetLaunchScreenCache).Return(tc.mockCacheReturn, tc.mockCacheError).Build()

			if tc.mockIsCacheExpire {
				dbCount := int64(0)
				if tc.mockDbReturn != nil {
					dbCount = int64(len(*tc.mockDbReturn))
				}
				mockey.Mock((*launchScreenDB.DBLaunchScreen).GetImageBySType).Return(tc.mockDbReturn, dbCount, nil).Build()
				mockey.Mock((*launchScreenCache.CacheLaunchScreen).SetLaunchScreenCache).Return(nil).Build()
				mockey.Mock((*launchScreenCache.CacheLaunchScreen).SetLastLaunchScreenIdCache).Return(nil).Build()
			} else {
				dbCount := int64(0)
				if tc.mockDbReturn != nil {
					dbCount = int64(len(*tc.mockDbReturn))
				}
				mockey.Mock((*launchScreenDB.DBLaunchScreen).GetImageByIdList).Return(tc.mockDbReturn, dbCount, tc.mockDbError).Build()
			}
			mockey.Mock((*launchScreenDB.DBLaunchScreen).AddImageListShowTime).Return(tc.mockAddShowTimeError).Build()

			// 得到结果
			result, _, err := launchScreenService.MobileGetImage(req)

			// 比对结果与错误
			if tc.expectingError {
				assert.Nil(t, result)
				if tc.expectErrorMsg != "" {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.expectErrorMsg)
				} else {
					assert.Equal(t, err, errno.NoRunningPictureError)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestLaunchScreenService_shouldGetFromMySQL(t *testing.T) {
	type testCase struct {
		name                 string
		studentId            string
		sType                int64
		device               string
		mockCacheExist       bool
		mockExpireCacheExist bool
		mockDbLastId         int64
		mockDbLastIdErr      error
		mockCacheLastId      int64
		mockCacheLastIdErr   error
		expectGetFromMySQL   bool
		expectError          bool
		expectErrorContains  string
	}

	testCases := []testCase{
		{
			name:               "cache not exist",
			studentId:          "102301517",
			sType:              3,
			device:             "android",
			mockCacheExist:     false,
			expectGetFromMySQL: true,
		},
		{
			name:                 "expire cache not exist",
			studentId:            "102301517",
			sType:                3,
			device:               "android",
			mockCacheExist:       true,
			mockExpireCacheExist: false,
			expectGetFromMySQL:   true,
		},
		{
			name:                 "db and cache id match",
			studentId:            "102301517",
			sType:                3,
			device:               "android",
			mockCacheExist:       true,
			mockExpireCacheExist: true,
			mockDbLastId:         100,
			mockCacheLastId:      100,
			expectGetFromMySQL:   false,
		},
		{
			name:                 "db and cache id mismatch",
			studentId:            "102301517",
			sType:                3,
			device:               "android",
			mockCacheExist:       true,
			mockExpireCacheExist: true,
			mockDbLastId:         101,
			mockCacheLastId:      100,
			expectGetFromMySQL:   true,
		},
		{
			name:                 "db get last id error",
			studentId:            "102301517",
			sType:                3,
			device:               "android",
			mockCacheExist:       true,
			mockExpireCacheExist: true,
			mockDbLastIdErr:      assert.AnError,
			expectGetFromMySQL:   true,
			expectError:          true,
			expectErrorContains:  "GetLastImageId",
		},
		{
			name:                 "cache get last id error",
			studentId:            "102301517",
			sType:                3,
			device:               "android",
			mockCacheExist:       true,
			mockExpireCacheExist: true,
			mockDbLastId:         100,
			mockCacheLastIdErr:   assert.AnError,
			expectGetFromMySQL:   true,
			expectError:          true,
			expectErrorContains:  "GetLastLaunchScreenIdCache",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient:    &db.Database{LaunchScreen: new(launchScreenDB.DBLaunchScreen)},
				CacheClient: &cache.Cache{LaunchScreen: new(launchScreenCache.CacheLaunchScreen)},
				SFClient:    new(utils.Snowflake),
				OssSet:      &oss.OSSSet{Provider: oss.UpYunProvider, Upyun: new(oss.UpYunConfig)},
			}

			svc := NewLaunchScreenService(context.Background(), mockClientSet)

			mockey.Mock((*cache.Cache).IsKeyExist).To(func(_ *cache.Cache, _ context.Context, _ string) bool {
				return tc.mockCacheExist
			}).Build()

			mockey.Mock((*launchScreenCache.CacheLaunchScreen).IsLastLaunchScreenIdCacheExist).To(
				func(_ *launchScreenCache.CacheLaunchScreen, _ context.Context, _ string) bool {
					return tc.mockExpireCacheExist
				},
			).Build()

			mockey.Mock((*launchScreenDB.DBLaunchScreen).GetLastImageId).To(
				func(_ *launchScreenDB.DBLaunchScreen, _ context.Context) (int64, error) {
					return tc.mockDbLastId, tc.mockDbLastIdErr
				},
			).Build()

			mockey.Mock((*launchScreenCache.CacheLaunchScreen).GetLastLaunchScreenIdCache).To(
				func(_ *launchScreenCache.CacheLaunchScreen, _ context.Context, _ string) (int64, error) {
					return tc.mockCacheLastId, tc.mockCacheLastIdErr
				},
			).Build()

			result, err := svc.shouldGetFromMySQL(tc.studentId, tc.sType, tc.device)

			if tc.expectError {
				assert.Error(t, err)
				if tc.expectErrorContains != "" {
					assert.Contains(t, err.Error(), tc.expectErrorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectGetFromMySQL, result)
			}
		})
	}
}

func TestLaunchScreenService_getImagesFromMySQL(t *testing.T) {
	type testCase struct {
		name                string
		studentId           string
		sType               int64
		device              string
		mockDbImages        *[]model.Picture
		mockDbCount         int64
		mockDbErr           error
		mockAddShowTimeErr  error
		mockSetCacheErr     error
		mockSetExpireErr    error
		mockGetLastIdReturn int64
		mockGetLastIdErr    error
		expectResult        *[]model.Picture
		expectCount         int64
		expectError         bool
		expectErrorContains string
	}

	matchingPicture := []model.Picture{
		{
			ID:        1,
			Url:       "url1",
			Text:      "text1",
			Regex:     `{"device": "android", "student_id": "102301517"}`,
			StartAt:   time.Now().Add(-24 * time.Hour),
			EndAt:     time.Now().Add(24 * time.Hour),
			StartTime: 0,
			EndTime:   24,
		},
	}

	noMatchPicture := []model.Picture{
		{
			ID:        2,
			Url:       "url2",
			Text:      "text2",
			Regex:     `{"device": "ios", "student_id": "999999999"}`,
			StartAt:   time.Now().Add(-24 * time.Hour),
			EndAt:     time.Now().Add(24 * time.Hour),
			StartTime: 0,
			EndTime:   24,
		},
	}

	mixedPictures := []model.Picture{
		matchingPicture[0],
		noMatchPicture[0],
		{
			ID:        3,
			Url:       "url3",
			Text:      "text3",
			Regex:     `{"device": "android,ios", "student_id": ""}`,
			StartAt:   time.Now().Add(-24 * time.Hour),
			EndAt:     time.Now().Add(24 * time.Hour),
			StartTime: 0,
			EndTime:   24,
		},
	}

	testCases := []testCase{
		{
			name:                "no images in db",
			studentId:           "102301517",
			sType:               3,
			device:              "android",
			mockDbImages:        &[]model.Picture{},
			mockDbCount:         0,
			expectCount:         0,
			expectError:         true,
			expectErrorContains: "没有可用图片",
		},
		{
			name:                "db query error",
			studentId:           "102301517",
			sType:               3,
			device:              "android",
			mockDbErr:           assert.AnError,
			expectError:         true,
			expectErrorContains: "GetImageBySType",
		},
		{
			name:                "matching images found",
			studentId:           "102301517",
			sType:               3,
			device:              "android",
			mockDbImages:        &matchingPicture,
			mockDbCount:         1,
			mockGetLastIdReturn: 1,
			expectResult:        &matchingPicture,
			expectCount:         1,
		},
		{
			name:                "no matching images",
			studentId:           "102301517",
			sType:               3,
			device:              "android",
			mockDbImages:        &noMatchPicture,
			mockDbCount:         1,
			expectCount:         0,
			expectError:         true,
			expectErrorContains: "没有可用图片",
		},
		{
			name:                "mixed images with matches",
			studentId:           "102301517",
			sType:               3,
			device:              "android",
			mockDbImages:        &mixedPictures,
			mockDbCount:         3,
			mockGetLastIdReturn: 3,
			expectCount:         2,
		},
		{
			name:                "set cache error",
			studentId:           "102301517",
			sType:               3,
			device:              "android",
			mockDbImages:        &matchingPicture,
			mockDbCount:         1,
			mockSetCacheErr:     assert.AnError,
			expectError:         true,
			expectErrorContains: "set cache error",
		},
		{
			name:                "add show time error",
			studentId:           "102301517",
			sType:               3,
			device:              "android",
			mockDbImages:        &matchingPicture,
			mockDbCount:         1,
			mockAddShowTimeErr:  assert.AnError,
			expectError:         true,
			expectErrorContains: "set cache error",
		},
		{
			name:                "get last id error",
			studentId:           "102301517",
			sType:               3,
			device:              "android",
			mockDbImages:        &matchingPicture,
			mockDbCount:         1,
			mockGetLastIdErr:    assert.AnError,
			expectError:         true,
			expectErrorContains: "set cache error",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient:    &db.Database{LaunchScreen: new(launchScreenDB.DBLaunchScreen)},
				CacheClient: &cache.Cache{LaunchScreen: new(launchScreenCache.CacheLaunchScreen)},
				SFClient:    new(utils.Snowflake),
				OssSet:      &oss.OSSSet{Provider: oss.UpYunProvider, Upyun: new(oss.UpYunConfig)},
			}

			svc := NewLaunchScreenService(context.Background(), mockClientSet)

			mockey.Mock((*launchScreenDB.DBLaunchScreen).GetImageBySType).To(
				func(_ *launchScreenDB.DBLaunchScreen, _ context.Context, _ int64) (*[]model.Picture, int64, error) {
					return tc.mockDbImages, tc.mockDbCount, tc.mockDbErr
				},
			).Build()

			mockey.Mock((*launchScreenDB.DBLaunchScreen).AddImageListShowTime).To(
				func(_ *launchScreenDB.DBLaunchScreen, _ context.Context, _ *[]model.Picture) error {
					return tc.mockAddShowTimeErr
				},
			).Build()

			mockey.Mock((*launchScreenCache.CacheLaunchScreen).SetLaunchScreenCache).To(
				func(_ *launchScreenCache.CacheLaunchScreen, _ context.Context, _ string, _ *[]int64) error {
					return tc.mockSetCacheErr
				},
			).Build()

			mockey.Mock((*launchScreenDB.DBLaunchScreen).GetLastImageId).To(
				func(_ *launchScreenDB.DBLaunchScreen, _ context.Context) (int64, error) {
					return tc.mockGetLastIdReturn, tc.mockGetLastIdErr
				},
			).Build()

			mockey.Mock((*launchScreenCache.CacheLaunchScreen).SetLastLaunchScreenIdCache).To(
				func(_ *launchScreenCache.CacheLaunchScreen, _ context.Context, _ int64, _ string) error {
					return tc.mockSetExpireErr
				},
			).Build()

			result, count, err := svc.getImagesFromMySQL(tc.studentId, tc.sType, tc.device)

			if tc.expectError {
				assert.Error(t, err)
				if tc.expectErrorContains != "" {
					assert.Contains(t, err.Error(), tc.expectErrorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectCount, count)
				if tc.expectResult != nil {
					assert.Equal(t, len(*tc.expectResult), len(*result))
				}
			}
		})
	}
}
