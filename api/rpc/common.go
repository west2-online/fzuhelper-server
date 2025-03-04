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

package rpc

import (
	"context"

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitCommonRPC() {
	client, err := client.InitCommonRPC()
	if err != nil {
		logger.Fatalf("api.rpc.common InitCommonRPC failed, err  %v", err)
	}
	commonClient = *client
}

func GetCSSRPC(ctx context.Context, req *common.GetCSSRequest) (*[]byte, error) {
	resp, err := commonClient.GetCSS(ctx, req)
	if err != nil {
		logger.Errorf("GetCSSRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if len(resp.Css) < 1 {
		return nil, errno.InternalServiceError
	}
	return &resp.Css, nil
}

func GetHtmlRPC(ctx context.Context, req *common.GetHtmlRequest) (*[]byte, error) {
	resp, err := commonClient.GetHtml(ctx, req)
	if err != nil {
		logger.Errorf("GetHtmlRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if len(resp.Html) < 1 {
		return nil, errno.InternalServiceError
	}
	return &resp.Html, nil
}

func GetUserAgreementRPC(ctx context.Context, req *common.GetUserAgreementRequest) (*[]byte, error) {
	resp, err := commonClient.GetUserAgreement(ctx, req)
	if err != nil {
		logger.Errorf("GetUserAgreementRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if len(resp.UserAgreement) < 1 {
		return nil, errno.InternalServiceError
	}
	return &resp.UserAgreement, nil
}

func GetTermsListRPC(ctx context.Context, req *common.TermListRequest) (*model.TermList, error) {
	resp, err := commonClient.GetTermsList(ctx, req)
	if err != nil {
		logger.Errorf("GetTermsListRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}

	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, err
	}

	return resp.TermLists, nil
}

func GetTermRPC(ctx context.Context, req *common.TermRequest) (*model.TermInfo, error) {
	resp, err := commonClient.GetTerm(ctx, req)
	if err != nil {
		logger.Errorf("GetTermRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, err
	}

	return resp.TermInfo, nil
}

func GetNoticesRPC(ctx context.Context, req *common.NoticeRequest) ([]*model.NoticeInfo, int64, error) {
	resp, err := commonClient.GetNotices(ctx, req)
	if err != nil {
		logger.Errorf("GetNoticesRPC: RPC called failed: %v", err.Error())
		return nil, 0, errno.InternalServiceError.WithMessage(err.Error())
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, 0, err
	}
	return resp.Notices, resp.Total, nil
}

func GetContributorRPC(ctx context.Context, req *common.GetContributorInfoRequest) (*common.GetContributorInfoResponse, error) {
	resp, err := commonClient.GetContributorInfo(ctx, req)
	if err != nil {
		logger.Errorf("GetContributorRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, err
	}
	return resp, nil
}
