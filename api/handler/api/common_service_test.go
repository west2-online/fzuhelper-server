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
	"testing"

	"github.com/bytedance/mockey"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func TestGetCSS(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *[]byte
		mockErr        error
		expectContains string
	}

	css := []byte("body{color:#000;}")

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v2/url/onekey/FZUHelper.css",
			mockResp:       &css,
			mockErr:        nil,
			expectContains: "body{color:#000;}",
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/onekey/FZUHelper.css",
			mockResp:       nil,
			mockErr:        errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/url/onekey/FZUHelper.css", GetCSS)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetCSSRPC).To(func(ctx context.Context, req *common.GetCSSRequest) (*[]byte, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetHtml(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *[]byte
		mockErr        error
		expectContains string
	}

	html := []byte("<html></html>")

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v2/url/onekey/FZUHelper.html",
			mockResp:       &html,
			mockErr:        nil,
			expectContains: "<html></html>",
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/onekey/FZUHelper.html",
			mockResp:       nil,
			mockErr:        errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/url/onekey/FZUHelper.html", GetHtml)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetHtmlRPC).To(func(ctx context.Context, req *common.GetHtmlRequest) (*[]byte, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetUserAgreement(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *[]byte
		mockErr        error
		expectContains string
	}

	agreement := []byte("agreement")

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v2/url/onekey/UserAgreement.html",
			mockResp:       &agreement,
			mockErr:        nil,
			expectContains: "agreement",
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/onekey/UserAgreement.html",
			mockResp:       nil,
			mockErr:        errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/url/onekey/UserAgreement.html", GetUserAgreement)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetUserAgreementRPC).To(func(ctx context.Context, req *common.GetUserAgreementRequest) (*[]byte, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetNotice(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockNotices    []*model.NoticeInfo
		mockTotal      int64
		mockErr        error
		expectContains string
	}

	notices := []*model.NoticeInfo{{}}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/common/notice?pageNum=1",
			mockNotices:    notices,
			mockTotal:      1,
			mockErr:        nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/common/notice?pageNum=1",
			mockNotices:    nil,
			mockTotal:      0,
			mockErr:        errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/common/notice",
			mockNotices:    nil,
			mockTotal:      0,
			mockErr:        nil,
			expectContains: `{"code":"20001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/common/notice", GetNotice)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetNoticesRPC).To(func(ctx context.Context, req *common.NoticeRequest) ([]*model.NoticeInfo, int64, error) {
				return tc.mockNotices, tc.mockTotal, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetContributorInfo(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *common.GetContributorInfoResponse
		mockErr        error
		expectContains string
	}

	mockResp := &common.GetContributorInfoResponse{}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/vi/common/contributor",
			mockResp:       mockResp,
			mockErr:        nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/vi/common/contributor",
			mockResp:       nil,
			mockErr:        errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/vi/common/contributor", GetContributorInfo)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetContributorRPC).To(func(ctx context.Context, req *common.GetContributorInfoRequest) (*common.GetContributorInfoResponse, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetToolboxConfig(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       []*model.ToolboxConfig
		mockErr        error
		expectContains string
	}

	configs := []*model.ToolboxConfig{{}}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/toolbox/config?version=1",
			mockResp:       configs,
			mockErr:        nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/toolbox/config?version=1",
			mockResp:       nil,
			mockErr:        errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/toolbox/config?version=abc",
			mockResp:       nil,
			mockErr:        nil,
			expectContains: `{"code":"20001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/toolbox/config", GetToolboxConfig)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetToolboxConfigRPC).To(func(ctx context.Context, req *common.GetToolboxConfigRequest) ([]*model.ToolboxConfig, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestPutToolboxConfig(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *common.PutToolboxConfigResponse
		mockErr        error
		expectContains string
	}

	configId := int64(123)
	mockResp := &common.PutToolboxConfigResponse{ConfigId: &configId}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/toolbox/config?secret=abc&tool_id=1",
			mockResp:       mockResp,
			mockErr:        nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/toolbox/config?secret=abc&tool_id=1",
			mockResp:       nil,
			mockErr:        errors.New("rpc error"),
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/toolbox/config",
			mockResp:       nil,
			mockErr:        nil,
			expectContains: `{"code":"20001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.PUT("/api/v1/toolbox/config", PutToolboxConfig)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.PutToolboxConfigRPC).To(func(ctx context.Context, req *common.PutToolboxConfigRequest) (*common.PutToolboxConfigResponse, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPut, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetTerm(t *testing.T) {
	f := func(str string) *string {
		return &str
	}

	type TestCase struct {
		Name              string
		expectedError     bool
		expectedErrorInfo error
		expectedResult    string
		expectedTermInfo  *model.TermInfo
		url               string
	}

	expectedTermInfo := &model.TermInfo{
		TermId:     f("201501"),
		Term:       f("201501"),
		SchoolYear: f("2015"),
		Events: []*model.TermEvent{
			{
				Name:      f("学生注册"),
				StartDate: f("2015-08-29"),
				EndDate:   f("2015-08-30"),
			},
			{
				Name:      f("学生补考"),
				StartDate: f("2015-08-29"),
				EndDate:   f("2015-09-06"),
			},
			{
				Name:      f("正式上课"),
				StartDate: f("2015-08-31"),
				EndDate:   f("2015-08-31"),
			},
			{
				Name:      f("新生报到"),
				StartDate: f("2015-09-07"),
				EndDate:   f("2015-09-07"),
			},
			{
				Name:      f("校运会"),
				StartDate: f("2015-11-12"),
				EndDate:   f("2015-11-14"),
			},
			{
				Name:      f("期末考试"),
				StartDate: f("2016-01-16"),
				EndDate:   f("2016-01-22"),
			},
			{
				Name:      f("寒假"),
				StartDate: f("2016-01-23"),
				EndDate:   f("2016-02-28"),
			},
		},
	}

	data, err := sonic.Marshal(expectedTermInfo)
	assert.Nil(t, err)

	testCases := []TestCase{
		{
			Name:              "GetTermSuccessfully",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult:    `{"code":"10000","message":"Success","data":` + string(data) + `}`,
			expectedTermInfo:  expectedTermInfo,
			url:               "/api/v1/terms/info?term=201501",
		},
		{
			Name:              "GetTermError",
			expectedError:     true,
			expectedErrorInfo: errors.New("etTermRPC: RPC called failed"),
			expectedResult:    `{"code":"50001","message":"etTermRPC: RPC called failed"}`,
			expectedTermInfo:  nil,
			url:               "/api/v1/terms/info?term=201501",
		},
		{
			Name:              "BindAndValidateError",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult: `{"code":"20001","message":"参数错误, 'term' field is a 'required' parameter,` +
				` but the request body does not have this parameter 'term'"}`,
			expectedTermInfo: nil,
			url:              "/api/v1/terms/info",
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/terms/info", GetTerm)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock(rpc.GetTermRPC).To(func(ctx context.Context, req *common.TermRequest) (*model.TermInfo, error) {
				return tc.expectedTermInfo, tc.expectedErrorInfo
			}).Build()

			result := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Equal(t, tc.expectedResult, string(result.Result().Body()))
		})
	}
}

