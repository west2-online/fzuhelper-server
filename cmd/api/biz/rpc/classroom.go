package rpc

import (
	"context"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitClassroomRPC() {
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}
	classroomClient, err = classroomservice.NewClient("classroom", client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

func GetEmptyRoomRPC(ctx context.Context, req *classroom.EmptyRoomRequest) (emptyRooms []*classroom.Classroom, err error) {
	resp, err := classroomClient.GetEmptyRoom(ctx, req)
	if err != nil {
		utils.LoggerObj.Errorf("api.rpc.classroom GetEmptyRoomRPC received rpc error %v", err)
		return nil, err
	}
	if resp.Base.Code != errno.SuccessCode {
		utils.LoggerObj.Errorf("api.rpc.classroom GetEmptyRoomRPC received failed")
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Rooms, nil
}
