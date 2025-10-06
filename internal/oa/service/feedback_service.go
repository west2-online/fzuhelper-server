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
	"strings"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *OAService) CreateFeedback(req *CreateFeedbackReq) (int64, error) {
	// 检验 not null 部分（选择性检验）
	if req.StuId == "" || req.Name == "" || req.College == "" ||
		req.ContactPhone == "" || req.ContactQQ == "" || req.ContactEmail == "" ||
		req.OsName == "" || req.OsVersion == "" || req.ProblemDesc == "" || req.AppVersion == "" {
		return 0, errno.Errorf(errno.InternalServiceErrorCode, "missing required fields")
	}

	// 将空的或不合法的 json 列替换为 [] 或 {}
	req.Screenshots = utils.EnsureJSONArray(req.Screenshots)
	req.VersionHistory = utils.EnsureJSONArray(req.VersionHistory)
	req.NetworkTraces = utils.EnsureJSON(req.NetworkTraces)
	req.Events = utils.EnsureJSONArray(req.Events)
	req.UserSettings = utils.EnsureJSONObject(req.UserSettings)

	// 将不合法 NetworkEnv 归为 Unknown
	switch req.NetworkEnv {
	case string(model.Network2G), string(model.Network3G), string(model.Network4G),
		string(model.Network5G), string(model.NetworkWifi), string(model.NetworkUnknown):
	default:
		logger.Warnf("invalid NetworkEnv=%q, fallback=%q (stu_id=%s)",
			req.NetworkEnv, model.NetworkUnknown, req.StuId)
		req.NetworkEnv = string(model.NetworkUnknown)
	}

	// 生成 reportID
	reportID, err := s.sf.NextVal()
	if err != nil {
		return 0, errno.Errorf(errno.InternalServiceErrorCode, "generate report_id failed: %v", err)
	}

	fb := &model.Feedback{
		ReportId:       reportID,
		StuId:          req.StuId,
		Name:           req.Name,
		College:        req.College,
		ContactPhone:   req.ContactPhone,
		ContactQQ:      req.ContactQQ,
		ContactEmail:   req.ContactEmail,
		NetworkEnv:     model.NetworkEnv(req.NetworkEnv),
		IsOnCampus:     req.IsOnCampus,
		OsName:         req.OsName,
		OsVersion:      req.OsVersion,
		Manufacturer:   req.Manufacturer,
		DeviceModel:    req.DeviceModel,
		ProblemDesc:    req.ProblemDesc,
		Screenshots:    req.Screenshots,
		AppVersion:     req.AppVersion,
		VersionHistory: req.VersionHistory,
		NetworkTraces:  req.NetworkTraces,
		Events:         req.Events,
		UserSettings:   req.UserSettings,
	}

	if err := s.db.OA.CreateFeedback(s.ctx, fb); err != nil {
		logger.Errorf("service.CreateFeedback dal error: %v", err)
		return 0, errno.Errorf(errno.InternalDatabaseErrorCode, "service.CreateFeedback error: %v", err)
	}

	return reportID, nil
}

func (s *OAService) GetFeedbackById(id int64) (*model.Feedback, error) {
	if id <= 0 {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "invalid id: %d", id)
	}
	ok, fb, err := s.db.OA.GetFeedbackById(s.ctx, id)
	if !ok || err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "service.GetFeedback error: %v", err)
	}

	fb.Screenshots = utils.EnsureJSONArray(fb.Screenshots)
	fb.VersionHistory = utils.EnsureJSONArray(fb.VersionHistory)
	fb.NetworkTraces = utils.EnsureJSON(fb.NetworkTraces)
	fb.Events = utils.EnsureJSONArray(fb.Events)
	fb.UserSettings = utils.EnsureJSONObject(fb.UserSettings)

	switch fb.NetworkEnv {
	case model.Network2G, model.Network3G, model.Network4G,
		model.Network5G, model.NetworkWifi, model.NetworkUnknown:
	default:
		logger.Warnf("feedback has invalid stored NetworkEnv, coercing to %q (report_id=%d, original=%q, stu_id=%s)",
			model.NetworkUnknown, fb.ReportId, fb.NetworkEnv, fb.StuId)
		fb.NetworkEnv = model.NetworkUnknown
	}

	return fb, nil
}

func (s *OAService) GetFeedbackList(req *FeedbackListReq) ([]model.FeedbackListItem, int64, error) {
	if req == nil {
		logger.Errorf("service.GetFeedbackList error: request is nil")
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "service.GetFeedbackList error: request is nil")
	}

	// 调整limit
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		logger.Warnf("service.GetFeedbackList: limit out of range, fix to 20 (limit=%d)", limit)
		limit = 20
	}

	// 去空白，防止空格
	if req.ProblemDesc != "" {
		req.ProblemDesc = strings.TrimSpace(req.ProblemDesc)
	}
	if req.StuId != "" {
		req.StuId = strings.TrimSpace(req.StuId)
	}
	if req.Name != "" {
		req.Name = strings.TrimSpace(req.Name)
	}
	if req.OsName != "" {
		req.OsName = strings.TrimSpace(req.OsName)
	}
	if req.AppVersion != "" {
		req.AppVersion = strings.TrimSpace(req.AppVersion)
	}
	if req.NetworkEnv != "" {
		req.NetworkEnv = strings.TrimSpace(req.NetworkEnv)
	}

	// 时间范围校验
	if req.BeginTime != nil && req.EndTime != nil && !req.EndTime.After(*req.BeginTime) {
		logger.Errorf("service.GetFeedbackList: invalid time range, begin=%v end=%v (swapping ignored)",
			*req.BeginTime, *req.EndTime)
		return nil, 0, errno.Errorf(errno.InternalServiceErrorCode, "invalid time range")
	}

	// 排序方式（默认为升序）
	orderDesc := true
	if req.OrderDesc != nil {
		orderDesc = *req.OrderDesc
	}

	req.Limit = limit
	req.OrderDesc = &orderDesc

	if req.NetworkEnv != "" {
		switch req.NetworkEnv {
		case string(model.Network2G), string(model.Network3G), string(model.Network4G),
			string(model.Network5G), string(model.NetworkWifi), string(model.NetworkUnknown):
		default:
			logger.Warnf("service.GetFeedbackList: invalid NetworkEnv=%q, coerce to %q",
				req.NetworkEnv, string(model.NetworkUnknown))
			req.NetworkEnv = string(model.NetworkUnknown)
		}
	}

	listReq := model.FeedbackListReq{
		StuId:       req.StuId,
		Name:        req.Name,
		NetworkEnv:  model.NetworkEnv(req.NetworkEnv),
		IsOnCampus:  req.IsOnCampus,
		OsName:      req.OsName,
		ProblemDesc: req.ProblemDesc,
		AppVersion:  req.AppVersion,
		Limit:       req.Limit,
		PageToken:   req.PageToken,
		OrderDesc:   req.OrderDesc,
		BeginTime:   req.BeginTime,
		EndTime:     req.EndTime,
	}
	items, next, err := s.db.OA.ListFeedback(s.ctx, listReq)
	if err != nil {
		logger.Errorf("service.GetFeedbackList dal error: %v", err)
		return nil, 0, errno.Errorf(errno.InternalDatabaseErrorCode, "list feedback error: %v", err)
	}

	if items == nil {
		items = []model.FeedbackListItem{}
	}

	logger.Infof("service: %s", items[0].ProblemDesc)
	return items, next, nil
}
