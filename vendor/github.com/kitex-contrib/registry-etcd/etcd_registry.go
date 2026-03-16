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
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/kitex-contrib/registry-etcd/retry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ttlKey              = "KITEX_ETCD_REGISTRY_LEASE_TTL"
	defaultTTL          = 60
	kitexIpToRegistry   = "KITEX_IP_TO_REGISTRY"
	kitexPortToRegistry = "KITEX_PORT_TO_REGISTRY"
)

type etcdRegistry struct {
	etcdClient  *clientv3.Client
	leaseTTL    int64
	meta        *registerMeta
	retryConfig *retry.Config
	stop        chan struct{}
	address     net.Addr
	prefix      string
}

type registerMeta struct {
	leaseID clientv3.LeaseID
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewEtcdRegistry creates an etcd based registry.
func NewEtcdRegistry(endpoints []string, opts ...Option) (registry.Registry, error) {
	cfg := &Config{
		EtcdConfig: &clientv3.Config{
			Endpoints: endpoints,
		},
		Prefix: "kitex/registry-etcd",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	etcdClient, err := clientv3.New(*cfg.EtcdConfig)
	if err != nil {
		return nil, err
	}
	retryConfig := retry.NewRetryConfig()
	return &etcdRegistry{
		etcdClient:  etcdClient,
		leaseTTL:    getTTL(),
		retryConfig: retryConfig,
		stop:        make(chan struct{}, 1),
		prefix:      cfg.Prefix,
	}, nil
}

// SetFixedAddress sets the fixed address for registering
// setting address to nil to clear the previous address
func SetFixedAddress(r registry.Registry, address net.Addr) {
	if er, ok := r.(*etcdRegistry); ok {
		er.address = address
		return
	}
	panic("invalid registry type: not etcdRegistry")
}

// NewEtcdRegistryWithRetry creates an etcd based registry with given custom retry configs
func NewEtcdRegistryWithRetry(endpoints []string, retryConfig *retry.Config, opts ...Option) (registry.Registry, error) {
	cfg := &Config{
		EtcdConfig: &clientv3.Config{
			Endpoints: endpoints,
		},
		Prefix: "kitex/registry-etcd",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	etcdClient, err := clientv3.New(*cfg.EtcdConfig)
	if err != nil {
		return nil, err
	}
	return &etcdRegistry{
		etcdClient:  etcdClient,
		leaseTTL:    getTTL(),
		retryConfig: retryConfig,
		stop:        make(chan struct{}, 1),
		prefix:      cfg.Prefix,
	}, nil
}

// NewEtcdRegistryWithAuth creates an etcd based registry with given username and password.
// Deprecated: Use WithAuthOpt instead.
func NewEtcdRegistryWithAuth(endpoints []string, username, password string) (registry.Registry, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}
	retryConfig := retry.NewRetryConfig()
	return &etcdRegistry{
		etcdClient:  etcdClient,
		leaseTTL:    getTTL(),
		retryConfig: retryConfig,
		stop:        make(chan struct{}, 1),
	}, nil
}

// Register registers a server with given registry info.
func (e *etcdRegistry) Register(info *registry.Info) error {
	if err := validateRegistryInfo(info); err != nil {
		return err
	}
	leaseID, err := e.grantLease()
	if err != nil {
		return err
	}

	if err := e.register(info, leaseID); err != nil {
		return err
	}
	meta := registerMeta{
		leaseID: leaseID,
	}
	meta.ctx, meta.cancel = context.WithCancel(context.Background())
	if err := e.keepalive(&meta); err != nil {
		return err
	}
	e.meta = &meta
	return nil
}

// Deregister deregisters a server with given registry info.
func (e *etcdRegistry) Deregister(info *registry.Info) error {
	if info.ServiceName == "" {
		return fmt.Errorf("missing service name in Deregister")
	}
	if err := e.deregister(info); err != nil {
		return err
	}
	e.meta.cancel()
	return nil
}

func (e *etcdRegistry) register(info *registry.Info, leaseID clientv3.LeaseID) error {
	addr, err := e.getAddressOfRegistration(info)
	if err != nil {
		return err
	}
	network := info.Addr.Network()
	if e.address != nil {
		network = e.address.Network()
		addr = e.address.String()
	}
	val, err := json.Marshal(&instanceInfo{
		Network: network,
		Address: addr,
		Weight:  info.Weight,
		Tags:    info.Tags,
	})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err = e.etcdClient.Put(ctx, serviceKey(e.prefix, info.ServiceName, addr), string(val), clientv3.WithLease(leaseID))
	if err != nil {
		return err
	}

	go func(key, val string) {
		e.keepRegister(key, val, e.retryConfig)
	}(serviceKey(e.prefix, info.ServiceName, addr), string(val))

	return nil
}

// keepRegister keep service registered status
// maxRetry == 0 means retry forever
func (e *etcdRegistry) keepRegister(key, val string, retryConfig *retry.Config) {
	var failedTimes uint
	delay := retryConfig.ObserveDelay
	for retryConfig.MaxAttemptTimes == 0 || failedTimes < retryConfig.MaxAttemptTimes {
		select {
		case _, ok := <-e.stop:
			if !ok {
				close(e.stop)
			}
			klog.Infof("stop keep register service %s", key)
			return
		case <-time.After(delay):
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		resp, err := e.etcdClient.Get(ctx, key)
		if err != nil {
			klog.Warnf("keep register get %s failed with err: %v", key, err)
			delay = retryConfig.RetryDelay
			failedTimes++
			continue
		}

		if len(resp.Kvs) == 0 {
			klog.Infof("keep register service %s", key)
			delay = retryConfig.RetryDelay
			leaseID, err := e.grantLease()
			if err != nil {
				klog.Warnf("keep register grant lease %s failed with err: %v", key, err)
				failedTimes++
				continue
			}

			_, err = e.etcdClient.Put(ctx, key, val, clientv3.WithLease(leaseID))
			if err != nil {
				klog.Warnf("keep register put %s failed with err: %v", key, err)
				failedTimes++
				continue
			}

			meta := registerMeta{
				leaseID: leaseID,
			}
			meta.ctx, meta.cancel = context.WithCancel(context.Background())
			if err := e.keepalive(&meta); err != nil {
				klog.Warnf("keep register keepalive %s failed with err: %v", key, err)
				failedTimes++
				continue
			}
			e.meta.cancel()
			e.meta = &meta
			delay = retryConfig.ObserveDelay
		}

		failedTimes = 0
	}
	klog.Errorf("keep register service %s failed times:%d", key, failedTimes)
}

func (e *etcdRegistry) deregister(info *registry.Info) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	addr, err := e.getAddressOfRegistration(info)
	if err != nil {
		return err
	}
	_, err = e.etcdClient.Delete(ctx, serviceKey(e.prefix, info.ServiceName, addr))
	if err != nil {
		return err
	}
	e.stop <- struct{}{}
	return nil
}

func (e *etcdRegistry) grantLease() (clientv3.LeaseID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := e.etcdClient.Grant(ctx, e.leaseTTL)
	if err != nil {
		return clientv3.NoLease, err
	}
	return resp.ID, nil
}

func (e *etcdRegistry) keepalive(meta *registerMeta) error {
	keepAlive, err := e.etcdClient.KeepAlive(meta.ctx, meta.leaseID)
	if err != nil {
		return err
	}
	go func() {
		// eat keepAlive channel to keep related lease alive.
		klog.Infof("start keepalive lease %x for etcd registry", meta.leaseID)
		for range keepAlive {
			select {
			case <-meta.ctx.Done():
				klog.Infof("stop keepalive lease %x for etcd registry", meta.leaseID)
				return
			default:
			}
		}
	}()
	return nil
}

// getAddressOfRegistration returns the address of the service registration.
func (e *etcdRegistry) getAddressOfRegistration(info *registry.Info) (string, error) {

	host, port, err := net.SplitHostPort(info.Addr.String())
	if err != nil {
		return "", err
	}

	// if host is empty or "::", use local ipv4 address as host
	if host == "" || host == "::" {
		host, err = getLocalIPv4Host()
		if err != nil {
			return "", fmt.Errorf("parse registry info addr error: %w", err)
		}
	}

	// if env KITEX_IP_TO_REGISTRY is set, use it as host
	if ipToRegistry, exists := os.LookupEnv(kitexIpToRegistry); exists && ipToRegistry != "" {
		host = ipToRegistry
	}

	// if env KITEX_PORT_TO_REGISTRY is set, use it as port
	if portToRegistry, exists := os.LookupEnv(kitexPortToRegistry); exists && portToRegistry != "" {
		port = portToRegistry
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("parse registry info port error: %w", err)
	}

	return fmt.Sprintf("%s:%d", host, p), nil
}

func validateRegistryInfo(info *registry.Info) error {
	if info.ServiceName == "" {
		return fmt.Errorf("missing service name in Register")
	}
	if strings.Contains(info.ServiceName, "/") {
		return fmt.Errorf("service name registered with etcd should not include character '/'")
	}
	if info.Addr == nil {
		return fmt.Errorf("missing addr in Register")
	}
	return nil
}

func getTTL() int64 {
	var ttl int64 = defaultTTL
	if str, ok := os.LookupEnv(ttlKey); ok {
		if t, err := strconv.Atoi(str); err == nil {
			ttl = int64(t)
		}
	}
	return ttl
}

func getLocalIPv4Host() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addr {
		ipNet, isIpNet := addr.(*net.IPNet)
		if isIpNet && !ipNet.IP.IsLoopback() {
			ipv4 := ipNet.IP.To4()
			if ipv4 != nil {
				return ipv4.String(), nil
			}
		}
	}
	return "", fmt.Errorf("not found ipv4 address")
}
