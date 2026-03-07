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

// ReorderFriendList 按照传入的 friendIds 列表顺序，依次从小到大赋值 order_seq。
// 对于不在 friendIds 中但实际存在的好友（一般是脏数据），按原 order_seq ASC 顺序追加到末尾。
func (c *DBUser) ReorderFriendList(ctx context.Context, stuId string, friendIds []string) error {
	err := c.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 查询当前所有好友，按 order_seq ASC 保持原有顺序
		var allFriends []struct {
			FollowedId string
		}
		if err := tx.Table(constants.UserRelationTableName).
			Where("follower_id = ? AND status = ?", stuId, constants.RelationOKStatus).
			Select("followed_id").
			Order("order_seq ASC, created_at ASC").
			Find(&allFriends).Error; err != nil {
			return err
		}

		// 2. 构建请求中 friendIds 的集合
		requestedSet := make(map[string]struct{}, len(friendIds))
		for _, id := range friendIds {
			requestedSet[id] = struct{}{}
		}

		// 3. 计算补集：在 DB 中存在但不在请求列表中的好友，保持原有顺序
		var remainder []string
		for _, f := range allFriends {
			if _, ok := requestedSet[f.FollowedId]; !ok {
				remainder = append(remainder, f.FollowedId)
			}
		}

		// 4. 合并：请求列表在前，补集在后
		merged := make([]string, 0, len(friendIds)+len(remainder))
		merged = append(merged, friendIds...)
		merged = append(merged, remainder...)

		// 5. 依次赋值 order_seq
		for i, friendId := range merged {
			orderSeq := int64(i + 1)
			if err := tx.Table(constants.UserRelationTableName).
				Where("follower_id = ? AND followed_id = ? AND status = ?", stuId, friendId, constants.RelationOKStatus).
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
