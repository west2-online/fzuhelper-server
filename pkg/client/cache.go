package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/fzuhelper-server/config"
)

// NewRedisClient 传入dbName， 比如classroom将key放在db0中，user放在db1中
func NewRedisClient(dbName int) (*redis.Client, error) {
	// 首先判断config的redis是否初始化过。如果没有，则该field应该是空指针，我们需要报错返回
	if config.Redis == nil {
		return nil, errors.New("redis config is nil")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       dbName,
	})
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, fmt.Errorf("client.NewRedisClient: ping redis failed: %w", err)
	}
	return client, nil
}
