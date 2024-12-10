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

package common

import (
	"context"
	"fmt"

	"github.com/west2-online/fzuhelper-server/internal/common/pack"
	"github.com/west2-online/fzuhelper-server/internal/common/service"
	common "github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/pkg/base"
)

// CommonServiceImpl implements the last service interface defined in the IDL.
type CommonServiceImpl struct{}

// GetTermsList implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTermsList(ctx context.Context, req *common.TermListRequest) (resp *common.TermListResponse, err error) {
	resp = common.NewTermListResponse()

	res, err := service.NewTermService(ctx).GetTermList()
	if err != nil {
		resp.Base = base.BuildBaseResp(fmt.Errorf("Common.GetTermsList: get terms list failed: %w", err))
		return resp, nil
	}

	resp.Base = base.BuildBaseResp(nil)
	resp.TermLists = pack.BuildTermsList(res)
	return
}

// GetTerm implements the CommonServiceImpl interface.
func (s *CommonServiceImpl) GetTerm(ctx context.Context, req *common.TermRequest) (resp *common.TermResponse, err error) {
	resp = common.NewTermResponse()

	res, err := service.NewTermService(ctx).GetTerm(req)
	if err != nil {
		resp.Base = base.BuildBaseResp(fmt.Errorf("Common.GetTerm: get term failed: %w", err))
		return resp, nil
	}

	resp.Base = base.BuildBaseResp(nil)
	resp.TermInfo = pack.BuildTermInfo(res)
	return
}
