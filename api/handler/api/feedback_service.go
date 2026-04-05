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
	"github.com/west2-online/fzuhelper-server/api/model/model"
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

	reportID, err := rpc.CreateFeedbackRPC(ctx, &oa.CreateFeedbackRequest{
		StuId:          req.StuID,
		Name:           req.Name,
		College:        req.College,
		ContactPhone:   req.ContactPhone,
		ContactQq:      req.ContactQq,
		ContactEmail:   req.ContactEmail,
		NetworkEnv:     req.NetworkEnv,
		IsOnCampus:     req.IsOnCampus,
		OsName:         req.OsName,
		OsVersion:      req.OsVersion,
		Manufacturer:   req.Manufacturer,
		DeviceModel:    req.DeviceModel,
		ProblemDesc:    req.ProblemDesc,
		Screenshots:    req.Screenshots,
		AppVersion:     req.AppVersion,
		VersionHistory: req.VersionHistory,
		NetworkTraces:  req.NetworkTraces,
		Events:         req.Events,
		UserSettings:   req.UserSettings,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp := new(api.CreateFeedbackResponse)
	resp.ReportID = reportID
	pack.RespData(c, resp)
}

// GetFeedbackByID .
// @router /api/v1/feedbacks/detail [GET]
func GetFeedbackByID(ctx context.Context, c *app.RequestContext) {
	var req api.GetFeedbackByIDRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	data, err := rpc.GetFeedbackByIdRPC(ctx, &oa.GetFeedbackByIDRequest{
		ReportId: req.ReportID,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp := new(api.FeedbackDetailResponse)
	resp.Data = &model.Feedback{
		ReportID:       data.ReportId,
		StuID:          data.StuId,
		Name:           data.Name,
		College:        data.College,
		ContactPhone:   data.ContactPhone,
		ContactQq:      data.ContactQq,
		ContactEmail:   data.ContactEmail,
		NetworkEnv:     data.NetworkEnv,
		IsOnCampus:     data.IsOnCampus,
		OsName:         data.OsName,
		OsVersion:      data.OsVersion,
		Manufacturer:   data.Manufacturer,
		DeviceModel:    data.DeviceModel,
		ProblemDesc:    data.ProblemDesc,
		Screenshots:    data.Screenshots,
		AppVersion:     data.AppVersion,
		VersionHistory: data.VersionHistory,
		NetworkTraces:  data.NetworkTraces,
		Events:         data.Events,
		UserSettings:   data.UserSettings,
	}
	pack.RespData(c, resp)
}

// ListFeedback .
// @router /api/v1/feedbacks/get/list [GET]
func ListFeedback(ctx context.Context, c *app.RequestContext) {
	var req api.GetListFeedbackRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}

	data, pageToken, err := rpc.GetFeedbackListRPC(ctx, &oa.GetListFeedbackRequest{
		StuId:       req.StuID,
		Name:        req.Name,
		NetworkEnv:  req.NetworkEnv,
		IsOnCampus:  req.IsOnCampus,
		OsName:      req.OsName,
		ProblemDesc: req.ProblemDesc,
		AppVersion:  req.AppVersion,
		BeginTimeMs: req.BeginTimeMs,
		EndTimeMs:   req.EndTimeMs,
		Limit:       req.Limit,
		PageToken:   req.PageToken,
		OrderDesc:   req.OrderDesc,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp := new(api.GetListFeedbackResponse)
	resp.Data = pack.BuildFeedbackList(data)
	resp.PageToken = pageToken
	pack.RespData(c, resp)
}
