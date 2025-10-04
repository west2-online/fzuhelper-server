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

	"github.com/west2-online/fzuhelper-server/internal/oa/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/oa"
	"github.com/west2-online/fzuhelper-server/pkg/base"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type OAServiceImpl struct {
	ClientSet *base.ClientSet
}

func NewOAService(clientSet *base.ClientSet) *OAServiceImpl {
	return &OAServiceImpl{
		ClientSet: clientSet,
	}
}

func (s *OAServiceImpl) CreateFeedback(ctx context.Context, req *oa.CreateFeedbackRequest) (resp *oa.CreateFeedbackResponse, err error) {
	resp = new(oa.CreateFeedbackResponse)
	l := service.NewOAService(ctx, "", nil, s.ClientSet)
	err = l.CreateFeedback(&service.CreateFeedbackReq{
		ReportId:       req.GetReportId(),
		StuId:          req.GetStuId(),
		Name:           req.GetName(),
		College:        req.GetCollege(),
		ContactPhone:   req.GetContactPhone(),
		ContactQQ:      req.GetContactQQ(),
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
	})
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	return resp, nil
}

func (s *OAServiceImpl) GetFeedback(ctx context.Context, req *oa.GetFeedbackRequest) (resp *oa.GetFeedbackResponse, err error) {
	resp = new(oa.GetFeedbackResponse)
	l := service.NewOAService(ctx, "", nil, s.ClientSet)
	fb, err := l.GetFeedback(req.ReportId)
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Data = &model.Feedback{
		ReportId:       fb.ReportId,
		StuId:          fb.StuId,
		Name:           fb.Name,
		College:        fb.College,
		ContactPhone:   fb.ContactPhone,
		ContactQQ:      fb.ContactQQ,
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
	return
}
