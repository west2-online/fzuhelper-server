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
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (s *OaService) CreateFeedback(in *CreateFeedbackReq) error {

	// 检验 not null 部分（选择性检验）
	if in.ReportId == 0 || in.StuId == "" || in.Name == "" || in.College == "" ||
		in.ContactPhone == "" || in.ContactQQ == "" || in.ContactEmail == "" ||
		in.OsName == "" || in.OsVersion == "" || in.ProblemDesc == "" || in.AppVersion == "" {
		return errno.Errorf(errno.InternalServiceErrorCode, "missing required fields")
	}

	// 将空的或不合法的 json 列替换为 [] 或 {}
	in.Screenshots = ensureJSONArray(in.Screenshots)
	in.VersionHistory = ensureJSONArray(in.VersionHistory)
	in.NetworkTraces = ensureJSON(in.NetworkTraces)
	in.Events = ensureJSONArray(in.Events)
	in.UserSettings = ensureJSONObject(in.UserSettings)

	// 将不合法 NetworkEnv 归为 Unknown
	switch in.NetworkEnv {
	case string(model.Network2G), string(model.Network3G), string(model.Network4G),
		string(model.Network5G), string(model.NetworkWifi), string(model.NetworkUnknown):
	default:
		in.NetworkEnv = string(model.NetworkUnknown)
	}

	// 组装 model
	fb := &model.Feedback{
		ReportId:       in.ReportId,
		StuId:          in.StuId,
		Name:           in.Name,
		College:        in.College,
		ContactPhone:   in.ContactPhone,
		ContactQQ:      in.ContactQQ,
		ContactEmail:   in.ContactEmail,
		NetworkEnv:     model.NetworkEnv(in.NetworkEnv),
		IsOnCampus:     in.IsOnCampus,
		OsName:         in.OsName,
		OsVersion:      in.OsVersion,
		Manufacturer:   in.Manufacturer,
		DeviceModel:    in.DeviceModel,
		ProblemDesc:    in.ProblemDesc,
		Screenshots:    in.Screenshots,
		AppVersion:     in.AppVersion,
		VersionHistory: in.VersionHistory,
		NetworkTraces:  in.NetworkTraces,
		Events:         in.Events,
		UserSettings:   in.UserSettings,
	}

	if err := s.db.Oa.CreateFeedback(s.ctx, fb); err != nil {
		hlog.Errorf("service.CreateFeedback dal error: %v", err)
		return errno.Errorf(errno.InternalDatabaseErrorCode, "service.CreateFeedback error: %v", err)
	}

	return nil
}

func (s *OaService) GetFeedback(id int64) (*model.Feedback, error) {
	if id <= 0 {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "invalid id: %d", id)
	}
	ok, fb, err := s.db.Oa.GetFeedbackById(s.ctx, id)
	if ok == false || err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "service.GetFeedback error: %v", err)
	}

	fb.Screenshots = ensureJSONArray(fb.Screenshots)
	fb.VersionHistory = ensureJSONArray(fb.VersionHistory)
	fb.NetworkTraces = ensureJSON(fb.NetworkTraces)
	fb.Events = ensureJSONArray(fb.Events)
	fb.UserSettings = ensureJSONObject(fb.UserSettings)

	switch fb.NetworkEnv {
	case model.Network2G, model.Network3G, model.Network4G,
		model.Network5G, model.NetworkWifi, model.NetworkUnknown:
	default:
		fb.NetworkEnv = model.NetworkUnknown
	}

	return fb, nil
}

// 转换函数
func ensureJSONArray(s string) string {
	if s == "" {
		return "[]"
	}
	if json.Valid([]byte(s)) {
		return s
	}
	return "[]"
}
func ensureJSONObject(s string) string {
	if s == "" {
		return "{}"
	}
	if json.Valid([]byte(s)) {
		return s
	}
	return "{}"
}
func ensureJSON(s string) string {
	if s == "" {
		return "[]"
	}
	if json.Valid([]byte(s)) {
		return s
	}
	return "[]"
}
