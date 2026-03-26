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
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// CreateRelation 双向插入好友关系
// 每次添加好友都插入新记录，删除时通过 deleted_at 软删除
func (c *DBUser) CreateRelation(ctx context.Context, relation []*model.FollowRelation) error {
	err := c.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Table(constants.UserRelationTableName).
			Create(&relation).
			Error
	})
	if err != nil {
		logger.Errorf("dal.CreateRelation error: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.CreateRelation error: %v", err)
	}
	return nil
}
