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

package course

import (
	"context"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *CacheCourse) SetAndGetRefreshCount(ctx context.Context, key string) (int64, error) {
	refreshKey := key + constants.RefreshSuffixKey

	// 使用 Redis INCR 自增 refreshCount
	count, err := c.client.Incr(ctx, refreshKey).Result()
	if err != nil {
		return -1, errno.Errorf(errno.InternalRedisErrorCode, "course.SetAndGetRefreshCount: incr failed")
	}

	// 如果 count == 1，说明这个 key 是新创建，要加expire
	if count == 1 {
		if err = c.client.Expire(ctx, refreshKey, constants.RefreshCountExpire).Err(); err != nil {
			return -1, errno.Errorf(errno.InternalRedisErrorCode, "course.SetAndGetRefreshCount: set expire failed")
		}
	}

	return count, nil
}
