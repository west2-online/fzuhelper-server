package api

import (
	"context"
	"errors"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

var loginDataKey *model.LoginData

func GetLoginData(ctx context.Context) (*model.LoginData, error) {
	user, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("获取Header错误")
	}
	return user, nil
}

func NewContext(ctx context.Context, value *model.LoginData) context.Context {
	return context.WithValue(ctx, loginDataKey, value)
}

func FromContext(ctx context.Context) (*model.LoginData, bool) {
	u, ok := ctx.Value(loginDataKey).(*model.LoginData)
	return u, ok
}
