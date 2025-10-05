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

package api

import (
	"bytes"
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	oa "github.com/west2-online/fzuhelper-server/kitex_gen/oa"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestCreateFeedback(t *testing.T) {
	type testCase struct {
		name           string
		body           string
		mockRPCError   error
		expectingError bool
		expectingMsg   string
		url            string
	}

	okBody := `{
		"reportId": 1,
		"stuId": "102301000",
		"name": "张三",
		"college": "计算机与大数据学院",
		"contactPhone": "13800000000",
		"contactQQ": "10001",
		"contactEmail": "a@b.com",
		"networkEnv": "wifi",
		"isOnCampus": true,
		"osName": "Android",
		"osVersion": "14",
		"manufacturer": "Xiaomi",
		"deviceModel": "Mi 14",
		"problemDesc": "登录白屏",
		"screenshots": "[]",
		"appVersion": "1.2.3",
		"versionHistory": "[]",
		"networkTraces": "[]",
		"events": "[]",
		"userSettings": "{}"
	}`

	testCases := []testCase{
		{
			name:           "success",
			body:           okBody,
			mockRPCError:   nil,
			expectingError: false,
			expectingMsg:   `{"code":"10000","message":`,
			url:            "/api/v1/feedback/create",
		},
		{
			name:           "invalid json",
			body:           `{"reportId": 1,`, // 非法 JSON
			mockRPCError:   nil,
			expectingError: true,
			expectingMsg:   `{"code":"20001","message":`,
			url:            "/api/v1/feedback/create",
		},
		{
			name:           "rpc error",
			body:           okBody,
			mockRPCError:   errno.InternalServiceError,
			expectingError: true,
			expectingMsg:   `{"code":"50001","message":`,
			url:            "/api/v1/feedback/create",
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/feedback/create", CreateFeedback)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.CreateFeedbackRPC).To(func(ctx context.Context, req *oa.CreateFeedbackRequest) error {
				return tc.mockRPCError
			}).Build()

			result := ut.PerformRequest(router, "POST", tc.url,
				&ut.Body{Body: bytes.NewBufferString(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"})
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectingMsg)
		})
	}
}
