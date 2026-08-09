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
	"strings"
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
	"github.com/west2-online/fzuhelper-server/pkg/errno"
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
			expectContains: "body{color:#000;}",
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/onekey/FZUHelper.css",
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"`,
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
			expectContains: "<html></html>",
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/onekey/FZUHelper.html",
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"`,
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
			expectContains: "agreement",
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/onekey/UserAgreement.html",
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"`,
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

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/common/notice?pageNum=1",
			mockNotices:    []*model.NoticeInfo{{}},
			mockTotal:      1,
			expectContains: `{"code":"10000","message":"ok","data":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/common/notice?pageNum=1",
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/common/notice",
			expectContains: `{"code":"20001","message":"参数错误`,
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

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/common/contributor",
			mockResp:       &common.GetContributorInfoResponse{},
			expectContains: `{"code":"10000","message":"ok","data":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/common/contributor",
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/common/contributor", GetContributorInfo)

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

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/toolbox/config?version=1",
			mockResp:       []*model.ToolboxConfig{{}},
			expectContains: `{"code":"10000","message":"ok","data":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/toolbox/config?version=1",
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/toolbox/config?version=abc",
			expectContains: `{"code":"20001","message":"参数错误`,
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

func TestListToolboxConfigs(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       []*model.ToolboxConfigDetail
		mockTotal      int64
		mockErr        error
		expectToolID   *int64
		expectStudent  *string
		expectPlatform *string
		expectVersion  *int64
		expectContains []string
	}

	testCases := []testCase{
		{
			name:      "success",
			url:       "/api/v1/toolbox/configs?secret=abc&page_num=1&page_size=20",
			mockTotal: 1,
			mockResp: []*model.ToolboxConfigDetail{
				{
					ConfigId:  1,
					ToolId:    1,
					StudentId: new("102300217"),
				},
			},
			expectContains: []string{
				`"code":"10000"`,
				`"message":"ok"`,
				`"total":1`,
				`"config_id":1`,
				`"student_id":"102300217"`,
			},
		},
		{
			name:           "success_with_filters",
			url:            "/api/v1/toolbox/configs?secret=abc&tool_id=1&student_id=102300217&platform=android&version=2",
			expectToolID:   new(int64(1)),
			expectStudent:  new("102300217"),
			expectPlatform: new("android"),
			expectVersion:  new(int64(2)),
			mockResp:       []*model.ToolboxConfigDetail{},
			expectContains: []string{`"config":[]`, `"total":0`},
		},
		{
			name:           "rpc error",
			url:            "/api/v1/toolbox/configs?secret=abc&page_num=1&page_size=20",
			mockErr:        errno.InternalServiceError,
			expectContains: []string{`{"code":"50001","message":"内部服务错误"`},
		},
		{
			name:      "empty page",
			url:       "/api/v1/toolbox/configs?secret=abc&page_num=2&page_size=20",
			mockTotal: 0,
			mockResp:  nil,
			expectContains: []string{
				`"config":[]`,
				`"total":0`,
			},
		},
		{
			name:           "missing secret",
			url:            "/api/v1/toolbox/configs?page_num=1&page_size=20",
			expectContains: []string{`{"code":"20001","message":"参数错误`},
		},
		{
			name:           "bind error",
			url:            "/api/v1/toolbox/configs?secret=abc&page_size=abc",
			expectContains: []string{`{"code":"20001","message":"参数错误`},
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/toolbox/configs", ListToolboxConfigs)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.ListToolboxConfigsRPC).To(func(ctx context.Context, req *common.ListToolboxConfigsRequest) ([]*model.ToolboxConfigDetail, int64, error) {
				assert.Equal(t, tc.expectToolID, req.ToolId)
				assert.Equal(t, tc.expectStudent, req.StudentId)
				assert.Equal(t, tc.expectPlatform, req.Platform)
				assert.Equal(t, tc.expectVersion, req.Version)
				return tc.mockResp, tc.mockTotal, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			for _, expectContains := range tc.expectContains {
				assert.Contains(t, string(res.Result().Body()), expectContains)
			}
		})
	}
}

// Full-replacement request bodies for the toolbox config handlers, split across
// lines to stay within the line-length limit.
const (
	toolboxConfigNullTail = `"name":null,"icon":null,"type":null,"message":null,` +
		`"extra":null,"student_id":null,"platform":null,"version":null}`
	toolboxConfigAllNullsBody    = `{"secret":"abc","tool_id":1,"visible":false,` + toolboxConfigNullTail
	toolboxConfigNullVisibleBody = `{"secret":"abc","tool_id":1,"visible":null,` + toolboxConfigNullTail
)

