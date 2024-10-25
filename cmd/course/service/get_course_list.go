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

	"github.com/west2-online/fzuhelper-server/cmd/course/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *CourseService) GetCourseList(req *course.CourseListRequest) ([]*jwch.Course, error) {
	stu := jwch.NewStudent().WithLoginData(req.LoginData.Id, utils.ParseCookies(req.LoginData.Cookies))

	terms, err := stu.GetTerms()
	if err != nil {
		return nil, fmt.Errorf("service.GetCourseList: Get terms failed: %w", err)
	}

	// validate term
	if !slices.Contains(terms.Terms, req.Term) {
		return nil, errors.New("service.GetCourseList: Invalid term")
	}

	courses, err := stu.GetSemesterCourses(req.Term, terms.ViewState, terms.EventValidation)
	if err != nil {
		return nil, fmt.Errorf("service.GetCourseList: Get semester courses failed: %w", err)
	}

	// async put course list to db
	go func() {
		if err := s.putCourseListToDatabase(req.LoginData.Id, req.Term, courses); err != nil {
			// do nothing because error has already been logged
		}
	}()

	return courses, nil
}

func (s *CourseService) putCourseListToDatabase(id string, term string, courses []*jwch.Course) error {
	stuId, err := utils.ParseJwchStuId(id)
	if err != nil {
		logger.Errorf("service.putCourseList: ParseJwchStuId failed: %v", err)
		return err
	}

	old, err := db.GetUserTermCourseSha256ByStuIdAndTerm(s.ctx, stuId, term)
	if err != nil {
		logger.Errorf("service.putCourseList: GetUserTermCourseSha256ByStuIdAndTerm failed: %v", err)
		return err
	}

	json, err := utils.JSONEncode(courses)
	if err != nil {
		logger.Errorf("service.putCourseList: JSONEncode failed: %v", err)
		return err
	}

	newSha256 := utils.SHA256(json)

	if old == nil {
		dbId, err := db.SF.NextVal()
		if err != nil {
			logger.Errorf("service.putCourseList: SF.NextVal failed: %v", err)
			return err
		}

		_, err = db.CreateUserTermCourse(s.ctx, &db.UserCourse{
			Id:                dbId,
			StuId:             stuId,
			Term:              term,
			TermCourses:       json,
			TermCoursesSha256: newSha256,
		})
		if err != nil {
			logger.Errorf("service.putCourseList: CreateUserTermCourse failed: %v", err)
			return err
		}
	} else if old.TermCoursesSha256 != newSha256 {
		_, err = db.UpdateUserTermCourse(s.ctx, &db.UserCourse{
			Id:                old.Id,
			TermCourses:       json,
			TermCoursesSha256: newSha256,
		})
		if err != nil {
			logger.Errorf("service.putCourseList: UpdateUserTermCourse failed: %v", err)
			return err
		}
	}

	return nil
}
