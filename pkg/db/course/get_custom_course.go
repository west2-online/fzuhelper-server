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

// GetCustomCourses 获取用户指定学期的自定义课程列表
func (c *DBCourse) GetCustomCourses(ctx context.Context, stuId, term string) ([]*model.UserCustomCourse, error) {
	var courses []*model.UserCustomCourse
	err := c.client.WithContext(ctx).
		Where("stu_id = ? AND term = ? AND deleted_at IS NULL", stuId, term).
		Find(&courses).Error
	return courses, err
}

// GetCustomCourseByID 根据 courseId 获取单个自定义课程
func (c *DBCourse) GetCustomCourseByID(ctx context.Context, stuId, term, courseId string) (*model.UserCustomCourse, error) {
	var course model.UserCustomCourse
	err := c.client.WithContext(ctx).
		Where("stu_id = ? AND term = ? AND course_id = ? AND deleted_at IS NULL", stuId, term, courseId).
		First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}
