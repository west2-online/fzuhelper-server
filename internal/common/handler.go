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

package common

import (
	"context"
	"fmt"

	"github.com/west2-online/fzuhelper-server/internal/common/pack"
	"github.com/west2-online/fzuhelper-server/internal/common/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/db/toolbox"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/singleflight"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/jwch"
)

type termResult struct {
	success bool
	events  *jwch.CalTermEvents
	err     error
}

type noticeResult struct {
	list  []model.Notice
	total int
}

// CommonServiceImpl implements the last service interface defined in the IDL.
type CommonServiceImpl struct {
	ClientSet *base.ClientSet
	taskQueue taskqueue.TaskQueue
}

func NewCommonService(clientSet *base.ClientSet, taskQueue taskqueue.TaskQueue) *CommonServiceImpl {
	return &CommonServiceImpl{
		ClientSet: clientSet,
		taskQueue: taskQueue,
	}
}

// GetCSS implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetCSS(ctx context.Context, req *common.GetCSSRequest) (resp *common.GetCSSResponse, err error) {
	resp = new(common.GetCSSResponse)
	css, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetCSS()
	if err != nil {
		logger.WithCtx(ctx).Infof("Common.GetCSS: %v", err)
		return resp, nil
	}
	resp.Css = *css
	return resp, nil
}

// GetHtml implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetHtml(ctx context.Context, req *common.GetHtmlRequest) (resp *common.GetHtmlResponse, err error) {
	resp = new(common.GetHtmlResponse)
	html, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetHtml()
	if err != nil {
		logger.WithCtx(ctx).Infof("Common.GetHtml: %v", err)
		return resp, nil
	}
	resp.Html = *html
	return resp, nil
}

// GetUserAgreement implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetUserAgreement(ctx context.Context,
	req *common.GetUserAgreementRequest,
) (resp *common.GetUserAgreementResponse, err error) {
	resp = new(common.GetUserAgreementResponse)
	agreement, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetUserAgreement()
	if err != nil {
		logger.WithCtx(ctx).Infof("Common.GetUserAgreement: %v", err)
		return resp, nil
	}
	resp.UserAgreement = *agreement
	return resp, nil
}

// GetTermsList implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTermsList(ctx context.Context, req *common.TermListRequest) (resp *common.TermListResponse, err error) {
	resp = common.NewTermListResponse()

	res, err := singleflight.Do(constants.SingleflightTermListKey, func() (*jwch.SchoolCalendar, error) {
		return service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetTermList()
	})
	if err != nil {
		resp.Base = base.BuildBaseResp(fmt.Errorf("Common.GetTermsList: get terms list failed: %w", err))
		return resp, nil
	}

	resp.Base = base.BuildBaseResp(nil)
	resp.TermLists = pack.BuildTermsList(res)
	return resp, err
}

// GetTerm implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTerm(ctx context.Context, req *common.TermRequest) (resp *common.TermResponse, err error) {
	resp = common.NewTermResponse()

	key := singleflight.Key(constants.SingleflightTermPrefix, req.Term)
	result, err := singleflight.Do(key, func() (termResult, error) {
		success, events, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetTerm(req)
		if err != nil && !success {
			return termResult{}, err
		}
		return termResult{success: success, events: events, err: err}, nil
	})
	if err != nil {
		base.LogError(fmt.Errorf("Common.GetTerm: get term info failed: %w", err))
	}
	if result.err != nil {
		base.LogError(fmt.Errorf("Common.GetTerm: get term info partially failed: %w", result.err))
	}

	if !result.success {
		resp.Base = base.BuildBaseResp(fmt.Errorf("Common.GetTerm: get term failed: %w", err))
		return resp, nil
	}

	resp.Base = base.BuildBaseResp(nil)
	resp.TermInfo = pack.BuildTermInfo(result.events)
	return resp, err
}

func (s *CommonServiceImpl) GetNotices(ctx context.Context, req *common.NoticeRequest) (resp *common.NoticeResponse, err error) {
	resp = new(common.NoticeResponse)
	key := singleflight.Key(constants.SingleflightNoticePrefix, req.PageNum)
	result, err := singleflight.Do(key, func() (noticeResult, error) {
		list, total, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetNotice(int(req.PageNum))
		if err != nil {
			return noticeResult{}, err
		}
		return noticeResult{list: list, total: total}, nil
	})
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Notices = pack.BuildNoticeList(result.list)
	resp.Total = int64(result.total)
	return resp, err
}

func (s *CommonServiceImpl) GetContributorInfo(ctx context.Context,
	_ *common.GetContributorInfoRequest,
) (resp *common.GetContributorInfoResponse, err error) {
	resp = new(common.GetContributorInfoResponse)

	res, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetContributorInfo()
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.FzuhelperApp = res[constants.ContributorFzuhelperAppKey]
	resp.FzuhelperServer = res[constants.ContributorFzuhelperServerKey]
	resp.Jwch = res[constants.ContributorJwchKey]
	resp.Yjsy = res[constants.ContributorYJSYKey]
	return resp, nil
}

