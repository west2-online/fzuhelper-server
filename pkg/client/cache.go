package client

import (
	"context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// NewRedisClient 传入dbName， 比如classroom将key放在db0中，user放在db1中
func NewRedisClient(dbName int) (*redis.Client, error) {
	//首先判断config的redis是否初始化过。如果没有，则该field应该是空指针，我们需要报错返回
	conf := config.Redis
	if conf == nil {
		logger.LoggerObj.Fatalf("The redis config init failed")
		return nil, errors.New("The redis config init failed")
	}
	logger.LoggerObj.Infof("redis addr: %s", conf.Addr)
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       dbName,
	})
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		logger.LoggerObj.Fatalf("redis client ping failed: %v", err)
		return nil, errors.Wrap(err, "redis client ping failed")
	}
	return client, nil
}
