package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/west2-online/fzuhelper-server/pkg/client"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

var RedisClient *redis.Client

func Init() {
	redisClient, err := client.NewRedisClient(constants.RedisDBEmptyRoom)
	if err != nil {
		logger.LoggerObj.Fatalf("cache.Init failed, err is %v", err)
		panic(err)
	}
	RedisClient = redisClient
}