func (s *CommonServiceImpl) GetToolboxConfig(ctx context.Context,
	req *common.GetToolboxConfigRequest,
) (r *common.GetToolboxConfigResponse, err error) {
	r = new(common.GetToolboxConfigResponse)

	// 获取请求参数，如果为空则使用默认值
	studentID := ""
	if req.StudentId != nil {
		studentID = *req.StudentId
	}

	platform := ""
	if req.Platform != nil {
		platform = *req.Platform
	}

	version := int64(0)
	if req.Version != nil {
		version = *req.Version
	}

	// 调用service获取配置
	dbConfigs, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetToolboxConfig(ctx, studentID, platform, version)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}

	r.Base = base.BuildSuccessResp()
	r.Config = pack.BuildToolboxConfigList(dbConfigs)
	return r, nil
}

func (s *CommonServiceImpl) CreateToolboxConfig(ctx context.Context,
	req *common.CreateToolboxConfigRequest,
) (r *common.CreateToolboxConfigResponse, err error) {
	r = new(common.CreateToolboxConfigResponse)
	config, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).CreateToolboxConfig(
		ctx,
		req.Secret,
		&model.ToolboxConfig{
			ToolID:    req.ToolId,
			Visible:   req.Visible,
			Name:      req.Name,
			Icon:      req.Icon,
			Type:      req.Type,
			Message:   req.Message,
			Extra:     req.Extra,
			StudentID: req.StudentId,
			Platform:  req.Platform,
			Version:   req.Version,
		},
	)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}
	r.Base = base.BuildSuccessResp()
	r.Config = pack.BuildToolboxConfigDetail(config)
	return r, nil
}

func (s *CommonServiceImpl) ListToolboxConfigs(ctx context.Context,
	req *common.ListToolboxConfigsRequest,
) (r *common.ListToolboxConfigsResponse, err error) {
	r = new(common.ListToolboxConfigsResponse)
	configs, total, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).ListToolboxConfigs(
		ctx,
		req.Secret,
		req.GetPageNum(),
		req.GetPageSize(),
		toolbox.ListToolboxConfigsFilter{
			ToolID:     req.ToolId,
			StudentID:  req.StudentId,
			Platform:   req.Platform,
			MinVersion: req.Version,
		},
	)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}
	r.Base = base.BuildSuccessResp()
	r.Config = pack.BuildToolboxConfigDetailList(configs)
	r.Total = total
	return r, nil
}

func (s *CommonServiceImpl) GetToolboxConfigByID(ctx context.Context,
	req *common.GetToolboxConfigByIDRequest,
) (r *common.GetToolboxConfigByIDResponse, err error) {
	r = new(common.GetToolboxConfigByIDResponse)
	config, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetToolboxConfigByID(ctx, req.Secret, req.ConfigId)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}
	r.Base = base.BuildSuccessResp()
	r.Config = pack.BuildToolboxConfigDetail(config)
	return r, nil
}

func (s *CommonServiceImpl) UpdateToolboxConfig(ctx context.Context,
	req *common.UpdateToolboxConfigRequest,
) (r *common.UpdateToolboxConfigResponse, err error) {
	r = new(common.UpdateToolboxConfigResponse)
	config, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).UpdateToolboxConfig(
		ctx,
		req.Secret,
		req.ConfigId,
		&model.ToolboxConfig{
			ToolID:    req.ToolId,
			Visible:   req.Visible,
			Name:      req.Name,
			Icon:      req.Icon,
			Type:      req.Type,
			Message:   req.Message,
			Extra:     req.Extra,
			StudentID: req.StudentId,
			Platform:  req.Platform,
			Version:   req.Version,
		},
	)
	if err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}
	r.Base = base.BuildSuccessResp()
	r.Config = pack.BuildToolboxConfigDetail(config)
	return r, nil
}

func (s *CommonServiceImpl) DeleteToolboxConfig(ctx context.Context,
	req *common.DeleteToolboxConfigRequest,
) (r *common.DeleteToolboxConfigResponse, err error) {
	r = new(common.DeleteToolboxConfigResponse)
	err = service.NewCommonService(ctx, s.ClientSet, s.taskQueue).DeleteToolboxConfig(ctx, req.Secret, req.ConfigId)
	r.Base = base.BuildBaseResp(err)
	return r, nil
}

func (s *CommonServiceImpl) TracePing(ctx context.Context, req *common.TracePingRequest) (resp *common.TracePingResponse, err error) {
	// log with trace context
	logger.WithCtx(ctx).Info("RPC trace ping request received")

	resp = new(common.TracePingResponse)
	resp.Base = base.BuildSuccessResp()
	resp.Message = "pong"
	return resp, nil
}

func (s *CommonServiceImpl) GetSignedLocationApiUrl(
	ctx context.Context,
	req *common.GetSignedLocationApiUrlRequest,
) (resp *common.GetSignedLocationApiUrlResponse, err error) {
	resp = new(common.GetSignedLocationApiUrlResponse)

	signedURL, headers, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetSignedApiUrl(req.Location)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = base.BuildSuccessResp()
	resp.SignedUrl = signedURL
	resp.Headers = headers

	return resp, nil
}
