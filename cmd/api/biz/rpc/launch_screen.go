package rpc

import (
	"context"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen/launchscreenservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitLaunchScreenRPC() {
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}
	launchScreenClient, err = launchscreenservice.NewClient(constants.LaunchScreenServiceName, client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

func CreateImageRPC(ctx context.Context, req *launch_screen.CreateImageRequest) (image *model.Picture, err error) {
	resp, err := launchScreenClient.CreateImage(ctx, req)
	if ok, err := utils.IsSuccess(err, resp.Base); !ok {
		return nil, err
	}
	return resp.Picture, nil
}
