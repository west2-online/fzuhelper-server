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
	"sync"

	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	login_model "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *CourseService) GetCourseList(req *course.CourseListRequest) ([]*jwch.Course, error) {
	loginData, err := context.GetLoginData(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetCourseList: Get login data fail: %w", err)
	}

	// 对相同 ID 的请求加锁，防止疯狂刷新导致并发问题
	lock := s.getLock(context.ExtractIDFromLoginData(loginData))
	if lock != nil {
		if !lock.TryLock() {
			return nil, errno.Errorf(errno.BizLimitCode, "GetCourseList: refresh too fast")
		}
		defer lock.Unlock()
	}

	cnt, err := s.cache.Course.SetAndGetRefreshCount(s.ctx, context.ExtractIDFromLoginData(loginData))
	if err != nil {
		return nil, fmt.Errorf("service.GetCourseList: Set/Get refresh count fail: %w", err)
	}
	if cnt != constants.MaxRefreshCount {
		cachedCourses, needRequest, err := s.getCourseListFromCache(req, loginData)
		if err != nil {
			return nil, err
		}
		// 如果缓存可用，则直接返回
		if cachedCourses != nil && !needRequest {
			return cachedCourses, nil
		}
	} else {
		// 跳过从缓存拿取并删除计数
		task := model.NewDeleteRefreshCountTask(s.ctx, s.cache, context.ExtractIDFromLoginData(loginData))
		s.taskQueue.Add(task)
	}

	// jwch
	courses, terms, err := s.requestCoursesFromJwch(req, loginData)
	if err != nil {
		return nil, err
	}

	// 异步存入缓存
	if slices.Contains(pack.GetTop2Terms(terms).Terms, req.Term) {
		setCoursesTask := model.NewSetCoursesCacheTask(s.ctx, s.cache, context.ExtractIDFromLoginData(loginData), req.Term, courses)
		s.taskQueue.Add(setCoursesTask)

		setTermsTask := model.NewSetTermsCacheTask(s.ctx, s.cache, context.ExtractIDFromLoginData(loginData), terms.Terms)
		s.taskQueue.Add(setTermsTask)
	}

	// 异步存入数据库
	s.putCourseListToDBAsync(req, loginData, courses)

	return courses, nil
}

// 从缓存获取课程列表，并判断是否需要请求教务处
func (s *CourseService) getCourseListFromCache(req *course.CourseListRequest, loginData *login_model.LoginData) ([]*jwch.Course, bool, error) {
	// 判断学期缓存是否存在
	if s.cache.IsKeyExist(s.ctx, context.ExtractIDFromLoginData(loginData)) {
		termsList, err := s.cache.Course.GetTermsCache(s.ctx, context.ExtractIDFromLoginData(loginData))
		if err != nil {
			return nil, false, fmt.Errorf("service.GetCourseList: Get term fail: %w", err)
		}

		terms := &jwch.Term{Terms: termsList}
		key := fmt.Sprintf("course:%s:%s", context.ExtractIDFromLoginData(loginData), req.Term)

		// 如果请求的学期在最新的两个学期中，且课程缓存存在，则直接返回
		if slices.Contains(pack.GetTop2Terms(terms).Terms, req.Term) && s.cache.IsKeyExist(s.ctx, key) {
			courses, err := s.cache.Course.GetCoursesCache(s.ctx, key)
			if err != nil {
				return nil, false, fmt.Errorf("service.GetCourseList: Get courses fail: %w", err)
			}
			return *courses, false, nil
		}
	}
	// 缓存不可用或者需要拉取新数据
	return nil, true, nil
}

// requestCoursesFromJwch 请求教务处，获取学期与课程信息
func (s *CourseService) requestCoursesFromJwch(req *course.CourseListRequest, loginData *login_model.LoginData) ([]*jwch.Course, *jwch.Term, error) {
	stu := jwch.NewStudent().WithLoginData(loginData.GetId(), utils.ParseCookies(loginData.GetCookies()))

	// 获取所有学期
	terms, err := stu.GetTerms()
	if err = base.HandleJwchError(err); err != nil {
		return nil, nil, fmt.Errorf("service.GetCourseList: Get terms failed: %w", err)
	}

	// 验证学期是否存在
	if !slices.Contains(terms.Terms, req.Term) {
		return nil, nil, errors.New("service.GetCourseList: Invalid term")
	}

	// 获取指定学期课程
	courses, err := stu.GetSemesterCourses(req.Term, terms.ViewState, terms.EventValidation)
	if err = base.HandleJwchError(err); err != nil {
		return nil, nil, fmt.Errorf("service.GetCourseList: Get semester courses failed: %w", err)
	}

	return courses, terms, nil
}

// putCourseListToDBAsync 将课程列表异步存入数据库
func (s *CourseService) putCourseListToDBAsync(req *course.CourseListRequest, loginData *login_model.LoginData, courses []*jwch.Course) {
	task := model.NewPutCourseListToDatabaseTask(s.ctx, s.db, context.ExtractIDFromLoginData(loginData), s.sf, req.Term, courses)
	s.taskQueue.Add(task)
}

func (s *CourseService) getLock(key string) *sync.Mutex {
	actual, _ := s.courseLockMap.LoadOrStore(key, &sync.Mutex{})
	l, e := actual.(*sync.Mutex)
	if !e {
		return nil
	}
	return l
}
