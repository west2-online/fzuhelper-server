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

package course

import (
	"context"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *DBCourse) UpdateAutoAdjustCourse(ctx context.Context, id int64, updates map[string]any) error {
	result := c.client.WithContext(ctx).
		Table(constants.AutoAdjustCourseTableName).
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "dal.UpdateAutoAdjustCourse update error: id=%d, %v", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return errno.Errorf(errno.BizNotExist, "dal.UpdateAutoAdjustCourse: no record found with id=%d", id)
	}
	return nil
}
