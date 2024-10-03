package rpc

import (
	"context"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitUserRPC() {
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}
	userClient, err = userservice.NewClient(constants.UserServiceName, client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

func LoginRPC(ctx context.Context, req *user.LoginRequest) (token *string, err error) {
	resp, err := userClient.Login(ctx, req)
	if ok, err := utils.IsSuccess(err, resp.Base); !ok {
		return nil, err
	}
	return resp.Token, nil
}

func RegisterRPC(ctx context.Context, req *user.RegisterRequest) (uid *int64, err error) {
	resp, err := userClient.Register(ctx, req)
	if ok, err := utils.IsSuccess(err, resp.Base); !ok {
		return nil, err
	}
	return resp.UserId, nil
}
