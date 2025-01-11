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

package user

import (
	"context"

	"github.com/west2-online/fzuhelper-server/internal/user/pack"
	"github.com/west2-online/fzuhelper-server/internal/user/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	ClientSet *base.ClientSet
}

func NewUserService(clientSet *base.ClientSet) *UserServiceImpl {
	return &UserServiceImpl{
		ClientSet: clientSet,
	}
}

// GetLoginData implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetLoginData(ctx context.Context, req *user.GetLoginDataRequest) (resp *user.GetLoginDataResponse, err error) {
	resp = new(user.GetLoginDataResponse)
	l := service.NewUserService(ctx, "", nil, s.ClientSet)
	id, cookies, err := l.GetLoginData(req)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Id = id
	resp.Cookies = cookies
	return
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, request *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	resp = new(user.GetUserInfoResponse)
	userHeader, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	l := service.NewUserService(ctx, userHeader.Id, utils.ParseCookies(userHeader.Cookies), s.ClientSet)
	info, err := l.GetUserInfo(userHeader.Id[len(userHeader.Id)-9:])
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Data = pack.BuildInfoResp(info)
	return
}
