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

	"github.com/west2-online/fzuhelper-server/internal/common/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
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
		logger.Infof("Common.GetCSS: %v", err)
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
		logger.Infof("Common.GetHtml: %v", err)
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
		logger.Infof("Common.GetUserAgreement: %v", err)
		return resp, nil
	}
	resp.UserAgreement = *agreement
	return resp, nil
}
