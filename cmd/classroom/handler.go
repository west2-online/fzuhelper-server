package main

import (
	"context"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/pack"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/service"
	classroom "github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// ClassroomServiceImpl implements the last service interface defined in the IDL.
type ClassroomServiceImpl struct{}

// GetEmptyRoom implements the ClassroomServiceImpl interface.
func (s *ClassroomServiceImpl) GetEmptyRoom(ctx context.Context, req *classroom.EmptyRoomRequest) (resp *classroom.EmptyRoomResponse, err error) {
	resp = classroom.NewEmptyRoomResponse()
	l := service.NewClassroomService(ctx)
	res, err := l.GetEmptyRoom(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Rooms = pack.BuildClassRooms(res, req.Campus)
	logger.LoggerObj.Info("GetEmptyRoom success")
	return
}
