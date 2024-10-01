package client

import (
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

//定义一系列的RPC客户端

func InitUserRPC() (*userservice.Client, error) {
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("InitUserRPC etcd.NewEtcdResolver failed: %w", err)
	}
	client, err := userservice.NewClient(constants.UserServiceName, client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		return nil, fmt.Errorf("InitUserRPC NewClient failed: %w", err)
	}
	return &client, nil
}

func InitClassroomRPC() (*classroomservice.Client, error) {
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("InitClassroomRPC etcd.NewEtcdResolver failed: %w", err)
	}
	client, err := classroomservice.NewClient(constants.ClassroomServiceName, client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		return nil, fmt.Errorf("InitClassroomRPC NewClient failed: %w", err)
	}
	return &client, nil
}
