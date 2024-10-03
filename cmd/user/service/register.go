package service

import (
	"github.com/west2-online/fzuhelper-server/cmd/user/dal/db"
	"github.com/west2-online/fzuhelper-server/cmd/user/pack/pwd"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
)

func (s *UserService) Register(req *user.RegisterRequest) (*db.User, error) {
	PwdDigest := pwd.SetPassword(req.Password)
	userModel := &db.User{
		Account:  req.Account,
		Name:     req.Name,
		Password: PwdDigest,
	}
	return db.Register(s.ctx, userModel)
}