func TestCreateToolboxConfig(t *testing.T) {
	type testCase struct {
		name           string
		body           string
		mockResp       *model.ToolboxConfigDetail
		mockErr        error
		expectContains string
	}

	validBody := `{"secret":"abc","tool_id":1,"visible":false,"name":"","icon":"","type":"","message":"","extra":"","student_id":"","platform":"","version":0}`

	testCases := []testCase{
		{
			name:           "success",
			body:           validBody,
			mockResp:       &model.ToolboxConfigDetail{ConfigId: 123, ToolId: 1, Name: new("")},
			expectContains: `"config":{"config_id":123,"tool_id":1,"visible":false,"name":""`,
		},
		{
			name:           "success with explicit nulls",
			body:           toolboxConfigAllNullsBody,
			mockResp:       &model.ToolboxConfigDetail{ConfigId: 124, ToolId: 1},
			expectContains: `"name":null`,
		},
		{
			name:           "rpc error",
			body:           validBody,
			mockErr:        errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"`,
		},
		{
			name:           "missing full-replacement field",
			body:           `{"secret":"abc","tool_id":1}`,
			expectContains: `{"code":"20005","message":"visible is required"`,
		},
		{
			name:           "non-nullable field cannot be null",
			body:           toolboxConfigNullVisibleBody,
			expectContains: `{"code":"20001","message":"visible cannot be null"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/toolbox/configs", CreateToolboxConfig)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.CreateToolboxConfigRPC).To(func(ctx context.Context, req *common.CreateToolboxConfigRequest) (*model.ToolboxConfigDetail, error) {
				if tc.name == "success with explicit nulls" {
					assert.Nil(t, req.Name)
					assert.Nil(t, req.Icon)
					assert.Nil(t, req.Type)
					assert.Nil(t, req.Message)
					assert.Nil(t, req.Extra)
					assert.Nil(t, req.StudentId)
					assert.Nil(t, req.Platform)
					assert.Nil(t, req.Version)
				}
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(
				router,
				consts.MethodPost,
				"/api/v1/toolbox/configs",
				&ut.Body{Body: strings.NewReader(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"},
			)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetToolboxConfigByID(t *testing.T) {
	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/toolbox/configs/:id", GetToolboxConfigByID)
	type testCase struct {
		name string
		test func(*testing.T)
	}
	testCases := []testCase{
		{name: "path id is forwarded", test: func(t *testing.T) {
			mockey.Mock(rpc.GetToolboxConfigByIDRPC).To(
				func(_ context.Context, req *common.GetToolboxConfigByIDRequest) (*model.ToolboxConfigDetail, error) {
					assert.Equal(t, int64(123), req.ConfigId)
					assert.Equal(t, "abc", req.Secret)
					return &model.ToolboxConfigDetail{ConfigId: 123, ToolId: 1}, nil
				},
			).Build()
			res := ut.PerformRequest(router, consts.MethodGet, "/api/v1/toolbox/configs/123?secret=abc", nil)
			assert.Contains(t, string(res.Result().Body()), `"config_id":123`)
		}},
		{name: "invalid path id", test: func(t *testing.T) {
			res := ut.PerformRequest(router, consts.MethodGet, "/api/v1/toolbox/configs/invalid?secret=abc", nil)
			assert.Contains(t, string(res.Result().Body()), `{"code":"20001","message":"参数错误`)
		}},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		t.Run(tc.name, tc.test)
		mockey.UnPatchAll()
	}
}

func TestUpdateToolboxConfig(t *testing.T) {
	router := route.NewEngine(&config.Options{})
	router.PUT("/api/v1/toolbox/configs/:id", UpdateToolboxConfig)
	type testCase struct {
		name string
		body string
		test func(*testing.T, string)
	}
	testCases := []testCase{
		{name: "full replacement forwards explicit nulls", body: toolboxConfigAllNullsBody, test: func(t *testing.T, body string) {
			mockey.Mock(rpc.UpdateToolboxConfigRPC).To(
				func(_ context.Context, req *common.UpdateToolboxConfigRequest) (*model.ToolboxConfigDetail, error) {
					assert.Equal(t, int64(123), req.ConfigId)
					assert.False(t, req.Visible)
					assert.Nil(t, req.Name)
					assert.Nil(t, req.Icon)
					assert.Nil(t, req.Type)
					assert.Nil(t, req.Message)
					assert.Nil(t, req.Extra)
					assert.Nil(t, req.StudentId)
					assert.Nil(t, req.Platform)
					assert.Nil(t, req.Version)
					return &model.ToolboxConfigDetail{ConfigId: 123, ToolId: 1}, nil
				},
			).Build()
			res := ut.PerformRequest(
				router, consts.MethodPut, "/api/v1/toolbox/configs/123",
				&ut.Body{Body: strings.NewReader(body), Len: len(body)},
				ut.Header{Key: "Content-Type", Value: "application/json"})
			responseBody := string(res.Result().Body())
			assert.Contains(t, responseBody, `"config_id":123`)
			assert.Contains(t, responseBody, `"name":null`)
		}},
		{
			name: "missing property is rejected",
			body: `{"secret":"abc","tool_id":1,"visible":false,"name":null,"icon":null,"type":null,
					"message":null,"extra":null,"student_id":null,"platform":null}`,
			test: func(t *testing.T, body string) {
				res := ut.PerformRequest(
					router, consts.MethodPut, "/api/v1/toolbox/configs/123",
					&ut.Body{Body: strings.NewReader(body), Len: len(body)},
					ut.Header{Key: "Content-Type", Value: "application/json"})
				assert.Contains(t, string(res.Result().Body()), `{"code":"20005","message":"version is required"`)
			},
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { tc.test(t, tc.body) })
		mockey.UnPatchAll()
	}
}

