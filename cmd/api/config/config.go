package config

import (
	"github.com/spf13/viper"
)

var Config Conf

type Conf struct {
	System *System `mapstructure:"system"`
	Etcd   `mapstructure:"etcd"`
}

type System struct {
	Domain string `mapstructure:"domain"`
	Host   string `mapstructure:"host"`
	Port   string `mapstructure:"port"`
}

type Etcd struct {
	EtcdHost string `mapstruct:"etcdHost"`
	EtcdPort string `mapstruct:"etcdPort"`
}

// Init for testing called
func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}
}
