package main

import (
	"context"
	"github.com/west2-online/fzuhelper-server/cmd/user/pack"
	"github.com/west2-online/fzuhelper-server/cmd/user/service"
	user "github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// GetLoginData implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetLoginData(ctx context.Context, req *user.GetLoginDataRequest) (resp *user.GetLoginDataResponse, err error) {
	resp = new(user.GetLoginDataResponse)
	l := service.NewUserService(ctx, "", nil)
	id, cookies, err := l.GetLoginData(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Id = id
	resp.Cookies = cookies
	utils.LoggerObj.Info("GetLoginData success")
	return
}