func TestDeleteToolboxConfig(t *testing.T) {
	router := route.NewEngine(&config.Options{})
	router.DELETE("/api/v1/toolbox/configs/:id", DeleteToolboxConfig)
	type testCase struct {
		name string
		test func(*testing.T)
	}
	testCases := []testCase{
		{name: "success", test: func(t *testing.T) {
			mockey.Mock(rpc.DeleteToolboxConfigRPC).To(
				func(_ context.Context, req *common.DeleteToolboxConfigRequest) error {
					assert.Equal(t, int64(123), req.ConfigId)
					return nil
				},
			).Build()
			res := ut.PerformRequest(router, consts.MethodDelete, "/api/v1/toolbox/configs/123?secret=abc", nil)
			assert.Contains(t, string(res.Result().Body()), `{"code":"10000","message":"ok"}`)
		}},
		{name: "rpc error", test: func(t *testing.T) {
			mockey.Mock(rpc.DeleteToolboxConfigRPC).Return(errno.InternalServiceError).Build()
			res := ut.PerformRequest(router, consts.MethodDelete, "/api/v1/toolbox/configs/123?secret=abc", nil)
			assert.Contains(t, string(res.Result().Body()), `{"code":"50001","message":"内部服务错误"}`)
		}},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		t.Run(tc.name, tc.test)
		mockey.UnPatchAll()
	}
}

func TestGetTerm(t *testing.T) {
	f := func(str string) *string {
		return &str
	}

	type TestCase struct {
		Name              string
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
			Name:             "GetTermSuccessfully",
			expectedResult:   `{"code":"10000","message":"Success","data":` + string(data) + `}`,
			expectedTermInfo: expectedTermInfo,
			url:              "/api/v1/terms/info?term=201501",
		},
		{
			Name:              "GetTermError",
			expectedErrorInfo: errno.InternalServiceError,
			expectedResult:    `{"code":"50001","message":"内部服务错误"}`,
			url:               "/api/v1/terms/info?term=201501",
		},
		{
			Name: "BindAndValidateError",
			expectedResult: `{"code":"20001","message":"参数错误, 'term' field is a 'required' parameter,` +
				` but the request body does not have this parameter 'term'"}`,
			url: "/api/v1/terms/info",
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
			Name:             "GetTermsListSuccessfully",
			expectedResult:   `{"code":"10000","message":"Success","data":` + string(data) + `}`,
			expectedTermInfo: expectedTermList,
		},
		{
			Name:              "GetTermsListError",
			expectedErrorInfo: errno.InternalServiceError,
			expectedResult:    `{"code":"50001","message":"内部服务错误"}`,
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

func TestGetSignedLocationApiUrl(t *testing.T) {
	type testCase struct {
		name           string
		mockSignedURL  string
		mockHeaders    map[string]string
		mockErr        error
		body           *ut.Body
		expectContains string
	}

	validBody := `{"location":"119.262647,26.106131"}`

	testCases := []testCase{
		{
			name:           "success",
			mockSignedURL:  "https://restapi.amap.com/v3/place/around?key=xxx&scode=abc",
			mockHeaders:    map[string]string{"User-Agent": "AMAP_Location_SDK_Android"},
			body:           &ut.Body{Body: strings.NewReader(validBody), Len: len(validBody)},
			expectContains: `{"code":"10000","message":"Success","data":`,
		},
		{
			name:           "rpc error",
			mockErr:        errno.InternalServiceError,
			body:           &ut.Body{Body: strings.NewReader(validBody), Len: len(validBody)},
			expectContains: `{"code":"50001","message":"内部服务错误"`,
		},
		{
			name:           "bind error",
			body:           nil,
			expectContains: `{"code":"20001","message":"参数错误`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/common/signed-location-api-url", GetSignedLocationApiUrl)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			if tc.name != "bind error" {
				mockey.Mock(rpc.GetSignedLocationApiUrlRPC).To(func(ctx context.Context, req *common.GetSignedLocationApiUrlRequest) (string, map[string]string, error) {
					return tc.mockSignedURL, tc.mockHeaders, tc.mockErr
				}).Build()
			}

			res := ut.PerformRequest(router, consts.MethodPost, "/api/v1/common/signed-location-api-url", tc.body,
				ut.Header{Key: "Content-Type", Value: "application/json"})
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}
