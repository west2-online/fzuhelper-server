package main

import (
	"context"
	common "github.com/west2-online/fzuhelper-server/kitex_gen/common"
)

// CommonServiceImpl implements the last service interface defined in the IDL.
type CommonServiceImpl struct{}

// GetCSS implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetCSS(ctx context.Context, req *common.GetCSSRequest) (resp *common.GetCSSResponse, err error) {
	// TODO: Your code here...
	return
}

// GetHtml implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetHtml(ctx context.Context, req *common.GetHtmlRequest) (resp *common.GetHtmlResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserAgreement implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetUserAgreement(ctx context.Context, req *common.GetUserAgreementRequest) (resp *common.GetUserAgreementResponse, err error) {
	// TODO: Your code here...
	return
}

// GetTermsList implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTermsList(ctx context.Context, req *common.TermListRequest) (resp *common.TermListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetTerm implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTerm(ctx context.Context, req *common.TermRequest) (resp *common.TermResponse, err error) {
	// TODO: Your code here...
	return
}

// GetNotices implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetNotices(ctx context.Context, req *common.NoticeRequest) (resp *common.NoticeResponse, err error) {
	// TODO: Your code here...
	return
}

// GetContributorInfo implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetContributorInfo(ctx context.Context, req *common.GetContributorInfoRequest) (resp *common.GetContributorInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// GetToolboxConfig implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetToolboxConfig(ctx context.Context, req *common.GetToolboxConfigRequest) (resp *common.GetToolboxConfigResponse, err error) {
	// TODO: Your code here...
	return
}

// PutToolboxConfig implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) PutToolboxConfig(ctx context.Context, req *common.PutToolboxConfigRequest) (resp *common.PutToolboxConfigResponse, err error) {
	// TODO: Your code here...
	return
}
