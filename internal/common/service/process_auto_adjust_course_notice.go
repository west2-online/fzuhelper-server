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

// 处理教务通知中的课程调整信息
func (s *CommonService) ProcessAutoAdjustCourseNotice(info *jwch.NoticeInfo) error {
	logger.Infof("ProcessAutoAdjustCourseNotice: processing notice, title=%s url=%s", info.Title, info.URL)

	// 仅处理标题包含"课程调整"的通知，其余通知直接跳过
	if !strings.Contains(info.Title, "课程调整") {
		return nil
	}

	// 根据通知的 WbTreeId 和 WbNewsId 获取通知详情（含正文 HTML）
	detail, err := jwch.NewStudent().GetNoticeDetail(&jwch.NoticeDetailReq{
		WbTreeId: info.WbTreeId,
		WbNewsId: info.WbNewsId,
	})
	if err != nil {
		return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to get notice detail: %w", err)
	}

	// 调用 LLM 从通知标题和正文中提取结构化的课程调整条目
	result, err := ai.AutoAdjustCourse(s.ctx, ai.AutoAdjustCourseInput{
		Title:   info.Title,
		Content: detail.Content,
	})
	if err != nil {
		return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to auto adjust course: %w", err)
	}

	logger.Infof("ProcessAutoAdjustCourseNotice: AI extracted %+v", result.Items)

	// 获取学期列表，用于后续将日期映射到具体学期
	calendar, err := s.GetTermList()
	if err != nil {
		return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to get term list: %w", err)
	}

	// termsToRefresh 收集本次写入了新调课记录的学期，处理完所有条目后统一刷新缓存，
	// 使用 map 以学期标识去重，避免对同一学期重复刷新。
	termsToRefresh := make(map[string]jwch.CalTerm)
	for _, item := range result.Items {
		// 解析调课来源日期，同时提取所在自然年（用于数据库字段 Year）
		fromDate, err := utils.TimeParse(item.FromDate)
		if err != nil {
			logger.Errorf("ProcessAutoAdjustCourseNotice: invalid from date %s: %v", item.FromDate, err)
			continue
		}
		year := strconv.Itoa(fromDate.Year())

		// toDate 为空表示该课程被取消（无补课日期），否则校验目标日期合法性
		toDate := &item.ToDate
		if item.ToDate == "" {
			// 课程取消的情况：只有 FromDate（被取消的上课日），没有 ToDate（补课日）
			toDate = nil
		} else {
			_, err = utils.TimeParse(item.ToDate)
			if err != nil {
				logger.Errorf("ProcessAutoAdjustCourseNotice: invalid to date %s: %v", item.ToDate, err)
				continue
			}
		}

		// 根据来源日期查找所属学期；若日期不在任何已知学期范围内则跳过
		term, found := utils.FindTermByDate(calendar.Terms, fromDate)
		if !found {
			logger.Warnf("ProcessAutoAdjustCourseNotice: no term found for date %s, skipping", item.FromDate)
			continue
		}

		// 加入待刷新 Map，后续统一刷新该学期缓存
		termsToRefresh[term.Term] = term

		// 计算来源日期在该学期中对应的周次和星期几
		fromWeek, fromWeekday, err := utils.GetWeekdayByDate(term.StartDate, item.FromDate)
		if err != nil {
			logger.Errorf("ProcessAutoAdjustCourseNotice: failed to get week info for %s: %v", item.FromDate, err)
			continue
		}

		// 目标周次和星期默认为 nil；仅当存在补课日期时才计算并赋值
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

		// 构造调课记录并写入数据库；Enabled 默认为 false，等待人工审核后再启用
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

	// 遍历所有涉及的学期，重新从数据库读取完整的调课列表并刷新缓存，
	// 确保后续查询能立即感知到本次新增的调课记录。
	for _, term := range termsToRefresh {
		key := s.cache.Course.AutoAdjustCourseKey(term.Term)

		// 获取当前学期所有的课程调整信息（含本次新增）
		adjustCourses, err := s.db.Course.GetAutoAdjustCourseListByTerm(s.ctx, term.Term)
		if err != nil {
			return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to get auto adjust course list: %w", err)
		}

		// 将最新的调课列表写入缓存
		err = s.cache.Course.SetAutoAdjustCourseListCache(s.ctx, key, adjustCourses)
		if err != nil {
			return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to cache auto adjust course list: %w", err)
		}
	}

	return nil
}
