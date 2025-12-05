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
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

var (
	Server               *server
	MCP                  *mcp
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
	Umeng                *umeng
	VersionUploadService *url
	Vendors              *vendors
	Friend               *friend
	runtimeViper         = viper.New()
)

const (
	remoteProvider = "etcd3" // 使用 etcd3
	remotePath     = "/config"
	remoteFileName = "config"
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
	err := runtimeViper.AddRemoteProvider(remoteProvider, Etcd.Addr, remotePath)
	if err != nil {
		logger.Fatalf("config.Init: add remote provider error: %v", err)
	}
	runtimeViper.SetConfigName(remoteFileName)
	runtimeViper.SetConfigType(remoteFileType)
	if err := runtimeViper.ReadRemoteConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Fatal("config.Init: could not find config files")
		}
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
	Jaeger = &c.Jaeger
	Mysql = &c.MySQL
	Redis = &c.Redis
	Elasticsearch = &c.Elasticsearch
	Kafka = &c.Kafka
	DefaultUser = &c.DefaultUser
	VersionUploadService = &c.Url
	Umeng = &c.Umeng
	Friend = &c.Friend
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
