package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/config"
)

var RedisClient *redis.Client

func Init() {
	conf := config.Config.Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.RedisHost + ":" + conf.RedisPort,
		Password: conf.RedisPassword,
		DB:       conf.RedisDbName,
	})
	_, err := RedisClient.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
}
