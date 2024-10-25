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

package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserCourse struct {
	Id                int64
	StuId             string
	Term              string
	TermCourses       string
	TermCoursesSha256 string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `sql:"index"`
}

func CreateUserTermCourse(ctx context.Context, userCourseModel *UserCourse) (*UserCourse, error) {
	if err := DB.WithContext(ctx).Create(userCourseModel).Error; err != nil {
		return nil, fmt.Errorf("dal.CreateUserTermCourse error: %v", err)
	}
	return userCourseModel, nil
}

func GetUserTermCourseByStuIdAndTerm(ctx context.Context, stuId string, term string) (*UserCourse, error) {
	userCourseModel := new(UserCourse)
	if err := DB.WithContext(ctx).Where("stu_id = ? and term = ?", stuId, term).First(userCourseModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("dal.GetUserTermCourseByStuIdAndTerm error: %v", err)
	}
	return userCourseModel, nil
}

func GetUserTermCourseSha256ByStuIdAndTerm(ctx context.Context, stuId string, term string) (*UserCourse, error) {
	userCourseModel := new(UserCourse)
	if err := DB.WithContext(ctx).Select("id", "term_courses_sha256").Where("stu_id = ? and term = ?", stuId, term).First(userCourseModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("dal.GetUserTermCourseSha256ByStuIdAndTerm error: %v", err)
	}
	return userCourseModel, nil
}

func UpdateUserTermCourse(ctx context.Context, userCourseModel *UserCourse) (*UserCourse, error) {
	if err := DB.WithContext(ctx).Model(userCourseModel).Updates(userCourseModel).Error; err != nil {
		return nil, fmt.Errorf("dal.UpdateUserTermCourse error: %v", err)
	}
	return userCourseModel, nil
}
