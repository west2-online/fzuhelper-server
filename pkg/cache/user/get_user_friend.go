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

package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (c *CacheUser) GetInvitationCodeCache(ctx context.Context, key string) (code string, err error) {
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("dal.GetInvitationCodeCache: GetInvitationCode cache failed: %w", err)
	}
	return value, nil
}

func (c *CacheUser) GetCodeStuIdMappingCache(ctx context.Context, key string) (stuId string, err error) {
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("dal.GetCodeStuIdMappingCache: GetStuIdCodeMapping cache failed: %w", err)
	}
	return value, nil
}

func (c *CacheUser) GetUserFriendCache(ctx context.Context, key string) (friendIds []string, err error) {
	friendIds, err = c.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("dal.GetCodeStuIdMappingCache: GetStuIdCodeMapping cache failed: %w", err)
	}
	return friendIds, nil
}

func (c *CacheUser) IsFriendCache(ctx context.Context, stuId, friendId string) (bool, error) {
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	exists, err := c.client.SIsMember(ctx, userFriendKey, friendId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("IsFriendCache: check failed: %w", err)
	}
	return exists, nil
}
