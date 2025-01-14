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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/mw"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

// GetLoginData .
// @router /api/v1/user/login [GET]
func GetLoginData(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetLoginDataRequest
	err = c.BindAndValidate(&req)
	if err != nil {
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

// ValidateCode .
// @router /api/v1/jwch/user/validateCode [POST]
func ValidateCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.ValidateCodeRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	request := new(protocol.Request)
	request.SetMethod(consts.MethodPost)
	request.SetRequestURI(constants.ValidateCodeURL)
	request.SetFormData(
		map[string]string{
			"image": req.Image,
		},
	)

	res := new(protocol.Response)

	if err = clientSet.HzClient.Do(ctx, request, res); err != nil {
		pack.RespError(c, err)
		return
	}

	c.String(http.StatusOK, res.BodyBuffer().String())
}

// ValidateCodeForAndroid .
// @router /api/login/validateCode [POST]
func ValidateCodeForAndroid(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.ValidateCodeForAndroidRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	request := new(protocol.Request)
	request.SetMethod(consts.MethodPost)
	request.SetRequestURI(constants.ValidateCodeURL)
	request.SetFormData(
		map[string]string{
			"image": req.ValidateCode,
		},
	)

	res := new(protocol.Response)

	if err = clientSet.HzClient.Do(ctx, request, res); err != nil {
		pack.RespError(c, err)
		return
	}
	// 旧版 Android 使用 message 作为解析后的验证码结果
	var originalResponse struct {
		Code    string `json:"code"`
		Data    string `json:"data"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(res.BodyBuffer().Bytes(), &originalResponse); err != nil {
		logger.Errorf("api.ValidateCodeForAndroid: JSON unmarshal error %v", err)
		pack.RespError(c, err)
		return
	}

	// 构建兼容格式的响应
	compatResponse := map[string]string{
		"code":    originalResponse.Code,
		"message": originalResponse.Data, // 将解析的验证码作为 message 返回
	}

	c.JSON(http.StatusOK, compatResponse)
}

// RefreshToken 利用 RefreshToken 刷新 AccessToken，如果类型不是 RefreshToken 会拒绝刷新
// @router /api/v1/login/refreshToken [POST]
func RefreshToken(ctx context.Context, c *app.RequestContext) {
	token := string(c.GetHeader(constants.AuthHeader))
	if token == "" {
		pack.RespError(c, errno.AuthMissing)
		return
	}
	tokenType, err := mw.CheckToken(token)
	if err != nil {
		pack.RespError(c, err)
		return
	}
	if tokenType != constants.TypeRefreshToken {
		pack.RespError(c, errno.AuthMissing.WithMessage("token type is access token, need refresh token"))
		return
	}
	access, refresh, err := mw.CreateAllToken()
	if err != nil {
		pack.RespError(c, err)
		return
	}
	c.Header(constants.AccessTokenHeader, access)
	c.Header(constants.RefreshTokenHeader, refresh)
	pack.RespSuccess(c)
}

// GetToken 基于用户校验刷新两个 Token
// @router /api/v1/login/access-token [POST]
func GetToken(ctx context.Context, c *app.RequestContext) {
	// 这个 ID 通常是 202412615623052106112，可以明显看到学号和日期，我们截取后 9 位作为学号来验证活跃
	identifier := c.Request.Header.Get("id")
	id := identifier[len(identifier)-9:]
	cookies := c.Request.Header.Get("cookies")

	err := jwch.NewStudent().
		WithUser(id, "").
		WithLoginData(identifier, utils.ParseCookies(cookies)).
		CheckSession()
	if err != nil {
		pack.RespError(c, errno.AuthError.WithMessage(fmt.Sprintf("check id and session failed, err: %v", err)))
		return
	}

	access, refresh, err := mw.CreateAllToken()
	if err != nil {
		pack.RespError(c, err)
		return
	}
	c.Header(constants.AccessTokenHeader, access)
	c.Header(constants.RefreshTokenHeader, refresh)
	pack.RespSuccess(c)
}

// TestAuth 测试鉴权功能
// @router api/v1/login/ping [GET]
func TestAuth(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    errno.SuccessCode,
		"message": errno.SuccessMsg,
	})
}

// GetUserInfo .
// @router /api/v1/jwch/user/info [GET]
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	info, err := rpc.GetUserInfoRPC(ctx, &user.GetUserInfoRequest{})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.RespData(c, info)
}
