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

package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	api "github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/api/rpc"
	oa "github.com/west2-online/fzuhelper-server/kitex_gen/oa"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// CreateFeedback .
// @router /api/v1/feedback/create [POST]
func CreateFeedback(ctx context.Context, c *app.RequestContext) {
	var req api.CreateFeedbackRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	resp := new(api.CreateFeedbackResponse)
	err := rpc.CreateFeedbackRPC(ctx, &oa.CreateFeedbackRequest{
		ReportId:       req.GetReportID(),
		StuId:          req.GetStuID(),
		Name:           req.GetName(),
		College:        req.GetCollege(),
		ContactPhone:   req.GetContactPhone(),
		ContactQq:      req.GetContactQq(),
		ContactEmail:   req.GetContactEmail(),
		NetworkEnv:     req.GetNetworkEnv(),
		IsOnCampus:     req.GetIsOnCampus(),
		OsName:         req.GetOsName(),
		OsVersion:      req.GetOsVersion(),
		Manufacturer:   req.GetManufacturer(),
		DeviceModel:    req.GetDeviceModel(),
		ProblemDesc:    req.GetProblemDesc(),
		Screenshots:    req.GetScreenshots(),
		AppVersion:     req.GetAppVersion(),
		VersionHistory: req.GetVersionHistory(),
		NetworkTraces:  req.GetNetworkTraces(),
		Events:         req.GetEvents(),
		UserSettings:   req.GetUserSettings(),
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.RespData(c, resp)
}

// GetFeedback .
// @router /api/v1/feedback/get [POST]
func GetFeedback(ctx context.Context, c *app.RequestContext) {
	var req api.GetFeedbackRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	resp, err := rpc.GetFeedbackRPC(ctx, &oa.GetFeedbackRequest{ReportId: req.ReportID})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.RespData(c, resp)
}
