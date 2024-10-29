package cache

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func SetLaunchScreenCache(ctx context.Context, key string, pictureIdList *[]int64) error {
	pictureIdListJson, err := sonic.Marshal(pictureIdList)
	if err != nil {
		return fmt.Errorf("dal.SetLaunchScreenCache: Marshal pictureIdList failed: %w", err)
	}
	if err = RedisClient.Set(ctx, key, pictureIdListJson, constants.LaunchScreenKeyExpire).Err(); err != nil {
		return fmt.Errorf("dal.SetLaunchScreenCache: Set pictureIdList cache failed: %w", err)
	}
	return nil
}

func SetLastLaunchScreenIdCache(ctx context.Context, id int64) error {
	if err := RedisClient.Set(ctx, constants.LastLaunchScreenIdKey, id, constants.LaunchScreenKeyExpire).Err(); err != nil {
		return fmt.Errorf("dal.SetTotalLaunchScreenCountCache failed: %w", err)
	}
	return nil
}
