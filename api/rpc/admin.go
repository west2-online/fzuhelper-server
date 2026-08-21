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

package rpc

import (
	"context"

	"github.com/west2-online/fzuhelper-server/kitex_gen/admin"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitAdminRPC() {
	c, err := client.InitAdminRPC()
	if err != nil {
		logger.Fatalf("api.rpc.admin InitAdminRPC failed, err is %v", err)
	}
	adminClient = *c
}

func SSOLoginRPC(ctx context.Context, req *admin.LoginRequest) (string, error) {
	resp, err := adminClient.Login(ctx, req)
	if err != nil {
		logger.WithCtx(ctx).Errorf("SSOLoginRPC: RPC called failed: %v", err.Error())
		return "", errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return "", errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.AuthorizationUrl, nil
}

func AdminCallbackRPC(ctx context.Context, req *admin.CallbackRequest) (string, error) {
	resp, err := adminClient.Callback(ctx, req)
	if err != nil {
		logger.WithCtx(ctx).Errorf("AdminCallbackRPC: RPC called failed: %v", err.Error())
		return "", errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return "", errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.RedirectUrl, nil
}

func AdminExchangeTicketRPC(ctx context.Context, req *admin.ExchangeTicketRequest) (string, error) {
	resp, err := adminClient.ExchangeTicket(ctx, req)
	if err != nil {
		logger.WithCtx(ctx).Errorf("AdminExchangeTicketRPC: RPC called failed: %v", err.Error())
		return "", errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return "", errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.AccessToken, nil
}
