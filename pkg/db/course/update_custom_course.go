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

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

// UpdateCustomCourse 更新自定义课程，返回受影响的行数
func (c *DBCourse) UpdateCustomCourse(ctx context.Context, stuId string, id int64, updates map[string]interface{}) (int64, error) {
	result := c.client.WithContext(ctx).
		Model(&model.UserCustomCourse{}).
		Where("stu_id = ? AND id = ? AND active_flag = 1", stuId, id).
		Updates(updates)
	return result.RowsAffected, result.Error
}
