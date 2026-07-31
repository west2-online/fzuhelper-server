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
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *DBToolbox) UpdateToolboxConfig(ctx context.Context, id int64, config *model.ToolboxConfig) (*model.ToolboxConfig, error) {
	updates := map[string]any{
		"tool_id":    config.ToolID,
		"visible":    config.Visible,
		"name":       config.Name,
		"icon":       config.Icon,
		"type":       config.Type,
		"message":    config.Message,
		"extra":      config.Extra,
		"student_id": config.StudentID,
		"platform":   config.Platform,
		"version":    config.Version,
		"updated_at": time.Now(),
	}
	result := c.client.WithContext(ctx).
		Model(&model.ToolboxConfig{}).
		Table(constants.ToolboxConfigTableName).
		Where("id = ?", id).
		Updates(updates)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, errno.NewErrNo(errno.BizLogicCode, "toolbox config already exists")
	}
	if result.Error != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.UpdateToolboxConfig error: %v", result.Error))
	}
	return c.GetToolboxConfigByID(ctx, id)
}
