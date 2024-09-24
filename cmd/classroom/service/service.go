package service

import (
	"context"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"net/http"
)

type ClassroomService struct {
	ctx        context.Context
	Identifier string
	cookies    []*http.Cookie
}

func NewClassroomService(ctx context.Context, identifier string, cookies []string) *ClassroomService {
	return &ClassroomService{
		ctx:        ctx,
		Identifier: identifier,
		cookies:    utils.ParseCookies(cookies),
	}
}
