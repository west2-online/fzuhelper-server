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
		ReportId:       req.GetReportId(),
		StuId:          req.GetStuId(),
		Name:           req.GetName(),
		College:        req.GetCollege(),
		ContactPhone:   req.GetContactPhone(),
		ContactQQ:      req.GetContactQQ(),
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

	resp, err := rpc.GetFeedbackRPC(ctx, &oa.GetFeedbackRequest{ReportId: req.ReportId})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.RespData(c, resp)
}
