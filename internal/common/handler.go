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
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
)

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
	css, err := service.NewCommonService(ctx, s.ClientSet, nil).GetCSS()
	if err != nil {
		resp.Css = fmt.Append(nil, err)
		return resp, nil
	}
	resp.Css = *css
	return resp, nil
}

// GetHtml implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetHtml(ctx context.Context, req *common.GetHtmlRequest) (resp *common.GetHtmlResponse, err error) {
	resp = new(common.GetHtmlResponse)
	html, err := service.NewCommonService(ctx, s.ClientSet, nil).GetHtml()
	if err != nil {
		resp.Html = fmt.Append(nil, err)
		return resp, nil
	}
	resp.Html = *html
	return resp, nil
}

// GetUserAgreement implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetUserAgreement(ctx context.Context, req *common.GetUserAgreementRequest) (resp *common.GetUserAgreementResponse, err error) {
	resp = new(common.GetUserAgreementResponse)
	agreement, err := service.NewCommonService(ctx, s.ClientSet, nil).GetUserAgreement()
	if err != nil {
		resp.UserAgreement = fmt.Append(nil, err)
		return resp, nil
	}
	resp.UserAgreement = *agreement
	return resp, nil
}

// GetTermsList implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTermsList(ctx context.Context, req *common.TermListRequest) (resp *common.TermListResponse, err error) {
	resp = new(common.TermListResponse)
	termList, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetTermList()
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.TermLists = pack.BuildTermsList(termList)
	return resp, nil
}

// GetTerm implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTerm(ctx context.Context, req *common.TermRequest) (resp *common.TermResponse, err error) {
	resp = new(common.TermResponse)
	term, err := service.NewCommonService(ctx, s.ClientSet, s.taskQueue).GetTerm(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.TermInfo = pack.BuildTermInfo(term)
	return resp, nil
}

func (s *CommonServiceImpl) GetNotices(ctx context.Context, req *common.NoticeRequest) (resp *common.NoticeResponse, err error) {
	resp = new(common.NoticeResponse)
	notices, total, err := service.NewCommonService(ctx, s.ClientSet, nil).GetNotice(int(req.PageNum))
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Notices = pack.BuildNoticeList(notices)
	resp.Total = int64(total)
	return resp, nil
}

func (s *CommonServiceImpl) GetContributorInfo(ctx context.Context, _ *common.GetContributorInfoRequest) (resp *common.GetContributorInfoResponse, err error) {
	resp = new(common.GetContributorInfoResponse)
	contributors, err := service.NewCommonService(ctx, s.ClientSet, nil).GetContributorInfo()
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.FzuhelperApp = contributors[constants.ContributorFzuhelperAppKey]
	resp.FzuhelperServer = contributors[constants.ContributorFzuhelperServerKey]
	resp.Jwch = contributors[constants.ContributorJwchKey]
	resp.Yjsy = contributors[constants.ContributorYJSYKey]
	return resp, nil
}

func (s *CommonServiceImpl) GetToolboxConfig(ctx context.Context, req *common.GetToolboxConfigRequest) (resp *common.GetToolboxConfigResponse, err error) {
	resp = new(common.GetToolboxConfigResponse)
	Configs, err := service.NewCommonService(ctx, s.ClientSet, nil).GetToolboxConfig(ctx, req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Config = pack.BuildToolboxConfigList(Configs)
	return resp, nil
}

func (s *CommonServiceImpl) PutToolboxConfig(ctx context.Context, req *common.PutToolboxConfigRequest) (resp *common.PutToolboxConfigResponse, err error) {
	resp = new(common.PutToolboxConfigResponse)
	config, err := service.NewCommonService(ctx, s.ClientSet, nil).PutToolboxConfig(ctx, req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.ConfigId = &config.Id
	return resp, nil
}
