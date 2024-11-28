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

package custom

import (
	"context"
	"net/http"

	"github.com/west2-online/fzuhelper-server/pkg/constants"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	api "github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/url"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

var ClientSet *base.ClientSet

// 1127-custom by FantasyRL

// APILogin .
// @router /api/v1/url/login [POST]
func APILogin(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.LoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.APILogin: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	err = rpc.LoginRPC(ctx, &url.LoginRequest{Password: req.Password})
	if err != nil {
		if errNo := errno.ConvertErr(err); errNo.ErrorCode == http.StatusUnauthorized {
			c.String(consts.StatusOK, constants.UrlCustomErrorMsg)
			return
		}
		pack.RespError(c, err)
		return
	}

	c.JSON(consts.StatusOK, consts.StatusOK)
}

// UploadVersionInfo .
// @router /api/v1/url/api/upload [POST]
func UploadVersionInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.UploadRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.UploadVersion: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	// resp := new(api.UploadResponse)

	err = rpc.UploadVersionRPC(ctx, &url.UploadRequest{
		Version:  req.Version,
		Code:     req.Code,
		Url:      req.URL,
		Feature:  req.Feature,
		Type:     req.Type,
		Password: req.Password,
	})
	if err != nil {
		if errNo := errno.ConvertErr(err); errNo.ErrorCode == http.StatusUnauthorized {
			c.String(consts.StatusOK, constants.UrlCustomErrorMsg)
			return
		}
		pack.RespError(c, err)
		return
	}

	c.JSON(consts.StatusOK, consts.StatusOK)
}

// GetUploadParams .
// @router /api/v1/url/api/uploadparams [POST]
func GetUploadParams(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.UploadParamsRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.UploadParams: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	resp := new(api.UploadParamsResponse)
	policy, auth, err := rpc.UploadParamsRPC(ctx, &url.UploadParamsRequest{Password: req.Password})
	if err != nil {
		if errNo := errno.ConvertErr(err); errNo.ErrorCode == http.StatusUnauthorized {
			c.String(consts.StatusOK, constants.UrlCustomErrorMsg)
			return
		}
		pack.RespError(c, err)
		return
	}
	resp.Policy = policy
	resp.Authorization = auth
	c.JSON(consts.StatusOK, resp)
}

// GetReleaseVersionModify .
// @router /api/v1/url/version.json [GET]
func GetReleaseVersionModify(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetReleaseVersionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.GetReleaseVersionModify: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	resp := new(api.GetReleaseVersionResponse)

	rpcResp, err := rpc.GetReleaseVersionRPC(ctx, &url.GetReleaseVersionRequest{})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.Version = rpcResp.Version
	resp.URL = rpcResp.Url
	resp.Code = rpcResp.Code
	resp.Feature = rpcResp.Feature
	c.JSON(consts.StatusOK, resp)
}

// GetBetaVersionModify .
// @router /api/v1/url/versionbeta.json [GET]
func GetBetaVersionModify(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetBetaVersionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		logger.Errorf("api.GetBetaVersionModify: BindAndValidate error %v", err)
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	resp := new(api.GetBetaVersionResponse)

	rpcResp, err := rpc.GetBetaVersionRPC(ctx, &url.GetBetaVersionRequest{})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.Version = rpcResp.Version
	resp.URL = rpcResp.Url
	resp.Code = rpcResp.Code
	resp.Feature = rpcResp.Feature
	c.JSON(consts.StatusOK, resp)
}
