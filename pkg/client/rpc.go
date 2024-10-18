/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"errors"
	"fmt"

	"github.com/cloudwego/kitex/client/streamclient"

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen/launchscreenservice"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// 定义一系列的RPC客户端

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

func InitLaunchScreenRPC() (*launchscreenservice.Client, error) {
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("InitLaunchScreenRPC etcd.NewEtcdResolver failed: %w", err)
	}
	client, err := launchscreenservice.NewClient(constants.LaunchScreenServiceName, client.WithResolver(r), client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		return nil, fmt.Errorf("InitLaunchScreenRPC NewClient failed: %w", err)
	}
	return &client, nil
}

func InitLaunchScreenStreamRPC() (*launchscreenservice.StreamClient, error) {
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("InitLaunchScreenStreamRPC etcd.NewEtcdResolver failed: %w", err)
	}
	streamClient := launchscreenservice.MustNewStreamClient(constants.LaunchScreenServiceName, streamclient.WithResolver(r))

	return &streamClient, nil
}
