package service

import (
	"context"

	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/dal/cache"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
)

func TestGetEmptyRoom_RoomInfoNotExist(t *testing.T) {
	classroomService := NewClassroomService(context.Background())

	// 模拟请求参数
	req := &classroom.EmptyRoomRequest{
		Date:      "2024-10-01",
		Campus:    "旗山校区",
		StartTime: "1",
		EndTime:   "1",
	}

	// Mock cache.IsExistRoomInfo 返回 false
	mockey.Mock(cache.IsExistRoomInfo).Return(false).Build()

	result, err := classroomService.GetEmptyRoom(req)

	// 断言返回结果
	assert.Nil(t, result)
	assert.EqualError(t, err, "service.GetEmptyRoom: room info not exist")
}

func TestGetEmptyRoom_RoomInfoExist(t *testing.T) {
	classroomService := NewClassroomService(context.Background())

	req := &classroom.EmptyRoomRequest{
		Date:      "2024-10-01",
		Campus:    "旗山校区",
		StartTime: "1",
		EndTime:   "1",
	}

	// Mock cache.IsExistRoomInfo 返回 true
	mockey.Mock(cache.IsExistRoomInfo).Return(true).Build()

	expectedRooms := []string{"旗山东1"}
	mockey.Mock(cache.GetEmptyRoomCache).Return(expectedRooms, nil).Build()

	// 调用 GetEmptyRoom 方法
	result, err := classroomService.GetEmptyRoom(req)

	// 断言返回结果
	assert.NoError(t, err)
	assert.Equal(t, expectedRooms, result)
}
