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
	"errors"
	"fmt"
	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	login_model "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"slices"
)

func (s *CourseService) GetCourseList(req *course.CourseListRequest) ([]*jwch.Course, error) {
	var loginData *login_model.LoginData
	var err error

	loginData, err = context.GetLoginData(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetCourseList: Get login data fail: %w", err)
	}

	terms := new(jwch.Term)
	// 学期缓存存在
	if s.cache.IsKeyExist(s.ctx, context.ExtractIDFromLoginData(loginData)) {
		termsList, err := s.cache.Course.GetTermsCache(s.ctx, context.ExtractIDFromLoginData(loginData))
		if err != nil {
			return nil, fmt.Errorf("service.GetCourseList: Get term fail: %w", err)
		}
		terms.Terms = termsList

		key := fmt.Sprintf("course:%s:%s", context.ExtractIDFromLoginData(loginData), req.Term)
		// 只有最新的两个学期的课程才会被放入缓存
		if slices.Contains(pack.GetTop2Terms(terms).Terms, req.Term) &&
			s.cache.IsKeyExist(s.ctx, key) {
			courses, err := s.cache.Course.GetCoursesCache(s.ctx, key)
			if err != nil {
				return nil, fmt.Errorf("service.GetCourseList: Get courses fail: %w", err)
			}
			return *courses, nil
		}
	}

	stu := jwch.NewStudent().WithLoginData(loginData.GetId(), utils.ParseCookies(loginData.GetCookies()))

	terms, err = stu.GetTerms()
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetCourseList: Get terms failed: %w", err)
	}

	// validate term
	if !slices.Contains(terms.Terms, req.Term) {
		return nil, errors.New("service.GetCourseList: Invalid term")
	}

	courses, err := stu.GetSemesterCourses(req.Term, terms.ViewState, terms.EventValidation)
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetCourseList: Get semester courses failed: %w", err)
	}

	if slices.Contains(pack.GetTop2Terms(terms).Terms, req.Term) {
		// async put course list to cache
		setCoursesTask := model.NewSetCoursesCacheTask(s.ctx, s.cache, context.ExtractIDFromLoginData(loginData), req.Term, courses)
		s.taskQueue.Add(setCoursesTask)

		setTermsTask := model.NewSetTermsCacheTask(s.ctx, s.cache, context.ExtractIDFromLoginData(loginData), terms.Terms)
		s.taskQueue.Add(setTermsTask)
	}

	// async put course list to db
	putCourseListTask := model.NewPutCourseListToDatabaseTask(s.ctx, s.db, context.ExtractIDFromLoginData(loginData), s.sf, req.Term, courses)
	s.taskQueue.Add(putCourseListTask)

	return courses, nil
}
