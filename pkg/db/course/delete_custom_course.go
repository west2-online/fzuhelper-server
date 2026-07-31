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

// DeleteCustomCourse 删除自定义课程（软删除）
func (c *DBCourse) DeleteCustomCourse(ctx context.Context, stuId, term, courseId string) error {
	return c.client.WithContext(ctx).
		Where("stu_id = ? AND term = ? AND course_id = ?", stuId, term, courseId).
		Delete(&model.UserCustomCourse{}).Error
}
