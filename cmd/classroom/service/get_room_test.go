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
