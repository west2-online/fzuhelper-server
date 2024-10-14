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

package config

import (
	"os"

	"github.com/west2-online/fzuhelper-server/pkg/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	_ "github.com/spf13/viper/remote"
)

var (
	Server        *server
	Mysql         *mySQL
	Snowflake     *snowflake
	Service       *service
	Jaeger        *jaeger
	Etcd          *etcd
	RabbitMQ      *rabbitMQ
	Redis         *redis
	DefaultUser   *defaultUser
	OSS           *oss
	Elasticsearch *elasticsearch

	runtime_viper = viper.New()
)

func Init(path string, service string) {
	runtime_viper.SetConfigType("yaml")
	runtime_viper.AddConfigPath(path)

	etcdAddr := os.Getenv("ETCD_ADDR")

	if etcdAddr == "" {
		logger.Fatalf("config.Init: etcd addr is empty")
	}
	logger.Infof("config.Init: etcd addr: %v", etcdAddr)
	Etcd = &etcd{Addr: etcdAddr}

	// use etcd for config save
	err := runtime_viper.AddRemoteProvider("etcd3", Etcd.Addr, "/config/config.yaml")
	if err != nil {
		logger.Fatalf("config.Init: add remote provider error: %v", err)
	}
	logger.Infof("config.Init: config path: %v", path)

	if err := runtime_viper.ReadRemoteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Fatal("config.Init: could not find config files")
		} else {
			logger.Fatalf("config.Init: read config error: %v", err)
		}
		logger.Fatalf("config.Init: read config error: %v", err)
	}

	configMapping(service)
	// logger.Infof("all keys: %v\n", runtime_viper.AllKeys())
	// 持续监听配置
	runtime_viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Infof("config: config file changed: %v\n", e.String())
	})
	runtime_viper.WatchConfig()
}

func configMapping(srv string) {
	c := new(config)
	if err := runtime_viper.Unmarshal(&c); err != nil {
		logger.Fatalf("config.configMapping: config: unmarshal error: %v", err)
	}
	Snowflake = &c.Snowflake

	Server = &c.Server
	Server.Secret = []byte(runtime_viper.GetString("server.jwt-secret"))

	Jaeger = &c.Jaeger
	Mysql = &c.MySQL
	RabbitMQ = &c.RabbitMQ
	Redis = &c.Redis
	OSS = &c.OSS
	Elasticsearch = &c.Elasticsearch
	DefaultUser = &c.DefaultUser
	Service = GetService(srv)
}

func GetService(srvname string) *service {
	addrlist := runtime_viper.GetStringSlice("services." + srvname + ".addr")

	return &service{
		Name:     runtime_viper.GetString("services." + srvname + ".name"),
		AddrList: addrlist,
		LB:       runtime_viper.GetBool("services." + srvname + ".load-balance"),
	}
}
