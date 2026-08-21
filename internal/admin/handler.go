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

package admin

import (
	"context"

	"github.com/west2-online/fzuhelper-server/internal/admin/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/admin"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// AdminServiceImpl implements the last service interface defined in the IDL.
type AdminServiceImpl struct {
	ClientSet *base.ClientSet
}

func NewAdminService(clientSet *base.ClientSet) *AdminServiceImpl {
	return &AdminServiceImpl{ClientSet: clientSet}
}

// Login implements the AdminServiceImpl interface.
func (s *AdminServiceImpl) Login(ctx context.Context, req *admin.LoginRequest) (resp *admin.LoginResponse, err error) {
	resp = new(admin.LoginResponse)
	authorizeUrl, err := service.NewAdminService(ctx, s.ClientSet).Login(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Errorf("Admin.Login: %v", err)
		return resp, nil
	}
	resp.AuthorizationUrl = authorizeUrl
	return resp, nil
}

// Callback implements the AdminServiceImpl interface.
func (s *AdminServiceImpl) Callback(ctx context.Context, req *admin.CallbackRequest) (resp *admin.CallbackResponse, err error) {
	resp = new(admin.CallbackResponse)
	redirectUrl, err := service.NewAdminService(ctx, s.ClientSet).Callback(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Errorf("Admin.Callback: %v", err)
		return resp, nil
	}
	resp.RedirectUrl = redirectUrl
	return resp, nil
}

// ExchangeTicket implements the AdminServiceImpl interface.
func (s *AdminServiceImpl) ExchangeTicket(ctx context.Context, req *admin.ExchangeTicketRequest) (resp *admin.ExchangeTicketResponse, err error) {
	resp = new(admin.ExchangeTicketResponse)
	accessToken, err := service.NewAdminService(ctx, s.ClientSet).ExchangeTicket(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		logger.WithCtx(ctx).Errorf("Admin.ExchangeTicket: %v", err)
		return resp, nil
	}
	resp.AccessToken = accessToken
	return resp, nil
}
