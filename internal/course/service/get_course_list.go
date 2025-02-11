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
	"strings"

	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	login_model "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *CourseService) GetCourseList(req *course.CourseListRequest) ([]*jwch.Course, error) {
	var loginData *login_model.LoginData
	var err error

	loginData, err = context.GetLoginData(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetCourseList: Get login data fail: %w", err)
	}

	terms := new(jwch.Term)
	// 缓存存在
	if s.cache.IsKeyExist(s.ctx, loginData.GetId()) {
		termsList, err := s.cache.Course.GetTermsCache(s.ctx, loginData.GetId())
		if err != nil {
			return nil, fmt.Errorf("service.GetCourseList: Get term fail: %w", err)
		}
		terms.Terms = *termsList

		// 只有最新的两个学期的才会被放入缓存
		key := strings.Join([]string{loginData.GetId(), req.Term}, ":")
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
		go func() {
			err = s.cache.Course.SetCoursesCache(s.ctx, strings.Join([]string{loginData.GetId(), req.Term}, ":"), &courses)
			if err != nil {
				logger.Errorf("service.GetCourseList: SetCoursesCache failed: %v", err)
			}
		}()
		go func() {
			err = s.cache.Course.SetTermsCache(s.ctx, loginData.GetId(), &terms.Terms)
			if err != nil {
				logger.Errorf("service.GetCourseList: SetTermsCache failed: %v", err)
			}
		}()
	}
	// async put course list to db
	go func() {
		if err := s.putCourseListToDatabase(loginData.GetId(), req.Term, courses); err != nil {
			logger.Errorf("service.GetCourseList: putCourseListToDatabase failed: %v", err)
		}
	}()

	return courses, nil
}

func (s *CourseService) putCourseListToDatabase(id string, term string, courses []*jwch.Course) error {
	stuId, err := utils.ParseJwchStuId(id)
	if err != nil {
		return fmt.Errorf("service.putCourseListToDatabase: ParseJwchStuId failed: %w", err)
	}

	old, err := s.db.Course.GetUserTermCourseSha256ByStuIdAndTerm(s.ctx, stuId, term)
	if err != nil {
		return fmt.Errorf("service.putCourseListToDatabase: GetUserTermCourseSha256ByStuIdAndTerm failed: %w", err)
	}

	json, err := utils.JSONEncode(courses)
	if err != nil {
		return fmt.Errorf("service.putCourseListToDatabase: JSONEncode failed: %w", err)
	}

	newSha256 := utils.SHA256(json)

	if old == nil {
		dbId, err := s.sf.NextVal()
		if err != nil {
			return fmt.Errorf("service.putCourseListToDatabase: SF.NextVal failed: %w", err)
		}

		_, err = s.db.Course.CreateUserTermCourse(s.ctx, &model.UserCourse{
			Id:                dbId,
			StuId:             stuId,
			Term:              term,
			TermCourses:       json,
			TermCoursesSha256: newSha256,
		})
		if err != nil {
			return fmt.Errorf("service.putCourseListToDatabase: CreateUserTermCourse failed: %w", err)
		}
	} else if old.TermCoursesSha256 != newSha256 {
		_, err = s.db.Course.UpdateUserTermCourse(s.ctx, &model.UserCourse{
			Id:                old.Id,
			TermCourses:       json,
			TermCoursesSha256: newSha256,
		})
		if err != nil {
			return fmt.Errorf("service.putCourseListToDatabase: UpdateUserTermCourse failed: %w", err)
		}
	}

	return nil
}
