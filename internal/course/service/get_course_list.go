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

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func (s *CourseService) GetCourseList(req *course.CourseListRequest, loginData *kitexModel.LoginData) ([]*kitexModel.Course, error) {
	var err error
	stuId := context.ExtractIDFromLoginData(loginData)
	termKey := fmt.Sprintf("terms:%s", stuId)
	courseKey := fmt.Sprintf("course:%s:%s", stuId, req.Term)
	terms := new(jwch.Term)
	// 学期缓存存在
	isRefresh := false
	if req.IsRefresh != nil {
		isRefresh = *req.IsRefresh
	}
	// 不刷新且cache存在
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
			return s.removeDuplicateCourses(pack.BuildCourse(courses)), nil
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
		s.taskQueue.Add(courseKey, taskqueue.QueueTask{Execute: func() error {
			return cache.SetSliceCache(s.cache, s.ctx, courseKey, courses,
				constants.CourseTermsKeyExpire, "Course.SetCourseCache")
		}})
		s.taskQueue.Add(termKey, taskqueue.QueueTask{Execute: func() error {
			return cache.SetValueSliceCache(s.cache, s.ctx, termKey, terms.Terms, constants.CourseTermsKeyExpire, "Course.SetTermsCache")
		}})
	}

	// async put course list to db
	s.taskQueue.Add(fmt.Sprintf("putCourse:%s", stuId), taskqueue.QueueTask{Execute: func() error {
		return s.putCourseToDatabase(stuId, req.Term, pack.BuildCourse(courses))
	}})

	return s.removeDuplicateCourses(pack.BuildCourse(courses)), nil
}

func (s *CourseService) putCourseToDatabase(stuId string, term string, courses []*kitexModel.Course) error {
	old, err := s.db.Course.GetUserTermCourseSha256ByStuIdAndTerm(s.ctx, stuId, term)
	if err != nil {
		return err
	}

	json, err := utils.JSONEncode(courses)
	if err != nil {
		return err
	}

	newSha256 := utils.SHA256(json)

	if old == nil {
		dbId, err := s.sf.NextVal()
		if err != nil {
			return err
		}

		_, err = s.db.Course.CreateUserTermCourse(s.ctx, &model.UserCourse{
			Id:                dbId,
			StuId:             stuId,
			Term:              term,
			TermCourses:       json,
			TermCoursesSha256: newSha256,
		})
		if err != nil {
			return err
		}
	} else if old.TermCoursesSha256 != newSha256 {
		_, err = s.db.Course.UpdateUserTermCourse(s.ctx, &model.UserCourse{
			Id:                old.Id,
			TermCourses:       json,
			TermCoursesSha256: newSha256,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *CourseService) GetCourseListYjsy(req *course.CourseListRequest, loginData *kitexModel.LoginData) ([]*kitexModel.Course, error) {
	var err error

	stuId := context.ExtractIDFromLoginData(loginData)
	termKey := fmt.Sprintf("terms:%s", stuId)
	courseKey := fmt.Sprintf("course:%s:%s", stuId, req.Term)
	terms := new(yjsy.Term)
	// 学期缓存存在
	isRefresh := false
	if req.IsRefresh != nil {
		isRefresh = *req.IsRefresh
	}
	if !isRefresh && s.cache.IsKeyExist(s.ctx, termKey) {
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
			return pack.BuildCourseYjsy(courses), nil
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
		s.taskQueue.Add(courseKey, taskqueue.QueueTask{Execute: func() error {
			return cache.SetSliceCache(s.cache, s.ctx, courseKey, courses,
				constants.CourseTermsKeyExpire, "Course.SetCourseCache")
		}})
		s.taskQueue.Add(termKey, taskqueue.QueueTask{Execute: func() error {
			return cache.SetValueSliceCache(s.cache, s.ctx, termKey, terms.Terms, constants.CourseTermsKeyExpire, "Course.SetTermsCache")
		}})
	}

	// 异步将课程列表存入数据库
	s.taskQueue.Add(fmt.Sprintf("putCourse:%s", stuId), taskqueue.QueueTask{Execute: func() error {
		return s.putCourseToDatabase(stuId, req.Term, pack.BuildCourseYjsy(courses))
	}})

	return pack.BuildCourseYjsy(courses), nil
}

// removeDuplicateCourses 移除重复课程，只保留第一个出现的。
func (s *CourseService) removeDuplicateCourses(courses []*kitexModel.Course) []*kitexModel.Course {
	seen := make(map[string]struct{})
	var result []*kitexModel.Course

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

func (s *CourseService) getSemesterCourses(stuID string, term string) (course []*kitexModel.Course, err error) {
	courseKey := fmt.Sprintf("course:%s:%s", stuID, term)
	if s.cache.IsKeyExist(s.ctx, courseKey) {
		courses, err := s.cache.Course.GetCoursesCache(s.ctx, courseKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetSemesterCourses: Get courses fail: %w", err)
		}
		return s.removeDuplicateCourses(pack.BuildCourse(courses)), nil
	}
	// 从数据中获取课程表
	var courses *model.UserCourse
	courses, err = s.db.Course.GetUserTermCourseByStuIdAndTerm(s.ctx, stuID, term)
	if err != nil {
		return nil, fmt.Errorf("service.GetSemesterCourses: Get courses fail: %w", err)
	}
	if courses == nil {
		return nil, errno.NewErrNo(errno.InternalServiceErrorCode, "service.GetSemesterCourses: there is no course in database, please login app and retry")
	}
	// 将数据库中的课程表进行解析转化
	list := make([]*kitexModel.Course, 0)
	if err = sonic.Unmarshal([]byte(courses.TermCourses), &list); err != nil {
		return nil, fmt.Errorf("service.GetSemesterCourses: Unmarshal fail: %w", err)
	}
	// 写入 cache
	s.taskQueue.Add(courseKey, taskqueue.QueueTask{Execute: func() error {
		return cache.SetSliceCache(s.cache, s.ctx, courseKey, list,
			constants.CourseTermsKeyExpire, "Course.SetCourseCache")
	}})
	return list, nil
}
