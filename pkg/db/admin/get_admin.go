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

package admin

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *DBAdmin) GetAdminByFeishuUserId(ctx context.Context, userId string) (bool, *model.AdminUser, error) {
	adminUser := new(model.AdminUser)
	err := c.client.WithContext(ctx).
		Table(constants.AdminUserTableName).
		Where("feishu_user_id = ? AND enabled = ?", userId, true).
		First(adminUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetAdminByFeishuUserId error: %v", err)
	}
	return true, adminUser, nil
}
