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
	"errors"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *DBVersion) GetVersion(ctx context.Context, date string) (bool, int64, error) {
	versionModel := new(model.Visit)
	err := c.client.WithContext(ctx).Table(constants.VisitTableName).Where("date = ?", date).First(versionModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, 0, nil
		}
		return false, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "db.GetVersion error : %v", err)
	}
	return true, versionModel.Visits, nil
}

func (c *DBVersion) GetVersionList(ctx context.Context) ([]*model.Visit, error) {
	var versions []*model.Visit
	err := c.client.WithContext(ctx).Table(constants.VisitTableName).Order("created_at desc").Find(&versions).Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "db.GetVersionList error: %v", err)
	}

	return versions, nil
}
