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
)

// CommonServiceImpl implements the last service interface defined in the IDL.
type CommonServiceImpl struct {
	ClientSet *base.ClientSet
}

func NewCommonService(clientSet *base.ClientSet) *CommonServiceImpl {
	return &CommonServiceImpl{
		ClientSet: clientSet,
	}
}

// GetCSS implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetCSS(ctx context.Context, req *common.GetCSSRequest) (resp *common.GetCSSResponse, err error) {
	resp = new(common.GetCSSResponse)
	css, err := service.NewCommonService(ctx, s.ClientSet).GetCSS()
	if err != nil {
		resp.Css = fmt.Appendf(nil, "Common.GetCSS: %v", err)
		return resp, nil
	}
	resp.Css = *css
	return resp, nil
}

// GetHtml implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetHtml(ctx context.Context, req *common.GetHtmlRequest) (resp *common.GetHtmlResponse, err error) {
	resp = new(common.GetHtmlResponse)
	html, err := service.NewCommonService(ctx, s.ClientSet).GetHtml()
	if err != nil {
		resp.Html = fmt.Appendf(nil, "Common.GetHtml: %v", err)
		return resp, nil
	}
	resp.Html = *html
	return resp, nil
}

// GetUserAgreement implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetUserAgreement(ctx context.Context, req *common.GetUserAgreementRequest) (resp *common.GetUserAgreementResponse, err error) {
	resp = new(common.GetUserAgreementResponse)
	agreement, err := service.NewCommonService(ctx, s.ClientSet).GetUserAgreement()
	if err != nil {
		resp.UserAgreement = fmt.Appendf(nil, "Common.GetUserAgreement: %v", err)
		return resp, nil
	}
	resp.UserAgreement = *agreement
	return resp, nil
}

// GetTermsList implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTermsList(ctx context.Context, req *common.TermListRequest) (resp *common.TermListResponse, err error) {
	resp = new(common.TermListResponse)
	res, err := service.NewCommonService(ctx, s.ClientSet).GetTermList()
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.TermLists = pack.BuildTermsList(res)
	return resp, nil
}

// GetTerm implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTerm(ctx context.Context, req *common.TermRequest) (resp *common.TermResponse, err error) {
	resp = new(common.TermResponse)
	success, res, err := service.NewCommonService(ctx, s.ClientSet).GetTerm(req)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	if !success {
		resp.Base = base.BuildBaseResp(fmt.Errorf("Common.GetTerm: get term failed."))
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.TermInfo = pack.BuildTermInfo(res)
	return resp, nil
}

func (s *CommonServiceImpl) GetNotices(ctx context.Context, req *common.NoticeRequest) (resp *common.NoticeResponse, err error) {
	resp = new(common.NoticeResponse)
	res, total, err := service.NewCommonService(ctx, s.ClientSet).GetNotice(int(req.PageNum))
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Notices = pack.BuildNoticeList(res)
	resp.Total = int64(total)
	return resp, nil
}

func (s *CommonServiceImpl) GetContributorInfo(ctx context.Context, _ *common.GetContributorInfoRequest) (resp *common.GetContributorInfoResponse, err error) {
	resp = new(common.GetContributorInfoResponse)
	res, err := service.NewCommonService(ctx, s.ClientSet).GetContributorInfo()
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.FzuhelperApp = res[constants.ContributorFzuhelperAppKey]
	resp.FzuhelperServer = res[constants.ContributorFzuhelperServerKey]
	resp.Jwch = res[constants.ContributorJwchKey]
	resp.Yjsy = res[constants.ContributorYJSYKey]
	return resp, nil
}

func (s *CommonServiceImpl) GetToolboxConfig(ctx context.Context, req *common.GetToolboxConfigRequest) (r *common.GetToolboxConfigResponse, err error) {
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
	dbConfigs, err := service.NewCommonService(ctx, s.ClientSet).GetToolboxConfig(ctx, studentID, platform, version)
	r.Base = base.BuildBaseResp(err)
	if err != nil {
		return r, nil
	}
	r.Config = pack.BuildToolboxConfigList(dbConfigs)
	return r, nil
}

func (s *CommonServiceImpl) PutToolboxConfig(ctx context.Context, req *common.PutToolboxConfigRequest) (r *common.PutToolboxConfigResponse, err error) {
	r = new(common.PutToolboxConfigResponse)

	// 获取请求参数，处理可选字段
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

	// 调用service层创建或更新配置
	config, err := service.NewCommonService(ctx, s.ClientSet).PutToolboxConfig(
		ctx,
		req.Secret,
		req.ToolId,
		studentID,
		platform,
		version,
		req.Visible,
		req.Name,
		req.Icon,
		req.Type,
		req.Message,
		req.Extra,
	)
	r.Base = base.BuildBaseResp(err)
	if err != nil {
		return r, nil
	}
	r.ConfigId = &config.Id
	return r, nil
}
