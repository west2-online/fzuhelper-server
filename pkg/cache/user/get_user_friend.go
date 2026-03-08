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
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

// userFriendCacheValue is the JSON-serialized value stored in the Hash for each friend.
type userFriendCacheValue struct {
	OrderSeq  int64 `json:"order_seq"`
	CreatedAt int64 `json:"created_at"` // unix timestamp
}

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
	results, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("dal.GetUserFriendCache: Get cache failed: %w", err)
	}
	if len(results) == 0 {
		return nil, nil
	}
	friendList = make([]*model.UserFriend, 0, len(results))
	for friendId, raw := range results {
		var val userFriendCacheValue
		if err := json.Unmarshal([]byte(raw), &val); err != nil {
			return nil, fmt.Errorf("dal.GetUserFriendCache: unmarshal failed for %s: %w", friendId, err)
		}
		friendList = append(friendList, &model.UserFriend{
			FriendId:  friendId,
			OrderSeq:  val.OrderSeq,
			CreatedAt: time.Unix(val.CreatedAt, 0),
		})
	}
	// Hash is unordered, sort by OrderSeq ASC, then CreatedAt (updated_at) ASC
	sort.Slice(friendList, func(i, j int) bool {
		if friendList[i].OrderSeq != friendList[j].OrderSeq {
			return friendList[i].OrderSeq < friendList[j].OrderSeq
		}
		return friendList[i].CreatedAt.Before(friendList[j].CreatedAt)
	})
	return friendList, nil
}

func (c *CacheUser) IsFriendCache(ctx context.Context, stuId, friendId string) (bool, error) {
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	exists, err := c.client.HExists(ctx, userFriendKey, friendId).Result()
	if err != nil {
		return false, fmt.Errorf("IsFriendCache: check failed: %w", err)
	}
	return exists, nil
}
