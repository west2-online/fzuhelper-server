package main

import (
	"context"
	oa "github.com/west2-online/fzuhelper-server/kitex_gen/oa"
)

// OAServiceImpl implements the last service interface defined in the IDL.
type OAServiceImpl struct{}

// CreateFeedback implements the OAServiceImpl interface.
func (s *OAServiceImpl) CreateFeedback(ctx context.Context, request *oa.CreateFeedbackRequest) (resp *oa.CreateFeedbackResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFeedback implements the OAServiceImpl interface.
func (s *OAServiceImpl) GetFeedback(ctx context.Context, request *oa.GetFeedbackRequest) (resp *oa.GetFeedbackResponse, err error) {
	// TODO: Your code here...
	return
}
