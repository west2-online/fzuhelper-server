/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cache

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func IsLaunchScreenCacheExist(ctx context.Context, key string) bool {
	return RedisClient.Exists(ctx, key).Val() == 1
}

func IsLastLaunchScreenIdCacheExist(ctx context.Context) bool {
	return RedisClient.Exists(ctx, constants.LastLaunchScreenIdKey).Val() == 1
}

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

func SetLastLaunchScreenIdCache(ctx context.Context, id int64) error {
	if err := RedisClient.Set(ctx, constants.LastLaunchScreenIdKey, id, constants.LaunchScreenKeyExpire).Err(); err != nil {
		return fmt.Errorf("dal.SetTotalLaunchScreenCountCache failed: %w", err)
	}
	return nil
}

func GetLastLaunchScreenIdCache(ctx context.Context) (int64, error) {
	id, err := RedisClient.Get(ctx, constants.LastLaunchScreenIdKey).Int64()
	if err != nil {
		return -1, fmt.Errorf("dal.GetLaunchScreenCache: Get pictureIdList cache failed: %w", err)
	}

	return id, nil
}
