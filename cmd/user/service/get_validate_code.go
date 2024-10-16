package service

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/jwch"
)

func (s *UserService) GetValidateCode(req *user.GetValidateCodeRequest) (string, error) {
	return jwch.GetValidateCode(req.Image)
}
