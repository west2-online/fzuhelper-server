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

package friend_config

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/clause"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// UpsertFriendConfig 插入或更新好友配置
func (c *DBFriendConfig) UpsertFriendConfig(ctx context.Context, config *model.FriendConfig) (*model.FriendConfig, error) {
	config.UpdatedAt = time.Now()
	err := c.client.WithContext(ctx).Table(constants.FriendConfigTableName).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "config_key"},
				{Name: "student_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"value",
				"updated_at",
			}),
		}).Create(config).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.UpsertFriendConfig upsert error: %v", err))
	}

	return config, nil
}
