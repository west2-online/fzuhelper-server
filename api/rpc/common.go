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

package rpc

import (
	"context"

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitCommonRPC() {
	c, err := client.InitCommonRPC()
	if err != nil {
		logger.Fatalf("api.rpc.Common InitCommonRPC failed, err  %v", err)
	}
	commonClient = *c
}

func GetTermsListRPC(ctx context.Context, req *common.TermListRequest) (*model.TermList, error) {
	resp, err := commonClient.GetTermsList(ctx, req)
	if err != nil {
		logger.Errorf("GetTermsListRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}

	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, err
	}

	return resp.TermLists, nil
}

func GetTermRPC(ctx context.Context, req *common.TermRequest) (*model.TermInfo, error) {
	resp, err := commonClient.GetTerm(ctx, req)
	if err != nil {
		logger.Errorf("GetTermRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, err
	}

	return resp.TermInfo, nil
}
