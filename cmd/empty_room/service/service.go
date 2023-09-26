package service

import (
	"context"
)

type EmptyRoomService struct {
	ctx context.Context
}

func NewEmptyRoomService(ctx context.Context) *EmptyRoomService {
	return &EmptyRoomService{ctx: ctx}
}
