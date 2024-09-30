package rpc

import (
	"context"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/client"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitClassroomRPC() {
	client, err := client.InitClassroomRPC()
	if err != nil {
		logger.LoggerObj.Fatalf("api.rpc.classroom InitClassroomRPC failed, err is %v", err)
		panic(err)
	}
	classroomClient = *client
}

func GetEmptyRoomRPC(ctx context.Context, req *classroom.EmptyRoomRequest) (emptyRooms []*model.Classroom, err error) {
	resp, err := classroomClient.GetEmptyRoom(ctx, req)
	if err != nil {
		logger.LoggerObj.Errorf("api.rpc.classroom GetEmptyRoomRPC received rpc error %v", err)
		return nil, err
	}
	if err = utils.IsSuccess(resp.Base); err != nil {
		return nil, err
	}
	return resp.Rooms, nil
}
