package service

import (
	"context"
	"github.com/west2-online/fzuhelper-server/cmd/user/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

type UserService struct {
	ctx context.Context
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{ctx: ctx}
}

func BuildUserResp(dbUser *db.User) *model.User {
	return &model.User{
		Id:      dbUser.ID,
		Name:    dbUser.Name,
		Account: dbUser.Account,
	}
}
