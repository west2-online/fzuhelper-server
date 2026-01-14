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

	"github.com/west2-online/fzuhelper-server/internal/oa/pack"
	"github.com/west2-online/fzuhelper-server/internal/oa/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/oa"
	"github.com/west2-online/fzuhelper-server/pkg/base"
)

// OAServiceImpl implements the last service interface defined in the IDL.
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
	reportID, err := l.CreateFeedback(pack.BuildServiceCreateFeedbackReq(req))
	resp.Base = base.BuildBaseResp(err)
	resp.ReportId = reportID
	return resp, nil
}

func (s *OAServiceImpl) GetFeedbackById(ctx context.Context, req *oa.GetFeedbackByIDRequest) (resp *oa.FeedbackDetailResponse, err error) {
	resp = new(oa.FeedbackDetailResponse)
	l := service.NewOAService(ctx, "", nil, s.ClientSet)
	fb, err := l.GetFeedbackById(req.ReportId)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Data = pack.BuildOAFeedbackDetailResponse(fb)
	return resp, nil
}

func (s *OAServiceImpl) GetFeedbackList(ctx context.Context, req *oa.GetListFeedbackRequest) (resp *oa.GetListFeedbackResponse, err error) {
	resp = new(oa.GetListFeedbackResponse)
	l := service.NewOAService(ctx, "", nil, s.ClientSet)
	items, next, err := l.GetFeedbackList(pack.BuildServiceFeedbackListReq(req))
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Data = pack.BuildOAListItems(items)
	resp.PageToken = &next
	return resp, nil
}
