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
	"time"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// CreateRelation 目前进行一条关系进行双向插入
// 好友关系可能会有 删除-建立-删除-建立的过程，这边先进行一次status的更新尝试
func (c *DBUser) CreateRelation(ctx context.Context, followerId, followedId string) error {
	relation := []*model.FollowRelation{
		{
			FollowedId: followedId,
			FollowerId: followerId,
		},
		{
			FollowedId: followerId,
			FollowerId: followedId,
		},
	}
	err := c.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Table(constants.UserRelationTableName).
			Where("follower_id = ? AND followed_id = ? AND status = 1", followerId, followedId).
			Updates(map[string]interface{}{
				"status":     constants.RelationOKStatus,
				"updated_at": time.Now(),
			})
		if result.Error != nil {
			return result.Error
		}
		// 记录不存在则创建
		if result.RowsAffected == 0 {
			err := tx.Table(constants.UserRelationTableName).
				Create(&relation).
				Error
			if err != nil {
				return err
			}
			return nil
		}
		// 存在则相应地进行双向关系的更新
		result = tx.Table(constants.UserRelationTableName).
			Where("follower_id = ? AND followed_id = ? AND status = 1", followedId, followerId).
			Updates(map[string]interface{}{
				"status":     constants.RelationOKStatus,
				"updated_at": time.Now(),
			})
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		logger.Errorf("dal.CreateRelation error: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.CreateRelation error: %v", err)
	}
	return nil
}
