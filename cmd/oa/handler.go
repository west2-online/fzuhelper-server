package main

import (
	"context"

	oa "github.com/west2-online/fzuhelper-server/kitex_gen/oa"
)

// OaServiceImpl implements the last service interface defined in the IDL.
type OaServiceImpl struct{}

// CreateFeedback implements the OaServiceImpl interface.
func (s *OaServiceImpl) CreateFeedback(ctx context.Context, request *oa.CreateFeedbackRequest) (resp *oa.CreateFeedbackResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFeedback implements the OaServiceImpl interface.
func (s *OaServiceImpl) GetFeedback(ctx context.Context, request *oa.GetFeedbackRequest) (resp *oa.GetFeedbackResponse, err error) {
	// TODO: Your code here...
	return
}
