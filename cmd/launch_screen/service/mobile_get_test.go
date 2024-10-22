package service

import (
	"context"
	"github.com/bytedance/mockey"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/cache"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"testing"
	"time"
)

func TestLaunchScreenService_MobileGetImage(t *testing.T) {
	type testCase struct {
		name              string
		mockIsCacheExist  bool //用于控制返回值
		mockIsCacheExpire bool
		mockExistReturn   interface{} //当exist，此字段模拟依赖结果(dal,cache)应该返回的真实数据
		mockExpireReturn  interface{}
		mockReturn        interface{}
		expectedResult    []db.Picture //期望的输出，指的是本方法调用后的输出
		expectingError    bool
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
			Regex:      "{\"device\": \"android\", \"student_id\": \"102301517\", \"etc...\"}",
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
			Regex:      "{\"device\": \"android\", \"student_id\": \"102301517\", \"etc...\"}",
		},
	}
	testCases := []testCase{
		{
			name:              "NoCache",
			mockIsCacheExist:  false,
			mockIsCacheExpire: true,
			mockExistReturn:   false,
			mockExpireReturn:  false,
			expectedResult:    expectedResult,
		},
		{
			name:              "CacheExist",
			mockIsCacheExist:  true,
			mockIsCacheExpire: false,
			mockExistReturn:   true,
			mockExpireReturn:  true,
		},
		{
			name:              "CacheExpired",
			mockIsCacheExist:  true,
			mockIsCacheExpire: true,
			mockExistReturn:   true,
			mockExpireReturn:  false,
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

			mockey.Mock(cache.IsLaunchScreenCacheExist).Return(tc.mockExistReturn).Build()
			mockey.Mock(cache.IsLastLaunchScreenIdCacheExist).Return(tc.mockExpireReturn).Build()
			if tc.mockIsCacheExpire {
				mockey.Mock(db.GetLastImageId).Return(2147483647).Build()
				mockey.Mock(cache.GetLastLaunchScreenIdCache).Return(0).Build()
			} else {
				mockey.Mock(db.GetLastImageId).Return(2147483647).Build()
				mockey.Mock(cache.GetLastLaunchScreenIdCache).Return(0).Build()
			}

		})

	}
}
