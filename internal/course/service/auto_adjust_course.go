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
	"fmt"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (s *CourseService) GetAutoAdjustCourseList(term string) ([]*model.AutoAdjustCourse, error) {
	key := s.cache.Course.AutoAdjustCourseKey(term)

	if s.cache.IsKeyExist(s.ctx, key) {
		list, err := s.cache.Course.GetAutoAdjustCourseListCache(s.ctx, key)
		if err != nil {
			return nil, fmt.Errorf("service.GetAutoAdjustCourseList: Get cache failed: %w", err)
		}
		return list, nil
	}

	list, err := s.db.Course.GetAutoAdjustCourseListByTerm(s.ctx, term)
	if err != nil {
		return nil, fmt.Errorf("service.GetAutoAdjustCourseList: Get from db failed: %w", err)
	}

	go func() {
		if err := s.cache.Course.SetAutoAdjustCourseListCache(s.ctx, key, list); err != nil {
			logger.Errorf("service.GetAutoAdjustCourseList: Set cache failed: %v", err)
		}
	}()

	return list, nil
}
