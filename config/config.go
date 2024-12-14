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
	"errors"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"

	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

var (
	Server               *server
	Mysql                *mySQL
	Snowflake            *snowflake
	Service              *service
	Jaeger               *jaeger
	Etcd                 *etcd
	Redis                *redis
	DefaultUser          *defaultUser
	Elasticsearch        *elasticsearch
	Kafka                *kafka
	UpYun                *upyun
	VersionUploadService *url
	runtimeViper         = viper.New()
)

func Init(service string) {
	// 从环境变量中获取 etcd 地址
	etcdAddr := os.Getenv("ETCD_ADDR")
	if etcdAddr == "" {
		logger.Fatalf("config.Init: etcd addr is empty")
	}
	logger.Infof("config.Init: etcd addr: %v", etcdAddr)
	Etcd = &etcd{Addr: etcdAddr}

	// use etcd for config save
	err := runtimeViper.AddRemoteProvider("etcd3", Etcd.Addr, "/config")
	if err != nil {
		logger.Fatalf("config.Init: add remote provider error: %v", err)
	}
	runtimeViper.SetConfigName("config")
	runtimeViper.SetConfigType("yaml")
	if err := runtimeViper.ReadRemoteConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			logger.Fatal("config.Init: could not find config files")
		}
		logger.Fatal("config.Init: read config error: %v", err)
	}
	configMapping(service)
	// 持续监听配置
	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		logger.Infof("config: config file changed: %v\n", e.String())
	})
	runtimeViper.WatchConfig()
}

func configMapping(srv string) {
	c := new(config)
	if err := runtimeViper.Unmarshal(&c); err != nil {
		logger.Fatalf("config.configMapping: config: unmarshal error: %v", err)
	}
	Snowflake = &c.Snowflake

	Server = &c.Server

	Jaeger = &c.Jaeger
	Mysql = &c.MySQL
	Redis = &c.Redis
	Elasticsearch = &c.Elasticsearch
	Kafka = &c.Kafka
	DefaultUser = &c.DefaultUser
	VersionUploadService = &c.Url
	upy, ok := c.UpYuns[srv]
	if ok {
		UpYun = &upy
	}

	Service = GetService(srv)
}

func GetService(srvname string) *service {
	logger.Debugf("get service name: %v", srvname)
	addrlist := runtimeViper.GetStringSlice("services." + srvname + ".addr")
	logger.Debugf("get addrlist: %v", addrlist)

	return &service{
		Name:     runtimeViper.GetString("services." + srvname + ".name"),
		AddrList: addrlist,
		LB:       runtimeViper.GetBool("services." + srvname + ".load-balance"),
	}
}
