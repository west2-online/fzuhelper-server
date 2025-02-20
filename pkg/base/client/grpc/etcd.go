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
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// EtcdResolver 用来封装 etcd 客户端和解析出来的 gRPC endpoint 列表
type EtcdResolver struct {
	EtcdClient  *clientv3.Client
	ServiceKey  string       // etcd 里存服务地址列表的 key，如 "/services/ai_agent/grpc"
	Endpoints   []string     // 当前解析到的服务列表
	endpointsMu sync.RWMutex // 读写锁，保证多协程安全
}

// initEtcdClient 允许即使没有初始 endpoints 也能启动
func initEtcdClient(serviceKey string) (*EtcdResolver, error) {
	if config.Etcd == nil || config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is nil")
	}
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.Etcd.Addr},
		DialTimeout: constants.EtcdDialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("connect etcd failed: %w", err)
	}

	c := &EtcdResolver{
		EtcdClient: etcdCli,
		ServiceKey: serviceKey,
	}

	// 尝试加载 endpoints，但不因为空而失败
	if err := c.loadEndpoints(); err != nil {
		log.Printf("Warning: %v\n", err)
	}

	return c, nil
}

// loadEndpoints 尝试从 etcd 中加载 endpoints，并统一解析为 []string
func (c *EtcdResolver) loadEndpoints() error {
	resp, err := c.EtcdClient.Get(context.Background(), c.ServiceKey)
	if err != nil {
		return fmt.Errorf("etcd get failed: %w", err)
	}
	if len(resp.Kvs) == 0 {
		// 允许初始启动时 endpoints 为空
		return nil
	}
	val := resp.Kvs[0].Value
	var eps []string
	// 尝试解析为 JSON 数组
	if err := sonic.Unmarshal(val, &eps); err != nil {
		// 如果解析失败，认为是单个 endpoint
		eps = []string{string(val)}
	}
	c.endpointsMu.Lock()
	c.Endpoints = eps
	c.endpointsMu.Unlock()
	return nil
}

// watchEndpoints 监听 etcd 里 key 的变化
func (c *AIAgentClient) watchEndpoints() {
	watchCh := c.EtcdResolver.EtcdClient.Watch(context.Background(), c.EtcdResolver.ServiceKey)
	for wResp := range watchCh {
		for _, ev := range wResp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				var eps []string
				if err := sonic.Unmarshal(ev.Kv.Value, &eps); err != nil {
					// 如果解析失败，认为是单个 endpoint
					eps = []string{string(ev.Kv.Value)}
				}
				c.EtcdResolver.endpointsMu.Lock()
				c.EtcdResolver.Endpoints = eps
				c.EtcdResolver.endpointsMu.Unlock()
				// 进行refresh
				c.RefreshGRPC()
				// log.Printf("Endpoints updated: %v\n", eps)
			case clientv3.EventTypeDelete:
				c.EtcdResolver.endpointsMu.Lock()
				c.EtcdResolver.Endpoints = nil
				c.EtcdResolver.endpointsMu.Unlock()
				// log.Println("Endpoints key deleted in etcd, no valid servers now.")
			}
		}
	}
}

// getRandomEndpoint 随机返回一个 endpoint
func (c *EtcdResolver) getRandomEndpoint() (string, error) {
	c.endpointsMu.RLock()
	defer c.endpointsMu.RUnlock()
	if len(c.Endpoints) == 0 {
		return "", errors.New("no endpoints available")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return c.Endpoints[r.Intn(len(c.Endpoints))], nil
}
