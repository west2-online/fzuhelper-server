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

package grpc

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"

	"github.com/west2-online/fzuhelper-server/kitex_gen/ai_agent"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

type AIAgentClient struct {
	Cli          ai_agent.AIAgentClient
	EtcdResolver *EtcdResolver
	refreshMu    sync.RWMutex // 控制 RefreshGRPC 的互斥锁
	refreshing   int32        // 原子标记（0=未刷新，1=刷新中）,保证同时只有一个RefreshGRPC被调用
}

func InitAiAgentRPC() (*AIAgentClient, error) {
	aiAgentClient := new(AIAgentClient)
	c, err := initEtcdClient(constants.AiAgentServiceName)
	if err != nil {
		logger.Errorf("Failed to connect: %v", err)
	}
	// 首次连接建立(watchEndPoints只会在之后连接刷新时才触发，对于首次不适用)
	conn, err := c.initConn()
	if err != nil {
		logger.Infof("get conn failed: %v", err)
	}
	client := ai_agent.NewAIAgentClient(conn)
	aiAgentClient.Cli = client
	aiAgentClient.EtcdResolver = c
	go aiAgentClient.watchEndpoints()

	return aiAgentClient, nil
}

// initConn 初始化时获取一个 grpc 连接
func (c *EtcdResolver) initConn() (*grpc.ClientConn, error) {
	var ep string
	var err error
	// 尝试等待 endpoints 可用
	timeout := time.After(constants.RefreshGRPCTimeout)
	ticker := time.NewTicker(constants.GRPCGetEndPointTicker)
	defer ticker.Stop()
	for {
		ep, err = c.getRandomEndpoint()
		if err == nil {
			break
		}
		select {
		case <-timeout:
			return nil, nil
		case <-ticker.C:
			continue
		}
	}
	conn, err := grpc.NewClient(ep, grpc.WithIdleTimeout(constants.GRPCGetConnTimeout))
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}
	return conn, nil
}

// RefreshGRPC 用于在etcd节点数更新时刷新client，并且保证了读写安全
func (c *AIAgentClient) RefreshGRPC() {
	if !atomic.CompareAndSwapInt32(&c.refreshing, 0, 1) {
		return // 已有协程在刷新，直接退出
	}
	defer atomic.StoreInt32(&c.refreshing, 0)
	c.refreshMu.Lock() // 加锁
	defer c.refreshMu.Unlock()

	var ep string
	var err error
	// 尝试等待 endpoints 可用
	timeout := time.After(constants.RefreshGRPCTimeout)
	ticker := time.NewTicker(constants.GRPCGetEndPointTicker)
	defer ticker.Stop()
	for {
		ep, err = c.EtcdResolver.getRandomEndpoint()
		if err == nil {
			break
		}
		select {
		case <-timeout:
			logger.Warnf("no instance of ai_agent found")
		case <-ticker.C:
			continue
		}
	}

	conn, err := grpc.NewClient(ep, grpc.WithIdleTimeout(constants.GRPCGetConnTimeout))
	if err != nil {
		logger.Errorf("dial failed: %v", err)
	}

	c.Cli = ai_agent.NewAIAgentClient(conn)
}
