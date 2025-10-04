package main

import (
	"context"
	api "github.com/west2-online/fzuhelper-server/kitex_gen/api"
)

// FeedbackServiceImpl implements the last service interface defined in the IDL.
type FeedbackServiceImpl struct{}

// CreateFeedback implements the FeedbackServiceImpl interface.
func (s *FeedbackServiceImpl) CreateFeedback(ctx context.Context, request *api.CreateFeedbackRequest) (resp *api.CreateFeedbackResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFeedback implements the FeedbackServiceImpl interface.
func (s *FeedbackServiceImpl) GetFeedback(ctx context.Context, request *api.GetFeedbackRequest) (resp *api.GetFeedbackResponse, err error) {
	// TODO: Your code here...
	return
}
