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
	userClient, err = userservice.NewClient("user", client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		panic(err)
	}
	utils.LoggerObj.Info("InitUserRPC success")
}

func GetLoginDataRPC(ctx context.Context, req *user.GetLoginDataRequest) (string, []string, error) {
	resp, err := userClient.GetLoginData(ctx, req)
	if err != nil {
		utils.LoggerObj.Errorf("api.rpc.user GetLoginDataRPC received rpc error %v", err)
		return "", nil, err
	}
	if err = utils.IsSuccess(resp.Base); err != nil {
		return "", nil, err
	}
	return resp.Id, resp.Cookies, nil
}
