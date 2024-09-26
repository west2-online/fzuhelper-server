package api

import (
	"context"
	"github.com/pkg/errors"
)

var loginDataKey *LoginData

func GetLoginData(ctx context.Context) (*LoginData, error) {
	user, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("获取Header错误")
	}
	return user, nil
}

func NewContext(ctx context.Context, value *LoginData) context.Context {
	return context.WithValue(ctx, loginDataKey, value)
}

func FromContext(ctx context.Context) (*LoginData, bool) {
	u, ok := ctx.Value(loginDataKey).(*LoginData)
	return u, ok
}