func TestGetTermsList(t *testing.T) {
	f := func(str string) *string {
		return &str
	}

	type TestCase struct {
		Name              string
		expectedError     bool
		expectedErrorInfo error
		expectedResult    string
		expectedTermInfo  *model.TermList
	}

	expectedTermList := &model.TermList{
		CurrentTerm: f("202401"),
		Terms: []*model.Term{
			{
				TermId:     f("2024012024082620250117"),
				SchoolYear: f("2024"),
				Term:       f("202401"),
				StartDate:  f("2024-08-26"),
				EndDate:    f("2025-01-17"),
			},
			{
				TermId:     f("2024022025022420250704"),
				SchoolYear: f("2024"),
				Term:       f("202402"),
				StartDate:  f("2025-02-24"),
				EndDate:    f("2025-07-04"),
			},
		},
	}

	data, err := sonic.Marshal(expectedTermList)
	assert.Nil(t, err)

	testCases := []TestCase{
		{
			Name:              "GetTermsListSuccessfully",
			expectedError:     false,
			expectedErrorInfo: nil,
			expectedResult:    `{"code":"10000","message":"Success","data":` + string(data) + `}`,
			expectedTermInfo:  expectedTermList,
		},
		{
			Name:              "GetTermsListError",
			expectedError:     true,
			expectedErrorInfo: errors.New("etTermRPC: RPC called failed"),
			expectedResult:    `{"code":"50001","message":"etTermRPC: RPC called failed"}`,
			expectedTermInfo:  nil,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/terms/list", GetTermsList)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			mockey.Mock(rpc.GetTermsListRPC).To(func(ctx context.Context, req *common.TermListRequest) (*model.TermList, error) {
				return tc.expectedTermInfo, tc.expectedErrorInfo
			}).Build()
			url := "/api/v1/terms/list"
			result := ut.PerformRequest(router, consts.MethodGet, url, nil)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Equal(t, tc.expectedResult, string(result.Result().Body()))
		})
	}
}
