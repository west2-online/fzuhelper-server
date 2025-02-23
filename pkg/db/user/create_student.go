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

package user

import (
	"context"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (c *DBUser) CreateStudent(ctx context.Context, userModel *model.Student) error {
	if err := c.client.WithContext(ctx).Table(constants.UserTableName).Create(&userModel).Error; err != nil {
		logger.Errorf("dal.CreateStudent error: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.CreateStudent error: %v", err)
	}
	return nil
}

func (c *DBUser) UpdateStudent(ctx context.Context, userModel *model.Student) error {
	if err := c.client.WithContext(ctx).Table(constants.UserTableName).Where("stu_id = ?", userModel.StuId).Omit("created_at").Save(userModel).Error; err != nil {
		logger.Errorf("dal.CreateStudent error: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.CreateStudent error: %v", err)
	}
	return nil
}
