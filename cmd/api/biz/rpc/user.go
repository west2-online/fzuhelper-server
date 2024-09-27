package rpc

import (
	"context"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/west2-online/fzuhelper-server/cmd/api/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitUserRPC() {
	conf := config.Config
	r, err := etcd.NewEtcdResolver([]string{conf.EtcdHost + ":" + conf.EtcdPort})
	if err != nil {
		panic(err)
	}
	userClient, err = userservice.NewClient("user", client.WithResolver(r))
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
	if resp.Base.Code != errno.SuccessCode {
		utils.LoggerObj.Errorf("api.rpc.user GetLoginDataRPC received failed")
		return "", nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Id, resp.Cookies, nil
}
