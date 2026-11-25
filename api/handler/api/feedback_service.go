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
	"io"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	api "github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/model/model"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/api/rpc"
	oa "github.com/west2-online/fzuhelper-server/kitex_gen/oa"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func readFeedbackUpload(c *app.RequestContext, maxSize int64) ([]byte, error) {
	header, err := c.FormFile("file")
	if err != nil || header.Size <= 0 {
		return nil, errno.NewErrNo(errno.ParamFileNotExistCode, "上传文件不能为空")
	}
	if header.Size > maxSize {
		return nil, errno.NewErrNo(errno.ParamRangeCode, "上传文件超过大小限制")
	}

	file, err := header.Open()
	if err != nil {
		return nil, errno.NewErrNo(errno.ParamFileReadErrorCode, "读取上传文件失败")
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxSize+1))
	if err != nil {
		return nil, errno.NewErrNo(errno.ParamFileReadErrorCode, "读取上传文件失败")
	}
	if len(data) == 0 {
		return nil, errno.NewErrNo(errno.ParamFileNotExistCode, "上传文件不能为空")
	}
	if int64(len(data)) > maxSize {
		return nil, errno.NewErrNo(errno.ParamRangeCode, "上传文件超过大小限制")
	}
	return data, nil
}

// CreateFeedback .
// @router /api/v1/feedback/create [POST]
func CreateFeedback(ctx context.Context, c *app.RequestContext) {
	var req api.CreateFeedbackRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.RespError(c, errno.ParamError.WithError(err))
		return
	}
	stuID, ok := utils.GetStuID(c)
	if !ok {
		pack.RespError(c, errno.AuthMissing)
		return
	}

	resp := new(api.CreateFeedbackResponse)
	reportID, err := rpc.CreateFeedbackRPC(ctx, &oa.CreateFeedbackRequest{
		StuId:          stuID,
		Name:           req.GetName(),
		College:        req.GetCollege(),
		ContactPhone:   req.ContactPhone,
		ContactQq:      req.ContactQq,
		ContactEmail:   req.ContactEmail,
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
	stuID, ok := utils.GetStuID(c)
	if !ok {
		pack.RespError(c, errno.AuthMissing)
		return
	}

	resp := new(api.GetFeedbackByIDResponse)
	data, err := rpc.GetFeedbackByIDRPC(ctx, &oa.GetFeedbackByIDRequest{
		ReportId: req.ReportID,
		StuId:    stuID,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
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
	stuID, ok := utils.GetStuID(c)
	if !ok {
		pack.RespError(c, errno.AuthMissing)
		return
	}

	resp := new(api.GetListFeedbackResponse)
	data, pageToken, err := rpc.ListFeedbackRPC(ctx, &oa.GetListFeedbackRequest{
		StuId:     stuID,
		Limit:     req.Limit,
		PageToken: req.PageToken,
		OrderDesc: req.OrderDesc,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.Data = pack.BuildFeedbackList(data)
	resp.PageToken = pageToken
	pack.RespData(c, resp)
}

// UploadFeedbackScreenshot .
// @router /api/v1/feedback/upload [POST]
func UploadFeedbackScreenshot(ctx context.Context, c *app.RequestContext) {
	file, err := readFeedbackUpload(c, constants.FeedbackScreenshotMaxSize)
	if err != nil {
		pack.RespError(c, err)
		return
	}
	mimeType := http.DetectContentType(file)
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		pack.RespError(c, errno.NewErrNo(errno.ParamFormatCode, "反馈截图仅支持 JPEG 和 PNG"))
		return
	}

	resp := new(api.UploadFeedbackScreenshotResponse)
	url, err := rpc.UploadFeedbackScreenshotRPC(ctx, &oa.UploadFeedbackScreenshotRequest{File: file})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.URL = url
	pack.RespData(c, resp)
}

// UploadFeedbackLog .
// @router /api/v1/feedback/upload-log [POST]
func UploadFeedbackLog(ctx context.Context, c *app.RequestContext) {
	file, err := readFeedbackUpload(c, constants.FeedbackLogMaxSize)
	if err != nil {
		pack.RespError(c, err)
		return
	}
	if !utils.ValidateJSONLines(file) {
		pack.RespError(c, errno.NewErrNo(errno.ParamFormatCode, "反馈日志必须是 UTF-8 JSON Lines"))
		return
	}

	resp := new(api.UploadFeedbackLogResponse)
	url, err := rpc.UploadFeedbackLogRPC(ctx, &oa.UploadFeedbackLogRequest{File: file})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	resp.URL = url
	pack.RespData(c, resp)
}
