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

// CreateRelation 目前进行一条关系进行双向插入
// 好友关系可能会有 删除-建立-删除-建立的过程，这边先进行一次status的更新尝试
func (c *DBUser) CreateRelation(ctx context.Context, relation []*model.FollowRelation) error {
	err := c.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Table(constants.UserRelationTableName).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "follower_id"}, {Name: "followed_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"status", "updated_at"}),
			}).Create(&relation).
			Error
	})
	if err != nil {
		logger.Errorf("dal.CreateRelation error: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.CreateRelation error: %v", err)
	}
	return nil
}
