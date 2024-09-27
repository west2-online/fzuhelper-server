package service

import (
	"context"
	"net/http"
)

type UserService struct {
	ctx        context.Context
	Identifier string
	cookies    []*http.Cookie
}

func NewUserService(ctx context.Context, identifier string, cookies []*http.Cookie) *UserService {
	return &UserService{
		ctx:        ctx,
		Identifier: identifier,
		cookies:    cookies,
	}
}
