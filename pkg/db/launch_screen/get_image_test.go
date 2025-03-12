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

func TestDBLaunchScreen_GetImageById(t *testing.T) {
	type testCase struct {
		name           string
		id             int64
		mockError      error
		expectedResult *model.Picture
		expectingError bool
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
			name:           "GetImageById_Success",
			id:             2024,
			mockError:      nil,
			expectedResult: picture,
			expectingError: false,
		},
		{
			name:           "GetImageById_RecordNotFound",
			id:             9999,
			mockError:      gorm.ErrRecordNotFound,
			expectedResult: nil,
			expectingError: true,
		},
		{
			name:           "GetImageById_DBError",
			id:             3030,
			mockError:      errors.New("database error"),
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
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}

				if picture, ok := dest.(*model.Picture); ok && tc.expectedResult != nil {
					*picture = *tc.expectedResult
				}
				return mockGormDB
			}).Build()

			result, err := mockDBLaunchScreen.GetImageById(context.Background(), tc.id)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "dal.GetImageById error")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestDBLaunchScreen_GetImageBySType(t *testing.T) {
	type testCase struct {
		name           string
		sType          int64
		mockError      error
		expectedResult *[]model.Picture
		expectedCount  int64
		expectingError bool
	}

	Loc := utils.LoadCNLocation()
	now := time.Now().In(Loc)
	picture1 := model.Picture{ID: 1, SType: 1, StartAt: now.Add(-1 * time.Hour), EndAt: now.Add(2 * time.Hour), StartTime: 0, EndTime: 23}
	picture2 := model.Picture{ID: 2, SType: 1, StartAt: now.Add(-2 * time.Hour), EndAt: now.Add(1 * time.Hour), StartTime: 0, EndTime: 23}
	pictures := &[]model.Picture{picture1, picture2}

	testCases := []testCase{
		{
			name:           "GetImageBySType_Success",
			sType:          1,
			mockError:      nil,
			expectedResult: pictures,
			expectedCount:  2,
			expectingError: false,
		},
		{
			name:           "GetImageBySType_DBError",
			sType:          1,
			mockError:      errors.New("database error"),
			expectedResult: nil,
			expectedCount:  -1,
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
			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB { return mockGormDB }).Build()
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Count).To(func(count *int64) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				*count = tc.expectedCount
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Order).To(func(value interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Find).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				expectedPicture, ok := dest.(*[]model.Picture)
				if ok {
					*expectedPicture = *tc.expectedResult
				}
				return mockGormDB
			}).Build()

			result, count, err := mockDBLaunchScreen.GetImageBySType(context.Background(), tc.sType)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tc.expectedCount, count)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedCount, count)
			}
		})
	}
}

func TestDBLaunchScreen_GetImageByIdList(t *testing.T) {
	type testCase struct {
		name           string
		imgIdList      *[]int64
		mockError      error
		expectedResult *[]model.Picture
		expectedCount  int64
		expectingError bool
	}

	imgIdList := &[]int64{1, 2, 3}
	now := time.Now()
	pictures := &[]model.Picture{
		{
			ID:         1,
			Url:        "url1",
			Href:       "href1",
			Text:       "text1",
			StartAt:    now.Add(-24 * time.Hour),
			EndAt:      now.Add(24 * time.Hour),
			StartTime:  0,
			EndTime:    24,
			ShowTimes:  10,
			PointTimes: 5,
		},
		{
			ID:         2,
			Url:        "url2",
			Href:       "href2",
			Text:       "text2",
			StartAt:    now.Add(-24 * time.Hour),
			EndAt:      now.Add(24 * time.Hour),
			StartTime:  0,
			EndTime:    24,
			ShowTimes:  15,
			PointTimes: 10,
		},
	}

	testCases := []testCase{
		{
			name:           "GetImageByIdList_Success",
			imgIdList:      imgIdList,
			mockError:      nil,
			expectedResult: pictures,
			expectedCount:  2,
			expectingError: false,
		},
		{
			name:           "GetImageByIdList_DBError",
			imgIdList:      imgIdList,
			mockError:      errors.New("database error"),
			expectedResult: nil,
			expectedCount:  -1,
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
			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB { return mockGormDB }).Build()
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Count).To(func(count *int64) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				*count = tc.expectedCount
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Order).To(func(value interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Find).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				expectedPicture, ok := dest.(*[]model.Picture)
				if ok {
					*expectedPicture = *tc.expectedResult
				}
				return mockGormDB
			}).Build()

			result, count, err := mockDBLaunchScreen.GetImageByIdList(context.Background(), tc.imgIdList)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, int64(-1), count)
				assert.Contains(t, err.Error(), "dal.GetImageByIdList error")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResult, result)
				assert.Equal(t, tc.expectedCount, count)
			}
		})
	}
}

func TestDBLaunchScreen_GetLastImageId(t *testing.T) {
	type testCase struct {
		name           string
		mockError      error
		expectedResult int64
		expectingError bool
	}

	testCases := []testCase{
		{
			name:           "GetLastImageId_Success",
			mockError:      nil,
			expectedResult: 12345,
			expectingError: false,
		},
		{
			name:           "GetLastImageId_DBError",
			mockError:      errors.New("database error"),
			expectedResult: -1,
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
			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB { return mockGormDB }).Build()
			mockey.Mock((*gorm.DB).Last).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				expectedPicture, ok := dest.(*model.Picture)
				if ok {
					expectedPicture.ID = tc.expectedResult
				}
				return mockGormDB
			}).Build()

			result, err := mockDBLaunchScreen.GetLastImageId(context.Background())

			if tc.expectingError {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedResult, result)
				assert.Contains(t, err.Error(), "dal.GetLastImageId error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
