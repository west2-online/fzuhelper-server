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
	"fmt"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/oa"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitOARPC() {
	c, err := client.InitOARPC()
	if err != nil {
		logger.Fatalf("api.rpc.oa InitOARPC failed, err is %v", err)
	}
	logger.Infof("InitOARPC: etcd=%s service=%s", config.Etcd.Addr, constants.OAServiceName)
	oaClient = *c
}

func CreateFeedbackRPC(ctx context.Context, req *oa.CreateFeedbackRequest) (int64, error) {
	resp, err := oaClient.CreateFeedback(ctx, req)
	if err != nil {
		logger.Errorf("CreateFeedbackRPC: RPC called failed: %v", err.Error())
		return 0, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return 0, errno.BizError.WithMessage(fmt.Sprintf("创建反馈表单失败：%s", resp.Base.Msg))
	}
	return resp.ReportId, nil
}

func GetFeedbackByIdRPC(ctx context.Context, req *oa.GetFeedbackByIDRequest) (*model.Feedback, error) {
	resp, err := oaClient.GetFeedbackById(ctx, req)
	if err != nil {
		logger.Errorf("GetFeedbackByIdRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.BizError.WithMessage(fmt.Sprintf("查询反馈表单详情失败：%s", resp.Base.Msg))
	}
	return resp.Data, nil
}

func GetFeedbackListRPC(ctx context.Context, req *oa.GetListFeedbackRequest) ([]*model.FeedbackListItem, *int64, error) {
	resp, err := oaClient.GetFeedbackList(ctx, req)
	if err != nil {
		logger.Errorf("GetFeedbackListRPC: RPC called failed: %v", err.Error())
		return nil, nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, nil, errno.BizError.WithMessage(fmt.Sprintf("查询反馈表单列表失败：%s", resp.Base.Msg))
	}
	return resp.Data, resp.PageToken, nil
}
