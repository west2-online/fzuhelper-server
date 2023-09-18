package service

import "context"

type TemplateService struct {
	ctx context.Context
}

// NewTemplateService new TemplateService
func NewTemplateService(ctx context.Context) *TemplateService {
	return &TemplateService{ctx: ctx}
}
