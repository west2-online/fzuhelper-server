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

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *DBToolbox) CreateToolboxConfig(ctx context.Context, config *model.ToolboxConfig) error {
	err := c.client.WithContext(ctx).Table(constants.ToolboxConfigTableName).Create(config).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errno.NewErrNo(errno.BizLogicCode, "toolbox config already exists")
	}
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.CreateToolboxConfig error: %v", err))
	}
	return nil
}
