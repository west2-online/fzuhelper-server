package service

import (
	"context"
	"net/http"
)

type ClassroomService struct {
	ctx        context.Context
	Identifier string
	cookies    []*http.Cookie
}

func NewClassroomService(ctx context.Context, identifier string, cookies []*http.Cookie) *ClassroomService {
	return &ClassroomService{
		ctx:        ctx,
		Identifier: identifier,
		cookies:    cookies,
	}
}
