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

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// ReorderFriendList 按照传入的 friendIds 列表顺序，依次从大到小赋值 order_seq（越大越靠前）。
// 对于不在 friendIds 中但实际存在的好友（脏数据），将其 order_seq 置为 0
func (c *DBUser) ReorderFriendList(ctx context.Context, stuId string, friendIds []string) error {
	err := c.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 查询当前所有好友
		var allFriends []struct {
			FollowedId string
		}
		if err := tx.Table(constants.UserRelationTableName).
			Where("follower_id = ? AND deleted_at IS NULL", stuId).
			Select("followed_id").
			Find(&allFriends).Error; err != nil {
			return err
		}

		// 2. 构建请求中 friendIds 的集合
		requestedSet := make(map[string]struct{}, len(friendIds))
		for _, id := range friendIds {
			requestedSet[id] = struct{}{}
		}

		// 3. 将不在请求列表中的好友（脏数据）的 order_seq 置为 0
		for _, f := range allFriends {
			if _, ok := requestedSet[f.FollowedId]; !ok {
				if err := tx.Table(constants.UserRelationTableName).
					Where("follower_id = ? AND followed_id = ? AND deleted_at IS NULL", stuId, f.FollowedId).
					Update("order_seq", 0).Error; err != nil {
					return err
				}
			}
		}

		// 4. 按传入顺序从大到小赋值 order_seq（列表中第一个元素 order_seq 最大，排在最前）
		for i, friendId := range friendIds {
			orderSeq := int64(len(friendIds) - i)
			if err := tx.Table(constants.UserRelationTableName).
				Where("follower_id = ? AND followed_id = ? AND deleted_at IS NULL", stuId, friendId).
				Update("order_seq", orderSeq).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.Errorf("dal.ReorderFriendList error: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.ReorderFriendList error: %v", err)
	}
	return nil
}
