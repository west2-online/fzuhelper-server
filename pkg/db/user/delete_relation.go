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

// DeleteRelation 删除关系
// 但是实际上是更新status
func (c *DBUser) DeleteRelation(ctx context.Context, followerId, followedId string) error {
	err := c.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Table(constants.UserRelationTableName).
			Where("follower_id = ? and followed_id = ?", followerId, followedId).
			Update("status", constants.RelationDeletedStatus).
			Error
		if err != nil {
			return err
		}
		err = tx.Table(constants.UserRelationTableName).
			Where("follower_id = ? and followed_id = ?", followedId, followerId).
			Update("status", constants.RelationDeletedStatus).
			Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Errorf("dal.DeleteRelation error: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.DeleteRelation error: %v", err)
	}
	return nil
}
