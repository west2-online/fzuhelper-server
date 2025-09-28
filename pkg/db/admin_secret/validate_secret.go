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

package admin_secret

import (
	"context"
	"fmt"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// ValidateSecret 验证模块密钥是否有效
func (c *DBAdminSecret) ValidateSecret(ctx context.Context, moduleName, secretKey string) error {
	var secret model.AdminSecret

	err := c.client.WithContext(ctx).
		Table(constants.AdminSecretTableName).
		Where("module_name = ? AND secret_key = ?", moduleName, secretKey).
		First(&secret).Error
	if err != nil {
		return errno.NewErrNo(errno.AuthErrorCode, "invalid admin secret")
	}

	return nil
}

func (c *DBAdminSecret) GetSecretsByModule(ctx context.Context, moduleName string) ([]*model.AdminSecret, error) {
	var secrets []*model.AdminSecret

	err := c.client.WithContext(ctx).
		Table(constants.AdminSecretTableName).
		Where("module_name = ?", moduleName).
		Find(&secrets).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.GetSecretsByModule error: %v", err))
	}

	return secrets, nil
}
