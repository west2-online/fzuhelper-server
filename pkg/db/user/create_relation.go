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
	"gorm.io/gorm/clause"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// CreateRelation 双向插入好友关系
// 好友关系支持 删除-建立-删除-建立 的循环，冲突时通过 OnConflict 恢复 status
func (c *DBUser) CreateRelation(ctx context.Context, relation []*model.FollowRelation) error {
	err := c.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在事务中查询并分配 order_seq，确保在高并发场景下也能正确递增
		for _, r := range relation {
			var nextSeq int64
			if err := tx.Table(constants.UserRelationTableName).
				Where("follower_id = ? AND status = ?", r.FollowerId, constants.RelationOKStatus).
				Select("COALESCE(MAX(order_seq), 0) + 1").
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Scan(&nextSeq).Error; err != nil {
				return err
			}
			r.OrderSeq = nextSeq
		}

		return tx.Table(constants.UserRelationTableName).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "follower_id"}, {Name: "followed_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"status", "updated_at", "order_seq"}),
			}).Create(&relation).
			Error
	})
	if err != nil {
		logger.Errorf("dal.CreateRelation error: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.CreateRelation error: %v", err)
	}
	return nil
}
