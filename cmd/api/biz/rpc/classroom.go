package rpc

import (
	"context"
	"fmt"

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/client"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitClassroomRPC() {
	client, err := client.InitClassroomRPC()
	if err != nil {
		logger.Fatalf("api.rpc.classroom InitClassroomRPC failed, err  %v", err)
	}
	classroomClient = *client
}

func GetEmptyRoomRPC(ctx context.Context, req *classroom.EmptyRoomRequest) (emptyRooms []*model.Classroom, err error) {
	resp, err := classroomClient.GetEmptyRoom(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetEmptyRoomRPC: the rpc called failed, err: %w", err)
	}
	if err = utils.IsSuccess(resp.Base); err != nil {
		return nil, fmt.Errorf("GetEmptyRoomRPC: the base code is not successful, err: %w", err)
	}
	return resp.Rooms, nil
}
