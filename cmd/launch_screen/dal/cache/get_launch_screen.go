package cache

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func GetLaunchScreenCache(ctx context.Context, key string) (pictureIdList []int64, err error) {
	data, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("dal.GetLaunchScreenCache: Get pictureIdList cache failed: %w", err)
	}

	if err = sonic.Unmarshal([]byte(data), &pictureIdList); err != nil {
		return nil, fmt.Errorf("dal.GetLaunchScreenCache: Unmarshal pictureIdList failed: %w", err)
	}

	return pictureIdList, nil
}

func GetLastLaunchScreenIdCache(ctx context.Context) (int64, error) {
	id, err := RedisClient.Get(ctx, constants.LastLaunchScreenIdKey).Int64()
	if err != nil {
		return -1, fmt.Errorf("dal.GetLaunchScreenCache: Get pictureIdList cache failed: %w", err)
	}

	return id, nil
}
