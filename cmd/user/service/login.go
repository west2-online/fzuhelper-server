package service

import (
	"github.com/west2-online/fzuhelper-server/cmd/user/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
)

func (s *UserService) Login(req *user.LoginRequest /*, stTracer opentracing.Tracer, parentSpan opentracing.Span*/) (*db.User, error) {
	userModel := &db.User{
		Account:  req.Account,
		Password: req.Password,
	}

	userResp, err := db.Login(s.ctx, userModel)
	if err != nil {
		return nil, err
	}

	return userResp, nil
}
