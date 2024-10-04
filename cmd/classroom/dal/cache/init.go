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
		// 如果redis服务启动失败，直接exit
		logger.Fatalf("cache.Init failed, err is %v", err)
	}
	RedisClient = redisClient
}
