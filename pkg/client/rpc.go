package client

import (
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/pkg/errors"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

//定义一系列的RPC客户端

func InitUserRPC() (*userservice.Client, error) {
	if config.Etcd.Addr == "" {
		return nil, errors.Wrap(errno.InternalServiceError, "config.Etcd.Addr is nil")
	}
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, errors.Wrap(err, "InitUserRPC etcd.NewEtcdResolver failed")
	}
	client, err := userservice.NewClient("user", client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		return nil, errors.Wrap(err, "InitUserRPC NewClient failed")
	}
	return &client, nil
}

func InitClassroomRPC() (*classroomservice.Client, error) {
	if config.Etcd.Addr == "" {
		return nil, errors.Wrap(errno.InternalServiceError, "config.Etcd.Addr is nil")
	}
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, errors.Wrap(err, "pkg.rpc.rpc.InitClassroomRPC etcd.NewEtcdResolver failed")
	}
	client, err := classroomservice.NewClient(constants.ClassroomService, client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		return nil, errors.Wrap(err, "pkg.rpc.rpc.InitClassroomRPC NewClient failed")
	}
	return &client, nil
}
