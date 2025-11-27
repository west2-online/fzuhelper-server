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
	"fmt"

	"github.com/west2-online/fzuhelper-server/pkg/base/environment"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func (c *CacheUser) SetInvitationCodeCache(ctx context.Context, key string, code string) error {
	if environment.IsTestEnvironment() {
		return nil
	}
	if err := c.client.Set(ctx, key, code, constants.UserInvitationCodeKeyExpire).Err(); err != nil {
		return fmt.Errorf("dal.SetInvitationCodeCache: Set cache failed: %w", err)
	}
	return nil
}

func (c *CacheUser) SetCodeStuIdMappingCache(ctx context.Context, key, stuId string) error {
	if environment.IsTestEnvironment() {
		return nil
	}
	if err := c.client.Set(ctx, key, stuId, constants.UserInvitationCodeKeyExpire).Err(); err != nil {
		return fmt.Errorf("dal.SetCodeStuIdMappingCache: Set cache failed: %w", err)
	}
	return nil
}

func (c *CacheUser) RemoveCodeStuIdMappingCache(ctx context.Context, key string) error {
	if environment.IsTestEnvironment() {
		return nil
	}
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("dal.RemoveCodeStuIdMappingCache: Delete cache failed: %w", err)
	}
	return nil
}

func (c *CacheUser) SetUserFriendCache(ctx context.Context, stuId, friendId string) error {
	if environment.IsTestEnvironment() {
		return nil
	}
	pipe := c.client.Pipeline()
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	userFriendKey_ := fmt.Sprintf("user_friends:%v", friendId)
	pipe.SAdd(ctx, userFriendKey, friendId)
	pipe.SAdd(ctx, userFriendKey_, stuId)
	pipe.Expire(ctx, userFriendKey, constants.UserFriendKeyExpire)
	pipe.Expire(ctx, userFriendKey_, constants.UserFriendKeyExpire)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("dal.SetInvitationCodeCache: Set cache failed: %w", err)
	}
	return nil
}

func (c *CacheUser) SetUserFriendListCache(ctx context.Context, stuId string, friendIds []string) error {
	if environment.IsTestEnvironment() {
		return nil
	}
	pipe := c.client.Pipeline()
	userFriendKey := fmt.Sprintf("user_friends:%v", stuId)
	for _, id := range friendIds {
		pipe.SAdd(ctx, userFriendKey, id)
	}
	pipe.Expire(ctx, userFriendKey, constants.UserFriendKeyExpire)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("dal.SetUserFriendListCache: Set cache failed: %w", err)
	}
	return nil
}
