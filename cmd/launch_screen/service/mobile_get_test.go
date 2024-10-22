package service

import (
	"context"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/cache"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"testing"
	"time"
)

func TestLaunchScreenService_MobileGetImage(t *testing.T) {
	type testCase struct {
		name                  string
		mockIsCacheExist      bool //用于控制返回值
		mockIsCacheExpire     bool
		mockExistReturn       bool //当exist，此字段模拟依赖结果(dal,cache)应该返回的真实数据
		mockExpireReturn      bool
		mockCacheReturn       []int64
		mockSqlReturn         *[]db.Picture
		mockCacheLastIdReturn int64
		mockSqlLastIdReturn   int64
		expectedResult        *[]db.Picture //期望的输出，指的是本方法调用后的输出
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
			mockSqlReturn:     &expectedResult,
		},
		{
			name:                  "CacheExist",
			mockIsCacheExist:      true,
			mockIsCacheExpire:     false,
			mockExistReturn:       true,
			mockExpireReturn:      true,
			expectedResult:        &expectedResult,
			mockCacheReturn:       []int64{expectedResult[0].ID, expectedResult[1].ID},
			mockSqlReturn:         &expectedResult,
			mockSqlLastIdReturn:   expectedResult[1].ID,
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
			mockSqlReturn:         &expectedResult,
			mockCacheLastIdReturn: expectedResult[0].ID,
		},
		{
			name:              "noLaunchScreen",
			mockIsCacheExist:  false,
			mockIsCacheExpire: true,
			mockExistReturn:   false,
			mockExpireReturn:  false,
			mockSqlReturn:     &[]db.Picture{},
			expectingError:    true,
		},
	}
	req := &launch_screen.MobileGetImageRequest{
		SType:     3, //请确保该id对应picture存在
		StudentId: 102301517,
	}
	defer mockey.UnPatchAll() //撤销所有mock操作，不会影响其他测试

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			launchScreenService := NewLaunchScreenService(context.Background())

			mockey.Mock(cache.IsLaunchScreenCacheExist).Return(tc.mockIsCacheExist).Build()
			mockey.Mock(cache.IsLastLaunchScreenIdCacheExist).Return(tc.mockExpireReturn).Build()
			mockey.Mock(db.GetLastImageId).Return(tc.mockSqlLastIdReturn, nil).Build()
			mockey.Mock(cache.GetLastLaunchScreenIdCache).Return(tc.mockCacheLastIdReturn, nil).Build()
			mockey.Mock(cache.GetLaunchScreenCache).Return(tc.mockCacheReturn, nil).Build()
			if tc.mockIsCacheExpire {
				mockey.Mock(db.GetImageBySType).Return(tc.mockSqlReturn, len(*tc.mockSqlReturn), nil).Build()
				mockey.Mock(cache.SetLaunchScreenCache).Return(nil).Build()
				mockey.Mock(cache.SetLastLaunchScreenIdCache).Return(nil).Build()
			} else {
				mockey.Mock(db.GetImageByIdList).Return(tc.mockSqlReturn, len(*tc.mockSqlReturn), nil).Build()
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
