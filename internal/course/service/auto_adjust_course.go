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
	"strconv"
	"time"

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	rpcmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *CourseService) getAutoAdjustCourseList(term string) ([]*model.AutoAdjustCourse, error) {
	key := s.cache.Course.AutoAdjustCourseKey(term)

	if s.cache.IsKeyExist(s.ctx, key) {
		list, err := s.cache.Course.GetAutoAdjustCourseListCache(s.ctx, key)
		if err != nil {
			return nil, fmt.Errorf("service.getAutoAdjustCourseList: Get cache failed: %w", err)
		}
		return list, nil
	}

	list, err := s.db.Course.GetAutoAdjustCourseListByTerm(s.ctx, term)
	if err != nil {
		return nil, fmt.Errorf("service.getAutoAdjustCourseList: Get from db failed: %w", err)
	}

	go func() {
		if err := s.cache.Course.SetAutoAdjustCourseListCache(s.ctx, key, list); err != nil {
			logger.Errorf("service.getAutoAdjustCourseList: Set cache failed: %v", err)
		}
	}()

	return list, nil
}

func (s *CourseService) GetAutoAdjustCourseList(secret, term string) ([]*model.AutoAdjustCourse, error) {
	if err := s.db.AdminSecret.ValidateSecret(s.ctx, constants.AdjustCourseModuleName, secret); err != nil {
		return nil, fmt.Errorf("service.GetAutoAdjustCourseList: Validate secret failed: %w", err)
	}
	return s.getAutoAdjustCourseList(term)
}

func (s *CourseService) UpdateAutoAdjustCourse(req *course.UpdateAdjustCourseRequest) error {
	if err := s.db.AdminSecret.ValidateSecret(s.ctx, constants.AdjustCourseModuleName, req.Secret); err != nil {
		return fmt.Errorf("service.UpdateAutoAdjustCourse: Validate secret failed: %w", err)
	}

	adjustCourse := &model.AutoAdjustCourse{
		Id: req.Id,
	}

	if req.Enabled != nil {
		adjustCourse.Enabled = req.GetEnabled()
	}

	if req.FromDate != nil || req.ToDate != nil {
		resp, err := s.commonClient.GetTermsList(s.ctx, &common.TermListRequest{})
		if err != nil {
			return fmt.Errorf("service.UpdateAutoAdjustCourse: Get terms list failed: %w", err)
		}
		if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
			return fmt.Errorf("service.UpdateAutoAdjustCourse: term list resp error: %w", err)
		}

		if req.FromDate != nil {
			fromDateStr := req.GetFromDate()
			fromDate, err := utils.TimeParse(fromDateStr)
			if err != nil {
				return fmt.Errorf("service.UpdateAutoAdjustCourse: invalid from_date %s: %w", fromDateStr, err)
			}

			term, found := findTermByDate(resp.TermLists.Terms, fromDate)
			if !found {
				return fmt.Errorf("service.UpdateAutoAdjustCourse: no term found for date %s", fromDateStr)
			}

			fromWeek, fromWeekday, err := utils.GetWeekdayByDate(term.GetStartDate(), fromDateStr)
			if err != nil {
				return fmt.Errorf("service.UpdateAutoAdjustCourse: failed to get week info for %s: %w", fromDateStr, err)
			}

			adjustCourse.FromDate = fromDateStr
			adjustCourse.FromWeek = int64(fromWeek)
			adjustCourse.FromWeekday = int64(fromWeekday)
			adjustCourse.Term = term.GetTerm()
			adjustCourse.Year = strconv.Itoa(fromDate.Year())
		}

		if req.ToDate != nil {
			toDateStr := req.GetToDate()
			if toDateStr == "" {
				// 空字符串表示课程取消
				adjustCourse.ToDate = nil
			} else {
				toDate, err := utils.TimeParse(toDateStr)
				if err != nil {
					return fmt.Errorf("service.UpdateAutoAdjustCourse: invalid to_date %s: %w", toDateStr, err)
				}

				term, found := findTermByDate(resp.TermLists.Terms, toDate)
				if !found {
					return fmt.Errorf("service.UpdateAutoAdjustCourse: no term found for to_date %s", toDateStr)
				}

				toWeek, toWeekday, err := utils.GetWeekdayByDate(term.GetStartDate(), toDateStr)
				if err != nil {
					return fmt.Errorf("service.UpdateAutoAdjustCourse: failed to get week info for to_date %s: %w", toDateStr, err)
				}

				adjustCourse.ToDate = &toDateStr
				adjustCourse.ToWeek = int64(toWeek)
				adjustCourse.ToWeekday = int64(toWeekday)
			}
		}
	}

	// 获取原记录用于刷新缓存
	original, err := s.db.Course.GetAutoAdjustCourseByID(s.ctx, req.Id)
	if err != nil {
		return fmt.Errorf("service.UpdateAutoAdjustCourse: Get original record failed: %w", err)
	}
	oldTerm := original.Term

	if err := s.db.Course.UpdateAutoAdjustCourse(s.ctx, adjustCourse); err != nil {
		return fmt.Errorf("service.UpdateAutoAdjustCourse: Update failed: %w", err)
	}

	// 刷新缓存，如果改了学期，那旧的也要刷新
	termsToRefresh := []string{oldTerm}
	if adjustCourse.Term != "" && adjustCourse.Term != oldTerm {
		termsToRefresh = append(termsToRefresh, adjustCourse.Term)
	}

	go func() {
		for _, term := range termsToRefresh {
			key := s.cache.Course.AutoAdjustCourseKey(term)
			list, err := s.db.Course.GetAutoAdjustCourseListByTerm(s.ctx, term)
			if err != nil {
				logger.Errorf("service.UpdateAutoAdjustCourse: Refresh cache get list for term %s failed: %v", term, err)
				continue
			}
			if err = s.cache.Course.SetAutoAdjustCourseListCache(s.ctx, key, list); err != nil {
				logger.Errorf("service.UpdateAutoAdjustCourse: Set cache for term %s failed: %v", term, err)
			}
		}
	}()

	return nil
}

func findTermByDate(terms []*rpcmodel.Term, date time.Time) (*rpcmodel.Term, bool) {
	for _, term := range terms {
		if term.StartDate == nil || term.EndDate == nil {
			continue
		}
		startDate, err := utils.TimeParse(*term.StartDate)
		if err != nil {
			continue
		}
		endDate, err := utils.TimeParse(*term.EndDate)
		if err != nil {
			continue
		}
		if !date.Before(startDate) && !date.After(endDate) {
			return term, true
		}
	}
	return nil, false
}
