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

package version

import (
	"context"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *CacheVersion) AddVisit(ctx context.Context, date string) error {
	// 增加访问量
	count, err := c.client.Incr(ctx, date).Result()
	if err != nil {
		return errno.Errorf(errno.InternalRedisErrorCode, "version.AddVisit error: %v", err)
	}

	// 如果 count == 1，说明这个 key 是新创建的，设置过期时间为 24 小时
	if count == 1 {
		err = c.client.Expire(ctx, date, constants.VisitExpire).Err()
		if err != nil {
			return errno.Errorf(errno.InternalRedisErrorCode, "version.AddVisit Expire error: %v", err)
		}
	}

	return nil
}

func (c *CacheVersion) CreateVisitKey(ctx context.Context, date string) error {
	// 尝试设置 key 仅在它不存在时
	set, err := c.client.SetNX(ctx, date, "1", constants.VisitExpire).Result()
	if err != nil {
		return errno.Errorf(errno.InternalRedisErrorCode, "version.CreateVisitKey error: %v", err)
	}

	// 如果 set 为 false，说明 key 已经存在，无需处理
	if !set {
		return nil
	}

	return nil
}
