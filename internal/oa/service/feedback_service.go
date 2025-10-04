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
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *OAService) CreateFeedback(req *CreateFeedbackReq) error {

	// 检验 not null 部分（选择性检验）
	if req.ReportId == 0 || req.StuId == "" || req.Name == "" || req.College == "" ||
		req.ContactPhone == "" || req.ContactQQ == "" || req.ContactEmail == "" ||
		req.OsName == "" || req.OsVersion == "" || req.ProblemDesc == "" || req.AppVersion == "" {
		return errno.Errorf(errno.InternalServiceErrorCode, "missing required fields")
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
		req.NetworkEnv = string(model.NetworkUnknown)
	}

	// 组装 model
	fb := &model.Feedback{
		ReportId:       req.ReportId,
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
		hlog.Errorf("service.CreateFeedback dal error: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "service.CreateFeedback error: %v", err)
	}

	return nil
}

func (s *OAService) GetFeedback(id int64) (*model.Feedback, error) {
	if id <= 0 {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "invalid id: %d", id)
	}
	ok, fb, err := s.db.OA.GetFeedbackById(s.ctx, id)
	if ok == false || err != nil {
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
		fb.NetworkEnv = model.NetworkUnknown
	}

	return fb, nil
}
