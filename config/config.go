package config

import (
	"errors"
	"log"
	"os"

	"github.com/cloudwego/kitex/pkg/klog"
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
		panic(errors.New("not found etcd addr in env"))
	}

	Etcd = &etcd{Addr: etcdAddr}

	// use etcd for config save
	err := runtime_viper.AddRemoteProvider("etcd3", Etcd.Addr, "/config/config.yaml")

	if err != nil {
		panic(err)
	}

	klog.Infof("config path: %v\n", path)

	if err := runtime_viper.ReadRemoteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			klog.Fatal("could not find config files")
		} else {
			klog.Fatal("read config error: %v", err)
		}
		klog.Fatal(err)
	}

	configMapping(service)
	klog.Infof("all keys: %v\n", runtime_viper.AllKeys())
	// 持续监听配置
	runtime_viper.OnConfigChange(func(e fsnotify.Event) {
		klog.Infof("config file changed: %v\n", e.String())
	})
	runtime_viper.WatchConfig()
}

func configMapping(srv string) {
	c := new(config)
	if err := runtime_viper.Unmarshal(&c); err != nil {
		log.Fatal(err)
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
