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
	"math"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *DBToolbox) GetToolboxConfigs(ctx context.Context) ([]*model.ToolboxConfig, error) {
	toolboxConfigs := make([]*model.ToolboxConfig, 0)
	if err := c.client.WithContext(ctx).
		Model(&model.ToolboxConfig{}).
		Table(constants.ToolboxConfigTableName).
		Find(&toolboxConfigs).Error; err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.GetToolboxConfigs error: %v", err))
	}
	return toolboxConfigs, nil
}

func (c *DBToolbox) GetToolboxConfigByID(ctx context.Context, id int64) (*model.ToolboxConfig, error) {
	config := new(model.ToolboxConfig)
	err := c.client.WithContext(ctx).
		Model(&model.ToolboxConfig{}).
		Table(constants.ToolboxConfigTableName).
		Where("id = ?", id).
		First(config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errno.NewErrNo(errno.BizNotExist, "toolbox config not found")
	}
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.GetToolboxConfigByID error: %v", err))
	}
	return config, nil
}

func (c *DBToolbox) ListToolboxConfigs(ctx context.Context, pageNum, pageSize int) ([]*model.ToolboxConfig, int64, error) {
	if pageNum <= 0 || pageSize <= 0 {
		return nil, 0, errno.NewErrNo(errno.ParamErrorCode, "page_num and page_size must be positive")
	}
	if pageNum-1 > math.MaxInt/pageSize {
		return nil, 0, errno.NewErrNo(errno.ParamErrorCode, "page offset is too large")
	}

	var total int64
	if err := c.client.WithContext(ctx).
		Model(&model.ToolboxConfig{}).
		Table(constants.ToolboxConfigTableName).
		Count(&total).Error; err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.ListToolboxConfigs count error: %v", err))
	}

	toolboxConfigs := make([]*model.ToolboxConfig, 0)
	offset := (pageNum - 1) * pageSize
	if err := c.client.WithContext(ctx).
		Model(&model.ToolboxConfig{}).
		Table(constants.ToolboxConfigTableName).
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&toolboxConfigs).Error; err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.ListToolboxConfigs error: %v", err))
	}

	return toolboxConfigs, total, nil
}
