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
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/academic"
)

func TestGetPlan(t *testing.T) {
	type TestCase struct {
		Name           string
		Headers        map[string]string
		ExpectedError  bool
		ExpectedResult string
		Url            string
	}

	testCases := []TestCase{
		{
			Name: "获取计划成功",
			Headers: map[string]string{
				"id":      "202511012137102301000",
				"Cookies": "ASP.NET_SessionId=db123abcdefgh5ijklzjv2et",
			},
			ExpectedError:  false,
			ExpectedResult: "Plan Content",
			Url:            "/api/v1/jwch/academic/plan",
		},
		{
			Name: "RPC调用失败",
			Headers: map[string]string{
				"id":      "202511012137102301000",
				"Cookies": "ASP.NET_SessionId=db123abcdefgh5ijklzjv2et",
			},
			ExpectedError:  true,
			ExpectedResult: `{"code":"50001","message":"GetCultivatePlanRPC: RPC called failed:`,
			Url:            "/api/v1/jwch/academic/plan",
		},
		{
			Name: "RPC返回空HTML",
			Headers: map[string]string{
				"id":      "202511012137102301000",
				"Cookies": "ASP.NET_SessionId=db123abcdefgh5ijklzjv2et",
			},
			ExpectedError:  true,
			ExpectedResult: `{"code":"50001","message":`,
			Url:            "/api/v1/jwch/academic/plan",
		},
	}

	// 初始化路由引擎并注册GetPlan路由
	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/jwch/academic/plan", GetPlan)
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// 模拟RPC调用
			mockey.Mock(rpc.GetCultivatePlanRPC).To(func(ctx context.Context, req *academic.GetPlanRequest) (string, error) {
				if tc.ExpectedError {
					// 根据测试用例的不同，可以自定义返回不同的错误信息
					return "", errors.New("GetCultivatePlanRPC: RPC called failed: 错误的文件路径")
				}
				// 成功返回HTML内容
				html := "Plan Content"
				return html, nil
			}).Build()
			resp := ut.PerformRequest(router, consts.MethodGet, tc.Url, nil)

			// 断言响应
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Contains(t, string(resp.Result().Body()), tc.ExpectedResult)
		})
	}
}
