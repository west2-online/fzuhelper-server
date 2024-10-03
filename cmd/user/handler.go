package main

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/west2-online/fzuhelper-server/cmd/user/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	userCli userservice.Client
}

func NewUserClient(addr string) (userservice.Client, error) {
	return userservice.NewClient(constants.UserServiceName, client.WithHostPorts(addr))
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	resp = new(user.LoginResponse)
	userResp, err := service.NewUserService(ctx).Login(req)
	resp.Base = utils.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	token, err := utils.CreateToken(userResp.ID)
	resp.Base = utils.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Token = &token
	return resp, nil
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	resp = new(user.RegisterResponse)

	userResp, err := service.NewUserService(ctx).Register(req)

	resp.Base = utils.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}

	resp.UserId = &userResp.ID

	return resp, nil
}
