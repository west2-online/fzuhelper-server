package service

import (
	"context"
)

type ClassroomService struct {
	ctx context.Context
}

func NewClassroomService(ctx context.Context) *ClassroomService {
	return &ClassroomService{
		ctx: ctx,
	}
}
