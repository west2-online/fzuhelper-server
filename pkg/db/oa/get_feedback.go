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

package oa

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (c *DBOA) GetFeedbackById(ctx context.Context, fbId int64) (bool, *model.Feedback, error) {
	fbModel := new(model.Feedback)
	if err := c.client.WithContext(ctx).Table(constants.FeedbackTableName).Where("report_id = ?", fbId).First(fbModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		logger.Errorf("dal.GetFeedbackById error:%v", err)
		return false, nil, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.GetFeedbackById error:%v", err)
	}
	return true, fbModel, nil
}

func (c *DBOA) ListFeedback(ctx context.Context, req model.FeedbackListReq) (items []model.FeedbackListItem, nextPageToken int64, err error) {
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		logger.Warnf("service.GetFeedbackList: limit out of range, fix to 20 (limit=%d)", limit)
		limit = 20
	}

	tx := c.client.WithContext(ctx).
		Table(constants.FeedbackTableName).
		Select(
			"report_id",
			"name",
			"network_env",
			"problem_desc",
			"app_version",
			"created_at",
			"updated_at",
		)

	if req.StuId != "" {
		tx = tx.Where("stu_id = ?", req.StuId)
	}
	if req.Name != "" {
		tx = tx.Where("name = ?", req.Name)
	}
	if req.NetworkEnv != "" {
		tx = tx.Where("network_env = ?", req.NetworkEnv)
	}
	if req.IsOnCampus != nil {
		tx = tx.Where("is_on_campus = ?", *req.IsOnCampus)
	}
	if req.OsName != "" {
		tx = tx.Where("os_name = ?", req.OsName)
	}
	if req.AppVersion != "" {
		tx = tx.Where("app_version = ?", req.AppVersion)
	}
	if req.BeginTime != nil {
		tx = tx.Where("created_at >= ?", *req.BeginTime)
	}
	if req.EndTime != nil {
		tx = tx.Where("created_at < ?", *req.EndTime)
	}
	if req.ProblemDesc != "" {
		like := fmt.Sprintf("%%%s%%", req.ProblemDesc)
		tx = tx.Where("problem_desc LIKE ?", like)
	}

	orderDesc := true
	if req.OrderDesc != nil {
		orderDesc = *req.OrderDesc
	}
	order := "report_id DESC"
	cmp := "<"
	if !orderDesc {
		order = "report_id ASC"
		cmp = ">"
	}
	if req.PageToken != 0 {
		tx = tx.Where(fmt.Sprintf("report_id %s ?", cmp), req.PageToken)
	}

	// 判断是否还有下一页
	var rows []model.FeedbackListItem
	if err = tx.Order(order).Limit(limit + 1).Find(&rows).Error; err != nil {
		logger.Errorf("dal.ListFeedback error: %v", err)
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "dal.ListFeedback error:%v", err)
	}

	// 处理 nextPageToken
	if len(rows) > limit {
		nextPageToken = rows[limit-1].ReportId
		rows = rows[:limit]
	} else {
		nextPageToken = 0
	}

	return rows, nextPageToken, nil
}
