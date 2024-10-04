package service

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *UserService) GetLoginData(req *user.GetLoginDataRequest) (string, []string, error) {
	id, rawCookies := jwch.NewStudent().WithUser(req.Id, req.Password).GetIdentifierAndCookies()
	return id, utils.ParseCookiesToString(rawCookies), nil
}
