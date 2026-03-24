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
	"strings"

	"github.com/west2-online/fzuhelper-server/pkg/ai"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *CommonService) ProcessAutoAdjustCourseNotice(info *jwch.NoticeInfo) error {
	if !strings.Contains(info.Title, "课程调整") {
		return nil
	}

	detail, err := jwch.NewStudent().GetNoticeDetail(&jwch.NoticeDetailReq{
		WbTreeId: info.WbTreeId,
		WbNewsId: info.WbNewsId,
	})
	if err != nil {
		return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to get notice detail: %w", err)
	}

	content := detail.Content

	result, err := ai.AutoAdjustCourse(ai.AutoAdjustCourseInput{
		Title:   info.Title,
		Content: content,
	})
	if err != nil {
		return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to auto adjust course: %w", err)
	}

	// 获取学期列表
	calendar, err := s.GetTermList()
	if err != nil {
		return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to get term list: %w", err)
	}

	termsToRefresh := make(map[string]jwch.CalTerm)
	for _, item := range result.Items {
		fromDate, err := utils.TimeParse(item.FromDate)
		if err != nil {
			logger.Errorf("ProcessAutoAdjustCourseNotice: invalid from date %s: %v", item.FromDate, err)
			continue
		}
		year := strconv.Itoa(fromDate.Year())

		toDate := &item.ToDate
		if item.ToDate == "" {
			// 课程取消的情况
			toDate = nil
		} else {
			_, err = utils.TimeParse(item.ToDate)
			if err != nil {
				logger.Errorf("ProcessAutoAdjustCourseNotice: invalid to date %s: %v", item.ToDate, err)
				continue
			}
		}

		// 根据日期获取对应学期
		term, found := utils.FindTermByDate(calendar.Terms, fromDate)
		if !found {
			logger.Warnf("ProcessAutoAdjustCourseNotice: no term found for date %s, skipping", item.FromDate)
			continue
		}

		// 加入待刷新Map中
		termsToRefresh[term.Term] = term

		// 获取日期对应学期的周数和星期
		fromWeek, fromWeekday, err := utils.GetWeekdayByDate(term.StartDate, item.FromDate)
		if err != nil {
			logger.Errorf("ProcessAutoAdjustCourseNotice: failed to get week info for %s: %v", item.FromDate, err)
			continue
		}

		// 对应的周数和星期默认为nil，当toDate不为空时才会有值
		var toWeekPtr, toWeekdayPtr *int64
		if toDate != nil {
			toWeek, toWeekday, err := utils.GetWeekdayByDate(term.StartDate, *toDate)
			if err != nil {
				logger.Errorf("ProcessAutoAdjustCourseNotice: failed to get week info for to date %s: %v", *toDate, err)
				continue
			}
			toWeekVal := int64(toWeek)
			toWeekPtr = &toWeekVal
			toWeekdayVal := int64(toWeekday)
			toWeekdayPtr = &toWeekdayVal
		}

		adjustCourse := &model.AutoAdjustCourse{
			Year:        year,
			FromDate:    item.FromDate,
			ToDate:      toDate,
			Term:        term.Term,
			FromWeek:    int64(fromWeek),
			ToWeek:      toWeekPtr,
			FromWeekday: int64(fromWeekday),
			ToWeekday:   toWeekdayPtr,
			Enabled:     false,
		}

		_, err = s.db.Course.CreateAutoAdjustCourse(s.ctx, adjustCourse)
		if err != nil {
			return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to create auto adjust course: %w", err)
		}
	}

	// 刷新缓存
	for _, term := range termsToRefresh {
		key := s.cache.Course.AutoAdjustCourseKey(term.Term)

		// 获取当前学期所有的课程调整信息
		adjustCourses, err := s.db.Course.GetAutoAdjustCourseListByTerm(s.ctx, term.Term)
		if err != nil {
			return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to get auto adjust course list: %w", err)
		}

		// 刷写缓存
		err = s.cache.Course.SetAutoAdjustCourseListCache(s.ctx, key, adjustCourses)
		if err != nil {
			return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to cache auto adjust course list: %w", err)
		}
	}

	return nil
}
