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

	"github.com/west2-online/fzuhelper-server/pkg/logger"

	"github.com/west2-online/fzuhelper-server/cmd/user/pack"
	"github.com/west2-online/fzuhelper-server/cmd/user/service"
	user "github.com/west2-online/fzuhelper-server/kitex_gen/user"
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
	return
}

// GetValidateCode implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetValidateCode(ctx context.Context, req *user.GetValidateCodeRequest) (resp *user.GetValidateCodeResponse, err error) {
	resp = new(user.GetValidateCodeResponse)
	code, err := service.NewUserService(ctx, "", nil).GetValidateCode(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		logger.Infof("User.GetValidateCode: %v", err)
		return resp, nil
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Code = &code
	return resp, nil
}
