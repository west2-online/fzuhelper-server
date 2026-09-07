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
	"testing"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/kitex_gen/common/commonservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

type jobFairClientStub struct {
	commonservice.Client
	response *common.JobFairResponse
	err      error
}

func (c *jobFairClientStub) GetJobFair(
	context.Context,
	*common.JobFairRequest,
	...callopt.Option,
) (*common.JobFairResponse, error) {
	return c.response, c.err
}

func TestGetJobFairRPCPreservesServiceErrorCode(t *testing.T) {
	originalClient := commonClient
	defer func() { commonClient = originalClient }()

	commonClient = &jobFairClientStub{response: &common.JobFairResponse{
		Base: &model.BaseResp{Code: errno.ParamErrorCode, Msg: "invalid month"},
	}}

	_, err := GetJobFairRPC(context.Background(), &common.JobFairRequest{Month: "2026/09"})
	assert.Equal(t, int64(errno.ParamErrorCode), errno.ConvertErr(err).ErrorCode)
}
