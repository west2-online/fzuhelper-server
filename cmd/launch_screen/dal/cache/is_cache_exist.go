package cache

import (
	"context"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func IsLaunchScreenCacheExist(ctx context.Context, key string) bool {
	return RedisClient.Exists(ctx, key).Val() == 1
}

func IsLastLaunchScreenIdCacheExist(ctx context.Context) bool {
	return RedisClient.Exists(ctx, constants.LastLaunchScreenIdKey).Val() == 1
}
