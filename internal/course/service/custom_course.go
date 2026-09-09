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

package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
)

func (s *CourseService) GetCustomCourses(ctx context.Context, stuID, term string) ([]*kitexModel.CustomCourse, error) {
	key := s.cache.Course.CustomCourseKey(stuID, term)
	if s.cache.IsKeyExist(s.ctx, key) {
		items, err := s.cache.Course.GetCustomCoursesCache(s.ctx, key)
		if err != nil {
			return nil, fmt.Errorf("CourseService.GetCustomCourses: get custom courses cache failed: %w", err)
		}
		return items, nil
	}

	courses, err := s.db.Course.GetCustomCourses(ctx, stuID, term)
	if err != nil {
		return nil, fmt.Errorf("CourseService.GetCustomCourses: get custom courses failed: %w", err)
	}
	items := pack.BuildCustomCourseItems(courses)
	// async put custom courses to cache
	s.taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
		return s.cache.Course.SetCustomCoursesCache(s.ctx, key, items)
	}})
	return items, nil
}

// refreshCustomCourseCache 在自定义课程变更后，从 db 重读全量回填缓存，与调课缓存刷新策略一致
func (s *CourseService) refreshCustomCourseCache(stuID, term string) {
	key := s.cache.Course.CustomCourseKey(stuID, term)
	s.taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
		courses, err := s.db.Course.GetCustomCourses(s.ctx, stuID, term)
		if err != nil {
			return err
		}
		return s.cache.Course.SetCustomCoursesCache(s.ctx, key, pack.BuildCustomCourseItems(courses))
	}})
}

func (s *CourseService) UpsertCustomCourse(ctx context.Context, stuID string, req *course.UpsertCustomCourseRequest) (string, error) {
	item := req.Course
	if item.Id != nil && *item.Id != "" {
		return s.updateCustomCourse(ctx, stuID, *item.Id, item)
	}

	id, err := s.sf.NextVal()
	if err != nil {
		return "", err
	}
	customCourse := &model.UserCustomCourse{
		Id:         id,
		StuId:      stuID,
		Term:       req.Term,
		Name:       item.Name,
		Teacher:    getStringValue(item.Teacher),
		Location:   item.Location,
		StartClass: int(item.StartClass),
		EndClass:   int(item.EndClass),
		StartWeek:  int(item.StartWeek),
		EndWeek:    int(item.EndWeek),
		Weekday:    int(item.Weekday),
		IsSingle:   item.Single,
		IsDouble:   item.Double,
		Color:      getStringValueWithDefault(item.Color, "#FF5733"),
		Remark:     getStringValue(item.Remark),
	}
	if _, err := s.db.Course.CreateCustomCourse(ctx, customCourse); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", errno.BizError.WithMessage("重复添加课程！")
		}
		return "", err
	}
	s.refreshCustomCourseCache(stuID, req.Term)
	return strconv.FormatInt(customCourse.Id, 10), nil
}

func (s *CourseService) updateCustomCourse(
	ctx context.Context,
	stuID, courseID string,
	item *kitexModel.CustomCourse,
) (string, error) {
	id, err := strconv.ParseInt(courseID, 10, 64)
	if err != nil {
		return "", errno.CustomCourseNotFoundError
	}
	existing, err := s.db.Course.GetCustomCourseByID(ctx, stuID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errno.CustomCourseNotFoundError
		}
		return "", err
	}
	if _, err := s.db.Course.UpdateCustomCourse(ctx, stuID, id, map[string]any{
		"name":        item.Name,
		"teacher":     getStringValue(item.Teacher),
		"location":    item.Location,
		"start_class": int(item.StartClass),
		"end_class":   int(item.EndClass),
		"start_week":  int(item.StartWeek),
		"end_week":    int(item.EndWeek),
		"weekday":     int(item.Weekday),
		"is_single":   item.Single,
		"is_double":   item.Double,
		"color":       getStringValueWithDefault(item.Color, "#FF5733"),
		"remark":      getStringValue(item.Remark),
	}); err != nil {
		return "", err
	}
	s.refreshCustomCourseCache(stuID, existing.Term)
	return courseID, nil
}

func (s *CourseService) DeleteCustomCourse(ctx context.Context, stuID string, req *course.DeleteCustomCourseRequest) error {
	id, err := strconv.ParseInt(req.CourseId, 10, 64)
	if err != nil {
		return errno.CustomCourseNotFoundError
	}
	existing, err := s.db.Course.GetCustomCourseByID(ctx, stuID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.CustomCourseNotFoundError
		}
		return errno.InternalServiceError.WithError(err)
	}
	rows, err := s.db.Course.DeleteCustomCourse(ctx, stuID, id)
	if err != nil {
		return errno.InternalServiceError.WithError(err)
	}
	if rows == 0 {
		return errno.CustomCourseNotFoundError
	}
	s.refreshCustomCourseCache(stuID, existing.Term)
	return nil
}

func getStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func getStringValueWithDefault(value *string, defaultValue string) string {
	if value == nil || *value == "" {
		return defaultValue
	}
	return *value
}
