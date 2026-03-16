// Copyright 2021 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package etcd resolver
package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	defaultWeight = 10
)

// etcdResolver is a resolver using etcd.
type etcdResolver struct {
	etcdClient    *clientv3.Client
	prefix        string
	defaultWeight int
}

// NewEtcdResolver creates a etcd based resolver.
func NewEtcdResolver(endpoints []string, opts ...Option) (discovery.Resolver, error) {
	cfg := &Config{
		EtcdConfig: &clientv3.Config{
			Endpoints: endpoints,
		},
		Prefix:        "kitex/registry-etcd",
		DefaultWeight: defaultWeight,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	etcdClient, err := clientv3.New(*cfg.EtcdConfig)
	if err != nil {
		return nil, err
	}
	return &etcdResolver{
		etcdClient:    etcdClient,
		prefix:        cfg.Prefix,
		defaultWeight: cfg.DefaultWeight,
	}, nil
}

// NewEtcdResolverWithAuth creates a etcd based resolver with given username and password.
// Deprecated: Use WithAuthOpt instead.
func NewEtcdResolverWithAuth(endpoints []string, username, password string) (discovery.Resolver, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}
	return &etcdResolver{
		etcdClient: etcdClient,
	}, nil
}

// Target implements the Resolver interface.
func (e *etcdResolver) Target(ctx context.Context, target rpcinfo.EndpointInfo) (description string) {
	return target.ServiceName()
}

// Resolve implements the Resolver interface.
func (e *etcdResolver) Resolve(ctx context.Context, desc string) (discovery.Result, error) {
	prefix := serviceKeyPrefix(e.prefix, desc)
	resp, err := e.etcdClient.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return discovery.Result{}, err
	}
	var eps []discovery.Instance
	for _, kv := range resp.Kvs {
		var info instanceInfo
		err = json.Unmarshal(kv.Value, &info)
		if err != nil {
			klog.Warnf("fail to unmarshal with err: %v, ignore key: %v", err, string(kv.Key))
			continue
		}
		weight := info.Weight
		if weight <= 0 {
			weight = e.defaultWeight
		}
		eps = append(eps, discovery.NewInstance(info.Network, info.Address, weight, info.Tags))
	}
	if len(eps) == 0 {
		return discovery.Result{}, fmt.Errorf("no instance remains for %v", desc)
	}
	return discovery.Result{
		Cacheable: true,
		CacheKey:  desc,
		Instances: eps,
	}, nil
}

// Diff implements the Resolver interface.
func (e *etcdResolver) Diff(cacheKey string, prev, next discovery.Result) (discovery.Change, bool) {
	return discovery.DefaultDiff(cacheKey, prev, next)
}

// Name implements the Resolver interface.
func (e *etcdResolver) Name() string {
	return "etcd"
}

func (e *etcdResolver) GetPrefix() string {
	return e.prefix
}

func (e *etcdResolver) GetDefaultWeight() int {
	return e.defaultWeight
}
