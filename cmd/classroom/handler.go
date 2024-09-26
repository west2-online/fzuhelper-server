package main

import (
	"context"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/pack"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/service"
	classroom "github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// ClassroomServiceImpl implements the last service interface defined in the IDL.
type ClassroomServiceImpl struct{}

// GetEmptyRoom implements the ClassroomServiceImpl interface.
func (s *ClassroomServiceImpl) GetEmptyRoom(ctx context.Context, req *classroom.EmptyRoomRequest) (resp *classroom.EmptyRoomResponse, err error) {
	// TODO: Your code here...
	resp = classroom.NewEmptyRoomResponse()
	l := service.NewClassroomService(ctx, req.Logindata.Id, utils.ParseCookies(req.Logindata.Cookies))
	res, err := l.GetEmptyRooms(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Rooms = pack.BuildClassRooms(res, req.Campus)
	utils.LoggerObj.Info("GetEmptyRoom success")
	return
}
