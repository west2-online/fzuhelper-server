package rpc

import (
	"context"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitClassroomRPC() {
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}
	classroomClient, err = classroomservice.NewClient("classroom", client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		panic(err)
	}
	utils.LoggerObj.Info("InitClassroomRPC success")
}

func GetEmptyRoomRPC(ctx context.Context, req *classroom.EmptyRoomRequest) (emptyRooms []*model.Classroom, err error) {
	resp, err := classroomClient.GetEmptyRoom(ctx, req)
	if err != nil {
		utils.LoggerObj.Errorf("api.rpc.classroom GetEmptyRoomRPC received rpc error %v", err)
		return nil, err
	}
	if err = utils.IsSuccess(resp.Base); err != nil {
		return nil, err
	}
	return resp.Rooms, nil
}
