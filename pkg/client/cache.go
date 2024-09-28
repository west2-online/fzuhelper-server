package client

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/west2-online/fzuhelper-server/config"
)

var RedisClient *redis.Client

// NewRedisClient 传入dbName， 比如classroom将key放在db0中，user放在db1中
func NewRedisClient(dbName int) *redis.Client {
	conf := config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       dbName,
	})
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
	return client
}
