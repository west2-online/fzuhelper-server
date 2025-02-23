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
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

func (c *DBCourse) GetUserTermCourseByStuIdAndTerm(ctx context.Context, stuId string, term string) (*model.UserCourse, error) {
	userCourseModel := new(model.UserCourse)
	if err := c.client.WithContext(ctx).Table(constants.CourseTableName).Where("stu_id = ? and term = ?", stuId, term).First(userCourseModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("dal.GetUserTermCourseByStuIdAndTerm error: %w", err)
	}
	return userCourseModel, nil
}

func (c *DBCourse) GetUserTermCourseSha256ByStuIdAndTerm(ctx context.Context, stuId string, term string) (*model.UserCourse, error) {
	userCourseModel := new(model.UserCourse)
	if err := c.client.WithContext(ctx).
		Table(constants.CourseTableName).
		Select("id", "term_courses_sha256").
		Where("stu_id = ? and term = ?", stuId, term).
		First(userCourseModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("dal.GetUserTermCourseSha256ByStuIdAndTerm error: %w", err)
	}
	return userCourseModel, nil
}
