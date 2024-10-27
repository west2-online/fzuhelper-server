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

package main

import (
	"context"

	"github.com/west2-online/fzuhelper-server/cmd/common/pack"
	"github.com/west2-online/fzuhelper-server/cmd/common/service"
	common "github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// CommonServiceImpl implements the last service interface defined in the IDL.
type CommonServiceImpl struct{}

// GetTermsList implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTermsList(ctx context.Context, req *common.TermListRequest) (resp *common.TermListResponse, err error) {
	resp = common.NewTermListResponse()

	res, err := service.NewTermService(ctx).GetTermList()
	if err != nil {
		logger.Errorf("Common.GetTermsList: GetTermList failed, err: %v", err)
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.CurrentTerm = res.CurrentTerm
	resp.Terms = pack.BuildTermList(res.Terms)

	return resp, nil
}

// GetTerm implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTerm(ctx context.Context, req *common.TermRequest) (resp *common.TermResponse, err error) {
	resp = common.NewTermResponse()

	res, err := service.NewTermService(ctx).GetTerm(req)
	if err != nil {
		logger.Errorf("Common.GetTerm: GetTerm failed, err: %v", err)
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.TermId = res.TermId
	resp.SchoolYear = res.SchoolYear
	resp.Term = res.Term
	resp.Events = pack.BuildTermEvents(res.Events)

	return resp, nil
}
