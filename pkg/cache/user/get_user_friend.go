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
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

func (c *CacheUser) GetInvitationCodeCache(ctx context.Context, key string) (code string, createdAt int64, err error) {
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", -1, fmt.Errorf("dal.GetInvitationCodeCache: GetInvitationCode cache failed: %w", err)
	}
	part := strings.Split(value, "-")
	if len(part) != constants.UserInvitationCodeCachePartLength {
		return "", -1, fmt.Errorf("dal.GetInvitationCodeCache: GetInvitationCode cache failed: invaild code format")
	}
	code = part[0]
	createdAt, err = strconv.ParseInt(part[1], 10, 64)
	if err != nil {
		return "", -1, fmt.Errorf("dal.GetInvitationCodeCache: GetInvitationCode cache failed: %w", err)
	}
	return code, createdAt, nil
}

func (c *CacheUser) GetCodeStuIdMappingCache(ctx context.Context, key string) (stuId string, err error) {
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("dal.GetCodeStuIdMappingCache: GetStuIdCodeMapping cache failed: %w", err)
	}
	return value, nil
}

func (c *CacheUser) GetUserFriendCache(ctx context.Context, key string) (friendList []*model.UserFriend, err error) {
	results, err := c.client.ZRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("dal.GetCodeStuIdMappingCache: GetStuIdCodeMapping cache failed: %w", err)
	}
	friendList = make([]*model.UserFriend, 0)
	for _, z := range results {
		if friendId, ok := z.Member.(string); ok {
			friendList = append(friendList, &model.UserFriend{
				FriendId:  friendId,
				UpdatedAt: time.Unix(int64(z.Score), 0),
			})
		}
	}
	return friendList, nil
}

func (c *CacheUser) IsFriendCache(ctx context.Context, stuId, friendId string) (bool, error) {
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	score, err := c.client.ZScore(ctx, userFriendKey, friendId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("IsFriendCache: check failed: %w", err)
	}
	return score > 0, nil
}
