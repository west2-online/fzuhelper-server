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
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

// DeleteCustomCourse 删除自定义课程（软删除），返回受影响的行数
func (c *DBCourse) DeleteCustomCourse(ctx context.Context, stuId string, id int64) (int64, error) {
	result := c.client.WithContext(ctx).
		Table(constants.UserCustomCourseTableName).
		Where("stu_id = ? AND id = ?", stuId, id).
		Delete(&model.UserCustomCourse{})
	return result.RowsAffected, result.Error
}
