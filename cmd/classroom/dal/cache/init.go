package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/west2-online/fzuhelper-server/pkg/client"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

var RedisClient *redis.Client

func Init() {
	RedisClient = client.NewRedisClient(constants.RedisDBEmptyRoom)
}
