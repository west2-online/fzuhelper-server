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
	"bytes"
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

var (
	Server               *server
	SignedLocationApiUrl *signedLocationApiUrl
	MCP                  *mcp
	AI                   *ai
	Mysql                *mySQL
	Snowflake            *snowflake
	Admin                *admin
	Service              *service
	Jaeger               *jaeger
	Otel                 *otel
	Etcd                 *etcd
	Redis                *redis
	DefaultUser          *defaultUser
	Elasticsearch        *elasticsearch
	Kafka                *kafka
	UpYun                *upyun
	Umeng                *umeng
	VersionUploadService *url
	Vendors              *vendors
	Friend               *friend
	APIMonitor           *apiMonitorConfig
	runtimeViper         = viper.New()
)

const (
	// remotePath 是配置在 etcd 中的存储 key，value 为 yaml 内容。
	// 原实现通过 viper/remote + sagikazarmark/crypt 读取，但该依赖链会引入
	// golang.org/x/crypto/openpgp（GO-2026-5932，该包无人维护、无修复版本），
	// 因此改为直接使用 etcd client/v3 读取，行为保持一致。
	remotePath     = "/config"
	remoteFileType = "yaml"
)

func Init(service string) {
	DeployEnv := os.Getenv("DEPLOY_ENV")
	if DeployEnv == "k8s" {
		InitFromConfigMap(service)
	} else {
		InitFromETCD(service)
	}
}

// InitFromETCD 目的是初始化并读入配置，此时没有初始化Logger，但仍然可以用 logger 来输出，只是没有自定义配置
func InitFromETCD(service string) {
	// 从环境变量中获取 etcd 地址
	etcdAddr := os.Getenv("ETCD_ADDR")
	if etcdAddr == "" {
		logger.Fatalf("config.Init: etcd addr is empty")
	}
	logger.Infof("config.Init: etcd addr: %v", etcdAddr)
	Etcd = &etcd{Addr: etcdAddr}

	// 配置存储在 etcd 中
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{etcdAddr},
	})
	if err != nil {
		logger.Fatalf("config.Init: create etcd client error: %v", err)
	}
	defer func() {
		if err := cli.Close(); err != nil {
			logger.Errorf("config.Init: close etcd client error: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := cli.Get(ctx, remotePath)
	if err != nil {
		logger.Fatalf("config.Init: read config from etcd error: %v", err)
	}
	if len(resp.Kvs) != 1 {
		logger.Fatalf("config.Init: config not found in etcd, key: %s", remotePath)
	}

	runtimeViper.SetConfigType(remoteFileType)
	if err := runtimeViper.ReadConfig(bytes.NewReader(resp.Kvs[0].Value)); err != nil {
		logger.Fatalf("config.Init: read config error: %v", err)
	}
	configMapping(service)

	// 设置持续监听
	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		// 我们无法确定监听到配置变更时是否已经初始化完毕，所以此处需要做一个判断
		logger.Infof("config: notice config changed: %v\n", e.String())
		configMapping(service) // 重新映射配置
	})
	runtimeViper.WatchConfig()
}

// InitFromConfigMap 用于从 k8s 的 ConfigMap 中初始化配置
// 方式是通过 pod 去挂载 configMap，然后容器再读取本地的config.yaml来初始化配置
// 优点：不再依赖 etcd，并且 k8s 会自动更新 ConfigMap，所以配置也会自动更新（热更新），不需要另外设置 etcd 来自定义启动脚本
// config 默认在 /app/config/config.yaml
func InitFromConfigMap(service string) {
	runtimeViper.AddConfigPath("./config")
	runtimeViper.SetConfigName("config")
	runtimeViper.SetConfigType("yaml")
	if err := runtimeViper.ReadInConfig(); err != nil {
		logger.Fatalf("config.InitFromConfigMap: read config error: %v", err)
	}
	configMapping(service)
	// 设置持续监听
	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		logger.Infof("config: notice config changed: %v\n", e.String())
		configMapping(service) // 重新映射配置
	})
	runtimeViper.WatchConfig()
}

// configMapping 用于将配置映射到全局变量
func configMapping(srv string) {
	c := new(config)
	if err := runtimeViper.Unmarshal(&c); err != nil {
		// 由于这个函数会在配置重载时被再次触发，所以需要判断日志记录方式
		logger.Fatalf("config.configMapping: config: unmarshal error: %v", err)
	}
	Snowflake = &c.Snowflake
	Server = &c.Server
	MCP = &c.MCP
	SignedLocationApiUrl = &c.SignedLocationApiUrl
	Admin = &c.Admin
	AI = &c.AI
	Jaeger = &c.Jaeger
	Otel = &c.Otel
	Mysql = &c.MySQL
	Redis = &c.Redis
	Elasticsearch = &c.Elasticsearch
	Kafka = &c.Kafka
	DefaultUser = &c.DefaultUser
	VersionUploadService = &c.Url
	Umeng = &c.Umeng
	Friend = &c.Friend
	APIMonitor = &c.APIMonitor
	if upy, ok := c.UpYuns[srv]; ok {
		UpYun = &upy
	}
	Vendors = &c.Vendors
	Service = getService(srv)
}

func getService(name string) *service {
	addrList := runtimeViper.GetStringSlice("services." + name + ".addr")

	return &service{
		Name:     runtimeViper.GetString("services." + name + ".name"),
		AddrList: addrList,
		LB:       runtimeViper.GetBool("services." + name + ".load-balance"),
	}
}

// GetLoggerLevel 会返回服务的日志等级
func GetLoggerLevel() string {
	if Server == nil {
		return constants.DefaultLogLevel
	}
	return Server.LogLevel
}

// InitForTest 专门用于测试环境的配置初始化
// 会读取config.example.yaml文件
func InitForTest(service string) error {
	// 寻找项目根目录的config.example.yaml文件
	configPath := findConfigFile("config.example.yaml")
	if configPath == "" {
		logger.Fatalf("config.InitForTest: config.example.yaml not found")
	}

	// 直接指定配置文件的完整路径
	runtimeViper.SetConfigFile(configPath)

	if err := runtimeViper.ReadInConfig(); err != nil {
		logger.Fatalf("config.InitForTest: read config error: %v", err)
	}
	configMapping(service)

	return nil
}

// findConfigFile 从当前目录开始向上查找配置文件
func findConfigFile(filename string) string {
	// 首先尝试当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	// 向上查找直到找到文件或到达根目录
	for {
		configPath := filepath.Join(currentDir, "config", filename)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		// 尝试直接在当前目录查找
		configPath = filepath.Join(currentDir, filename)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		// 向上一级目录
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// 已经到达根目录
			break
		}
		currentDir = parentDir
	}

	return ""
}
