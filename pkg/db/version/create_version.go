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

package version

import (
	"context"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *DBVersion) CreateVersion(ctx context.Context, version *model.Visit) error {
	if err := c.client.WithContext(ctx).Table(constants.VisitTableName).Create(&version).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.CreateVersion error: %v", err)
	}
	return nil
}

func (c *DBVersion) UpdateVersion(ctx context.Context, version *model.Visit) error {
	if err := c.client.WithContext(ctx).Table(constants.VisitTableName).Where("date = ?", version.Date).UpdateColumn("visits", version.Visits).Error; err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.UpdateVersion error: %v", err)
	}
	return nil
}
