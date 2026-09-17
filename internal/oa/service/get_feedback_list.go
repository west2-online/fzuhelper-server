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

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func (s *OAService) GetFeedbackList(req *FeedbackListReq) ([]model.FeedbackListItem, int64, error) {
	if req == nil {
		logger.WithCtx(s.ctx).Errorf("service.GetFeedbackList error: request is nil")
		return nil, 0, fmt.Errorf("service.GetFeedbackList: request is nil")
	}

	stuID := strings.TrimSpace(req.StuId)
	if stuID == "" {
		return nil, 0, fmt.Errorf("service.GetFeedbackList: missing stu_id")
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		logger.WithCtx(s.ctx).Warnf("service.GetFeedbackList: limit out of range, fix to 20 (limit=%d)", limit)
		limit = 20
	}

	orderDesc := true
	if req.OrderDesc != nil {
		orderDesc = *req.OrderDesc
	}

	items, next, err := s.db.OA.ListFeedback(s.ctx, model.FeedbackListReq{
		StuId:     stuID,
		Limit:     limit,
		PageToken: req.PageToken,
		OrderDesc: &orderDesc,
	})
	if err != nil {
		logger.WithCtx(s.ctx).Errorf("service.GetFeedbackList dal error: %v", err)
		return nil, 0, fmt.Errorf("service.GetFeedbackList: list feedback failed: %w", err)
	}

	if items == nil {
		items = []model.FeedbackListItem{}
	}
	return items, next, nil
}
