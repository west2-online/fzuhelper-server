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
	"slices"
	"sort"
	"strings"

	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	loginmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func (s *CourseService) GetCourseList(req *course.CourseListRequest, loginData *loginmodel.LoginData) ([]*jwch.Course, error) {
	var err error

	termKey := fmt.Sprintf("terms:%s", context.ExtractIDFromLoginData(loginData))
	courseKey := strings.Join([]string{context.ExtractIDFromLoginData(loginData), req.Term}, ":")
	terms := new(jwch.Term)
	// 学期缓存存在
	isRefresh := false
	if req.IsRefresh != nil {
		isRefresh = *req.IsRefresh
	}
	if !isRefresh && s.cache.IsKeyExist(s.ctx, termKey) {
		termsList, err := s.cache.Course.GetTermsCache(s.ctx, termKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetCourseList: Get term fail: %w", err)
		}
		terms.Terms = termsList
		// 只有最新的两个学期的课程才会被放入缓存
		if slices.Contains(pack.GetTop2Terms(terms).Terms, req.Term) &&
			s.cache.IsKeyExist(s.ctx, courseKey) {
			courses, err := s.cache.Course.GetCoursesCache(s.ctx, courseKey)
			if err != nil {
				return nil, fmt.Errorf("service.GetCourseList: Get courses fail: %w", err)
			}
			return s.removeDuplicateCourses(*courses), nil
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
		setTermsTask := model.NewSetTermsCacheTask(s.ctx, s.cache, termKey, terms.Terms)
		s.taskQueue.Add(setTermsTask)
	}

	// async put course list to db
	putCourseListTask := model.NewPutCourseListToDatabaseTask(s.ctx, s.db, context.ExtractIDFromLoginData(loginData), s.sf, req.Term, courses)
	s.taskQueue.Add(putCourseListTask)

	return s.removeDuplicateCourses(courses), nil
}

func (s *CourseService) GetCourseListYjsy(req *course.CourseListRequest, loginData *loginmodel.LoginData) ([]*yjsy.Course, error) {
	var err error

	termKey := fmt.Sprintf("terms:%s", context.ExtractIDFromLoginData(loginData))
	courseKey := strings.Join([]string{context.ExtractIDFromLoginData(loginData), req.Term}, ":")
	terms := new(yjsy.Term)
	// 学期缓存存在
	if s.cache.IsKeyExist(s.ctx, termKey) {
		termsList, err := s.cache.Course.GetTermsCache(s.ctx, termKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetCourseListYjsy: Get terms fail: %w", err)
		}
		terms.Terms = termsList

		// 检查是否有该学期的课程缓存
		if slices.Contains(pack.GetTop2TermsYjsy(terms).Terms, req.Term) && s.cache.IsKeyExist(s.ctx, courseKey) {
			courses, err := s.cache.Course.GetCoursesCacheYjsy(s.ctx, courseKey)
			if err != nil {
				return nil, fmt.Errorf("service.GetCourseListYjsy: Get courses fail: %w", err)
			}
			return *courses, nil
		}
	}

	// 获取学期信息
	stu := yjsy.NewStudent().WithLoginData(utils.ParseCookies(loginData.Cookies))
	terms, err = stu.GetTerms()
	if err = base.HandleYjsyError(err); err != nil {
		return nil, fmt.Errorf("service.GetCourseListYjsy: Get terms failed: %w", err)
	}

	// 验证学期是否有效
	if !slices.Contains(terms.Terms, req.Term) {
		return nil, errors.New("service.GetCourseListYjsy: Invalid term")
	}

	// 获取该学期的课程
	courses, err := stu.GetSemesterCourses(req.Term)
	if err = base.HandleYjsyError(err); err != nil {
		return nil, fmt.Errorf("service.GetCourseListYjsy: Get semester courses failed: %w", err)
	}

	// 如果是前两个学期，异步缓存课程列表
	if slices.Contains(pack.GetTop2TermsYjsy(terms).Terms, req.Term) {
		setCoursesTask := model.NewSetCoursesCacheTaskYjsy(s.ctx, s.cache, context.ExtractIDFromLoginData(loginData), req.Term, courses)
		s.taskQueue.Add(setCoursesTask)

		setTermsTask := model.NewSetTermsCacheTask(s.ctx, s.cache, termKey, terms.Terms)
		s.taskQueue.Add(setTermsTask)
	}

	// 异步将课程列表存入数据库
	putCourseListTask := model.NewPutCourseListToDatabaseTaskYjsy(s.ctx, s.db, context.ExtractIDFromLoginData(loginData), s.sf, req.Term, courses)
	s.taskQueue.Add(putCourseListTask)

	return courses, nil
}

// removeDuplicateCourses 移除重复课程，只保留第一个出现的。
func (s *CourseService) removeDuplicateCourses(courses []*jwch.Course) []*jwch.Course {
	seen := make(map[string]struct{})
	var result []*jwch.Course

	for _, c := range courses {
		srIDs := make([]string, 0, len(c.ScheduleRules))
		for _, rule := range c.ScheduleRules {
			part := fmt.Sprintf("%d-%d-%d-%d",
				rule.StartClass, rule.EndClass,
				rule.StartWeek, rule.EndWeek)
			srIDs = append(srIDs, part)
		}
		sort.Strings(srIDs)

		// 把“课程名 + 教师 + 排课信息”拼成一个全局唯一的 key
		identifier := fmt.Sprintf("%s-%s-%s", c.Name, c.Teacher, strings.Join(srIDs, "|"))

		// 如果 map 里还没出现过这个标识，那就是新课程
		if _, exists := seen[identifier]; !exists {
			seen[identifier] = struct{}{}
			result = append(result, c)
		}
	}

	return result
}
