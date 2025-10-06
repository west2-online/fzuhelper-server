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

package pack

import (
	"github.com/west2-online/fzuhelper-server/internal/oa/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/oa"
	db "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func BuildServiceCreateFeedbackReq(req *oa.CreateFeedbackRequest) *service.CreateFeedbackReq {
	return &service.CreateFeedbackReq{
		StuId:          req.GetStuId(),
		Name:           req.GetName(),
		College:        req.GetCollege(),
		ContactPhone:   req.GetContactPhone(),
		ContactQQ:      req.GetContactQq(),
		ContactEmail:   req.GetContactEmail(),
		NetworkEnv:     req.GetNetworkEnv(),
		IsOnCampus:     req.GetIsOnCampus(),
		OsName:         req.GetOsName(),
		OsVersion:      req.GetOsVersion(),
		Manufacturer:   req.GetManufacturer(),
		DeviceModel:    req.GetDeviceModel(),
		ProblemDesc:    req.GetProblemDesc(),
		Screenshots:    req.GetScreenshots(),
		AppVersion:     req.GetAppVersion(),
		VersionHistory: req.GetVersionHistory(),
		NetworkTraces:  req.GetNetworkTraces(),
		Events:         req.GetEvents(),
		UserSettings:   req.GetUserSettings(),
	}
}

func BuildServiceFeedbackListReq(req *oa.GetListFeedbackRequest) *service.FeedbackListReq {
	return &service.FeedbackListReq{
		StuId:       utils.StrOrEmpty(req.StuId),
		Name:        utils.StrOrEmpty(req.Name),
		NetworkEnv:  utils.StrOrEmpty(req.NetworkEnv),
		IsOnCampus:  req.IsOnCampus,
		OsName:      utils.StrOrEmpty(req.OsName),
		ProblemDesc: utils.StrOrEmpty(req.ProblemDesc),
		AppVersion:  utils.StrOrEmpty(req.AppVersion),
		Limit:       int(utils.I64OrZero(req.Limit)),
		PageToken:   utils.I64OrZero(req.PageToken),
		OrderDesc:   req.OrderDesc,
		BeginTime:   utils.TimePtrFromMillis(req.BeginTimeMs),
		EndTime:     utils.TimePtrFromMillis(req.EndTimeMs),
	}
}

func BuildOAFeedbackDetailResponse(fb *db.Feedback) *model.Feedback {
	return &model.Feedback{
		ReportId:       fb.ReportId,
		StuId:          fb.StuId,
		Name:           fb.Name,
		College:        fb.College,
		ContactPhone:   fb.ContactPhone,
		ContactQq:      fb.ContactQQ,
		ContactEmail:   fb.ContactEmail,
		NetworkEnv:     string(fb.NetworkEnv),
		IsOnCampus:     fb.IsOnCampus,
		OsName:         fb.OsName,
		OsVersion:      fb.OsVersion,
		Manufacturer:   fb.Manufacturer,
		DeviceModel:    fb.DeviceModel,
		ProblemDesc:    fb.ProblemDesc,
		Screenshots:    fb.Screenshots,
		AppVersion:     fb.AppVersion,
		VersionHistory: fb.VersionHistory,
		NetworkTraces:  fb.NetworkTraces,
		Events:         fb.Events,
		UserSettings:   fb.UserSettings,
	}
}

func BuildOAListItems(dbItems []db.FeedbackListItem) []*model.FeedbackListItem {
	if len(dbItems) == 0 {
		return nil
	}
	out := make([]*model.FeedbackListItem, len(dbItems))
	for i := range dbItems {
		it := &dbItems[i]
		out[i] = &model.FeedbackListItem{
			ReportId:    it.ReportId,
			Name:        it.Name,
			NetworkEnv:  string(it.NetworkEnv),
			ProblemDesc: it.ProblemDesc,
			AppVersion:  it.AppVersion,
		}
	}
	return out
}
