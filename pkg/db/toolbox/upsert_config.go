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

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// UpsertToolboxConfig 插入或更新工具箱配置
// 如果存在相同的 tool_id + student_id + platform + version 组合，则更新；否则插入
func (c *DBToolbox) UpsertToolboxConfig(ctx context.Context, config *model.ToolboxConfig) (*model.ToolboxConfig, error) {
	var existingConfig model.ToolboxConfig

	// 构建查询条件
	query := c.client.WithContext(ctx).Table(constants.ToolboxConfigTableName).
		Where("tool_id = ?", config.ToolID)

	// 添加可选的查询条件
	if config.StudentID != "" {
		query = query.Where("student_id = ?", config.StudentID)
	} else {
		query = query.Where("student_id IS NULL OR student_id = ''")
	}

	if config.Version > 0 {
		query = query.Where("version = ?", config.Version)
	} else {
		query = query.Where("version = 0 OR version IS NULL")
	}

	if config.Platform != "" {
		query = query.Where("platform = ?", config.Platform)
	} else {
		query = query.Where("platform IS NULL OR platform = ''")
	}

	// 查找是否存在匹配的记录
	err := query.First(&existingConfig).Error

	if err == nil {
		// 记录存在，更新
		config.Id = existingConfig.Id
		config.CreatedAt = existingConfig.CreatedAt

		if err := c.client.WithContext(ctx).Table(constants.ToolboxConfigTableName).Save(config).Error; err != nil {
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.UpsertToolboxConfig update error: %v", err))
		}
		return config, nil
	} else {
		// 记录不存在，插入新记录
		if err := c.client.WithContext(ctx).Table(constants.ToolboxConfigTableName).Create(config).Error; err != nil {
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.UpsertToolboxConfig create error: %v", err))
		}
		return config, nil
	}
}
