package service

import (
	"context"
)

type ScreenService struct {
	ctx context.Context
}

func NewScreenService(ctx context.Context) *ScreenService {
	return &ScreenService{ctx: ctx}
}
