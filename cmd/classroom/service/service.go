package service

import (
	"context"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/jwch"
	"net/http"
)

type ClassroomService struct {
	ctx        context.Context
	Identifier string
	cookies    []*http.Cookie
}

func NewClassroomServiceInDefault(ctx context.Context) *ClassroomService {
	id, cookies := jwch.NewStudent().WithUser(constants.DefaultAccount, constants.DefaultPassword).GetIdentifierAndCookies()
	return &ClassroomService{
		ctx:        ctx,
		Identifier: id,
		cookies:    cookies,
	}
}

func NewClassroomService(ctx context.Context, identifier string, cookies []*http.Cookie) *ClassroomService {
	return &ClassroomService{
		ctx:        ctx,
		Identifier: identifier,
		cookies:    cookies,
	}
}
