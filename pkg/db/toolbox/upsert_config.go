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

package toolbox

import (
	"context"
	"fmt"

	"gorm.io/gorm/clause"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// UpsertToolboxConfig 插入或更新工具箱配置
func (c *DBToolbox) UpsertToolboxConfig(ctx context.Context, config *model.ToolboxConfig) (*model.ToolboxConfig, error) {
	// 使用GORM的OnConflict实现upsert
	err := c.client.WithContext(ctx).Table(constants.ToolboxConfigTableName).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "tool_id"},
				{Name: "student_id"},
				{Name: "platform"},
				{Name: "version"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"visible",
				"name",
				"icon",
				"type",
				"message",
				"extra",
				"updated_at",
			}),
		}).Create(config).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.UpsertToolboxConfig upsert error: %v", err))
	}

	return config, nil
}
