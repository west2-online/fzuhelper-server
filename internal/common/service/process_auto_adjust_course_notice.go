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
	"strings"

	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/pkg/ai"
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

	if len(result.Items) == 0 {
		logger.Infof("ProcessAutoAdjustCourseNotice: no adjust course item extracted, skip")
		return nil
	}

	// 只透传原始日期，学期与周次的换算、落库、缓存刷新都由 course 服务完成，
	// 避免 common 跨服务直接操作 course 的数据库与缓存。
	items := make([]*course.CreateAdjustCourseItem, 0, len(result.Items))
	for _, item := range result.Items {
		adjustItem := &course.CreateAdjustCourseItem{FromDate: item.FromDate}
		if item.ToDate != "" {
			toDate := item.ToDate
			adjustItem.ToDate = &toDate
		}
		items = append(items, adjustItem)
	}

	resp, err := s.courseClient.CreateAdjustCourse(s.ctx, &course.CreateAdjustCourseRequest{Items: items})
	if err != nil {
		return fmt.Errorf("ProcessAutoAdjustCourseNotice: failed to create adjust course: %w", err)
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return fmt.Errorf("ProcessAutoAdjustCourseNotice: create adjust course resp error: %w", err)
	}

	logger.Infof("ProcessAutoAdjustCourseNotice: created %d adjust course record(s)", resp.GetCreated())

	return nil
}
