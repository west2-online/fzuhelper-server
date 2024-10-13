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

// Code generated by hertz generator.

package api

import (
	"context"

	"github.com/west2-online/fzuhelper-server/cmd/api/biz/pack"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	api "github.com/west2-online/fzuhelper-server/cmd/api/biz/model/api"
)

// GetLoginData .
// @router /api/v1/user/login [GET]
func GetLoginData(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetLoginDataRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.GetLoginData: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}
	resp := new(api.GetLoginDataResponse)
	id, cookies, err := rpc.GetLoginDataRPC(ctx, &user.GetLoginDataRequest{
		Id:       req.ID,
		Password: req.Password,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.ID = id
	resp.Cookies = cookies
	pack.RespData(c, resp)
}

// Login .
// @Summary Login
// @Description login to get token
// @Accept json/form
// @Produce json
// @Param account query string true "account"
// @Param password query string true "密码"
// @router /launch_screen/api/login [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.LoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.LoginResponse)
	token, err := rpc.LoginRPC(ctx, &user.LoginRequest{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}

	resp.Base = pack.BuildSuccessResp()
	resp.Token = token
	c.JSON(consts.StatusOK, resp)
}

// Register .
// @Summary Register
// @Description userRegister
// @Accept json/form
// @Produce json
// @Param account query string true "account"
// @Param name query string true "name"
// @Param password query string true "密码"
// @router /launch_screen/api/register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.RegisterRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.RegisterResponse)
	uid, err := rpc.RegisterRPC(ctx, &user.RegisterRequest{
		Account:  req.Account,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}

	resp.Base = pack.BuildSuccessResp()
	resp.UserID = uid
	c.JSON(consts.StatusOK, resp)
}
