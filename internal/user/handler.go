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

	"github.com/west2-online/fzuhelper-server/internal/user/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/base"
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
