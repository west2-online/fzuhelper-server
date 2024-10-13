/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"

	"github.com/west2-online/fzuhelper-server/cmd/user/pack"
	"github.com/west2-online/fzuhelper-server/pkg/utils"

	"github.com/west2-online/fzuhelper-server/cmd/user/service"
	user "github.com/west2-online/fzuhelper-server/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

/*
userCli userservice.Client
func NewUserClient(addr string) (userservice.Client, error) {
	return userservice.NewClient(constants.UserServiceName, client.WithHostPorts(addr))
}
*/

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
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	resp = new(user.LoginResponse)
	userResp, err := service.NewUserService(ctx, "", nil).Login(req)
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

	userResp, err := service.NewUserService(ctx, "", nil).Register(req)

	resp.Base = utils.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}

	resp.UserId = &userResp.ID

	return resp, nil
}
