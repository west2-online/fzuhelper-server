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
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func TestGetScores(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/jwch/academic/scores",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/jwch/academic/scores",
			mockRPCError:   errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/jwch/academic/scores", GetScores)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetScoresRPC).To(func(ctx context.Context, req *academic.GetScoresRequest) ([]*model.Score, error) {
				if tc.mockRPCError != nil {
					return nil, tc.mockRPCError
				}
				return []*model.Score{}, nil
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetGPA(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/jwch/academic/gpa",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/jwch/academic/gpa",
			mockRPCError:   errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/jwch/academic/gpa", GetGPA)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetGPARPC).To(func(ctx context.Context, req *academic.GetGPARequest) (*model.GPABean, error) {
				if tc.mockRPCError != nil {
					return nil, tc.mockRPCError
				}
				return &model.GPABean{}, nil
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetCredit(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/jwch/academic/credit",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/jwch/academic/credit",
			mockRPCError:   errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/jwch/academic/credit", GetCredit)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetCreditRPC).To(func(ctx context.Context, req *academic.GetCreditRequest) ([]*model.Credit, error) {
				if tc.mockRPCError != nil {
					return nil, tc.mockRPCError
				}
				return []*model.Credit{}, nil
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetUnifiedExam(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/jwch/academic/unifiedExam",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/jwch/academic/unifiedExam",
			mockRPCError:   errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/jwch/academic/unifiedExam", GetUnifiedExam)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetUnifiedExamRPC).To(func(ctx context.Context, req *academic.GetUnifiedExamRequest) ([]*model.UnifiedExam, error) {
				if tc.mockRPCError != nil {
					return nil, tc.mockRPCError
				}
				return []*model.UnifiedExam{}, nil
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

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

func TestGetCreditV2(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockRPCError   error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v2/jwch/academic/credit",
			mockRPCError:   nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v2/jwch/academic/credit",
			mockRPCError:   errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/jwch/academic/credit", GetCreditV2)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetCreditV2RPC).To(func(ctx context.Context, req *academic.GetCreditV2Request) (*model.CreditResponse, error) {
				if tc.mockRPCError != nil {
					return nil, tc.mockRPCError
				}
				return &model.CreditResponse{}, nil
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}
