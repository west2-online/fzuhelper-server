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

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	api "github.com/west2-online/fzuhelper-server/api/model/api"
)

// Login .
// @router /api/v2/url/login [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.LoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.LoginResponse)

	c.JSON(consts.StatusOK, resp)
}

// UploadVersion .
// @router /api/v2/url/api/upload [POST]
func UploadVersion(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.UploadRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.UploadResponse)

	c.JSON(consts.StatusOK, resp)
}

// UploadParams .
// @router /api/v2/url/api/uploadparams [POST]
func UploadParams(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.UploadParamsRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.UploadResponse)

	c.JSON(consts.StatusOK, resp)
}

// DownloadReleaseApk .
// @router /api/v2/url/release.apk [GET]
func DownloadReleaseApk(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DownloadReleaseApkRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.DownloadReleaseApkResponse)

	c.JSON(consts.StatusOK, resp)
}

// DownloadBetaApk .
// @router /api/v2/url/beta.apk [GET]
func DownloadBetaApk(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DownloadBetaApkRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.DownloadBetaApkResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetReleaseVersion .
// @router /api/v2/url/version.json [GET]
func GetReleaseVersion(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetReleaseVersionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.GetReleaseVersionResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetBetaVersion .
// @router /api/v2/url/versionbeta.json [GET]
func GetBetaVersion(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetBetaVersionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.GetBetaVersionResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetSetting .
// @router /api/v2/url/settings.php [GET]
func GetSetting(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetSettingRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.GetSettingResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetTest .
// @router /api/v2/url/test [POST]
func GetTest(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetSettingRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.GetTestResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetCloud .
// @router /api/v2/url/getcloud [GET]
func GetCloud(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetCloudRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.GetCloudResponse)

	c.JSON(consts.StatusOK, resp)
}

// SetCloud .
// @router /api/v2/url/setcloud [GET]
func SetCloud(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.SetCloudRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.SetCloudResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetDump .
// @router /api/v2/url/dump [GET]
func GetDump(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetDumpRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.GetDumpResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetCSS .
// @router /api/v2/url/onekey/FZUHelper.css [GET]
func GetCSS(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetCSSRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.GetCSSResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetHtml .
// @router /api/v2/url/onekey/FZUHelper.html [GET]
func GetHtml(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetHtmlRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.GetHtmlResponse)

	c.JSON(consts.StatusOK, resp)
}

// GetUserAgreement .
// @router /api/v2/url/onekey/UserAgreement.html [GET]
func GetUserAgreement(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetUserAgreementRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(api.GetUserAgreementResponse)

	c.JSON(consts.StatusOK, resp)
}
