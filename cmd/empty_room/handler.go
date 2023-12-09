package main

import (
	"context"

	"github.com/west2-online/fzuhelper-server/cmd/empty_room/pack"
	"github.com/west2-online/fzuhelper-server/cmd/empty_room/service"
	empty_room "github.com/west2-online/fzuhelper-server/kitex_gen/empty_room"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// EmptyRoomServiceImpl implements the last service interface defined in the IDL.
type EmptyRoomServiceImpl struct{}

// GetEmptyRoom implements the EmptyRoomServiceImpl interface.
func (s *EmptyRoomServiceImpl) GetEmptyRoom(ctx context.Context, req *empty_room.EmptyRoomRequest) (resp *empty_room.EmptyRoomResponse, err error) {
	resp = new(empty_room.EmptyRoomResponse)

	if _, err := utils.CheckToken(req.Token); err != nil {
		resp.Base = pack.BuildBaseResp(errno.AuthFailedError)
		return resp, nil
	}

	if req.Account == nil && req.Password == nil {
		// 当未传入账号密码时(例如研究生获取空教室),使用默认账号密码
		*req.Account = constants.DefaultAccount
		*req.Password = constants.DefaultPassword
	}

	empty_room, err := service.NewEmptyRoomService(ctx).GetRoom(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.RoomName = empty_room
	return
}
