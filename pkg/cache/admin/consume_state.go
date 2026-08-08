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

package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func (c *CacheAdmin) ConsumeOAuthStateCache(ctx context.Context, state string) (returnTo string, exists bool, err error) {
	key := constants.AdminOAuthStateKeyPrefix + state
	returnTo, err = c.client.GetDel(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("dal.ConsumeOAuthStateCache: get and delete cache failed: %w", err)
	}
	return returnTo, true, nil
}
