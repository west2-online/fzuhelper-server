package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/cmd/classroom/dal/cache"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
)

func TestGetEmptyRoom(t *testing.T) {
	// 测试用例结构体
	type testCase struct {
		name           string
		mockIsExist    bool
		mockReturn     interface{}
		expectedResult []string
		expectingError bool
	}

	// 测试用例列表
	tests := []testCase{
		{
			name:           "RoomInfoNotExist",
			mockIsExist:    false,
			expectedResult: nil,
			expectingError: true,
		},
		{
			name:           "RoomInfoExist",
			mockIsExist:    true,
			mockReturn:     []string{"旗山东1"},
			expectedResult: []string{"旗山东1"},
			expectingError: false,
		},
	}

	// 通用请求参数
	req := &classroom.EmptyRoomRequest{
		Date:      "2024-10-01",
		Campus:    "旗山校区",
		StartTime: "1",
		EndTime:   "1",
	}

	defer mockey.UnPatchAll()

	// 运行所有测试用例
	for _, tc := range tests {
		mockey.PatchConvey(tc.name, t, func() {
			classroomService := NewClassroomService(context.Background())

			// 根据测试用例设置 Mock 行为
			mockey.Mock(cache.IsExistRoomInfo).Return(tc.mockIsExist).Build()
			if tc.mockIsExist {
				mockey.Mock(cache.GetEmptyRoomCache).Return(tc.mockReturn, nil).Build()
			}

			// 调用 GetEmptyRoom 方法
			result, err := classroomService.GetEmptyRoom(req)

			// 根据预期的错误存在与否进行断言
			if tc.expectingError {
				assert.Nil(t, result)
				assert.EqualError(t, err, "service.GetEmptyRoom: room info not exist")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
