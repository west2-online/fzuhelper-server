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

package custom

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/version"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func buildLoginForm(password string) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("password", password)
	_ = w.Close()
	return &buf, w.FormDataContentType()
}

func buildUploadForm(versionStr, code, url, feature, typeStr, password, force string) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("version", versionStr)
	_ = w.WriteField("code", code)
	_ = w.WriteField("url", url)
	_ = w.WriteField("feature", feature)
	_ = w.WriteField("type", typeStr)
	_ = w.WriteField("password", password)
	_ = w.WriteField("force", force)
	_ = w.Close()
	return &buf, w.FormDataContentType()
}

func buildUploadParamsForm(password string) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("password", password)
	_ = w.Close()
	return &buf, w.FormDataContentType()
}

func ptrStr(s string) *string { return &s }

func TestAPILogin(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		password       string
		mockRPCErr     error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/url/login",
			password:       "correct_password",
			expectContains: `200`,
		},
		{
			name:           "unauthorized",
			url:            "/api/v1/url/login",
			password:       "wrong_password",
			mockRPCErr:     errno.NewErrNo(http.StatusUnauthorized, "unauthorized"),
			expectContains: urlCustomErrorMsg,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/url/login",
			password:       "test_password",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `"code":"50001","message":"内部服务错误"`,
		},
		{
			name:           "bind error - missing password",
			url:            "/api/v1/url/login",
			expectContains: `"code":"20001","message":"参数错误,`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/url/login", APILogin)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.LoginRPC).To(func(ctx context.Context, req *version.LoginRequest) error {
				return tc.mockRPCErr
			}).Build()

			buf, contentType := buildLoginForm(tc.password)
			res := ut.PerformRequest(router, consts.MethodPost, tc.url,
				&ut.Body{Body: buf, Len: buf.Len()},
				ut.Header{Key: "Content-Type", Value: contentType})
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestUploadVersionInfo(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		version        string
		code           string
		urlStr         string
		feature        string
		typeStr        string
		password       string
		force          string
		mockRPCErr     error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/url/api/upload",
			version:        "1.0.0",
			code:           "100",
			urlStr:         "https://example.com/app.apk",
			feature:        "new features",
			typeStr:        "release",
			password:       "correct_password",
			force:          "true",
			expectContains: "200",
		},
		{
			name:           "bind error",
			url:            "/api/v1/url/api/upload",
			version:        "1.0.0",
			code:           "100",
			urlStr:         "https://example.com/app.apk",
			feature:        "new features",
			typeStr:        "release",
			force:          "true",
			expectContains: `"code":"20001","message":"参数错误,`,
		},
		{
			name:           "unauthorized",
			url:            "/api/v1/url/api/upload",
			version:        "1.0.0",
			code:           "100",
			urlStr:         "https://example.com/app.apk",
			feature:        "new features",
			typeStr:        "release",
			password:       "wrong_password",
			force:          "true",
			mockRPCErr:     errno.NewErrNo(http.StatusUnauthorized, "unauthorized"),
			expectContains: urlCustomErrorMsg,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/url/api/upload",
			version:        "1.0.0",
			code:           "100",
			urlStr:         "https://example.com/app.apk",
			feature:        "new features",
			typeStr:        "release",
			password:       "test_password",
			force:          "true",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `"code":"50001","message":"内部服务错误"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/url/api/upload", UploadVersionInfo)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.UploadVersionRPC).To(func(ctx context.Context, req *version.UploadRequest) error {
				return tc.mockRPCErr
			}).Build()

			buf, contentType := buildUploadForm(tc.version, tc.code, tc.urlStr, tc.feature, tc.typeStr, tc.password, tc.force)
			res := ut.PerformRequest(router, consts.MethodPost, tc.url,
				&ut.Body{Body: buf, Len: buf.Len()},
				ut.Header{Key: "Content-Type", Value: contentType})
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetUploadParams(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		password       string
		mockPolicy     string
		mockAuth       string
		mockRPCErr     error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/url/api/uploadparams",
			password:       "correct_password",
			mockPolicy:     "test_policy",
			mockAuth:       "test_auth",
			expectContains: "test_policy",
		},
		{
			name:           "bind error",
			url:            "/api/v1/url/api/uploadparams",
			expectContains: `"code":"20001","message":"参数错误,`,
		},
		{
			name:           "unauthorized",
			url:            "/api/v1/url/api/uploadparams",
			password:       "wrong_password",
			mockRPCErr:     errno.NewErrNo(http.StatusUnauthorized, "unauthorized"),
			expectContains: urlCustomErrorMsg,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/url/api/uploadparams",
			password:       "test_password",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `"code":"50001","message":"内部服务错误"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/url/api/uploadparams", GetUploadParams)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.UploadParamsRPC).To(func(ctx context.Context, req *version.UploadParamsRequest) (*string, *string, error) {
				return &tc.mockPolicy, &tc.mockAuth, tc.mockRPCErr
			}).Build()

			buf, contentType := buildUploadParamsForm(tc.password)
			res := ut.PerformRequest(router, consts.MethodPost, tc.url,
				&ut.Body{Body: buf, Len: buf.Len()},
				ut.Header{Key: "Content-Type", Value: contentType})
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetReleaseVersionModify(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *version.GetReleaseVersionResponse
		mockRPCErr     error
		expectContains string
	}

	testCases := []testCase{
		{
			name: "success",
			url:  "/api/v1/url/version.json",
			mockResp: &version.GetReleaseVersionResponse{
				Version: ptrStr("1.0.0"),
				Url:     ptrStr("https://example.com/app.apk"),
				Code:    ptrStr("100"),
				Feature: ptrStr("new features"),
			},
			expectContains: "1.0.0",
		},
		{
			name:           "rpc error",
			url:            "/api/v1/url/version.json",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `"code":"50001","message":"内部服务错误"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/url/version.json", GetReleaseVersionModify)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetReleaseVersionRPC).To(func(ctx context.Context, req *version.GetReleaseVersionRequest) (*version.GetReleaseVersionResponse, error) {
				return tc.mockResp, tc.mockRPCErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetBetaVersionModify(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *version.GetBetaVersionResponse
		mockRPCErr     error
		expectContains string
	}

	testCases := []testCase{
		{
			name: "success",
			url:  "/api/v1/url/versionbeta.json",
			mockResp: &version.GetBetaVersionResponse{
				Version: ptrStr("1.0.0-beta"),
				Url:     ptrStr("https://example.com/app-beta.apk"),
				Code:    ptrStr("99"),
				Feature: ptrStr("beta features"),
			},
			expectContains: "1.0.0-beta",
		},
		{
			name:           "rpc error",
			url:            "/api/v1/url/versionbeta.json",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `"code":"50001","message":"内部服务错误"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/url/versionbeta.json", GetBetaVersionModify)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetBetaVersionRPC).To(func(ctx context.Context, req *version.GetBetaVersionRequest) (*version.GetBetaVersionResponse, error) {
				return tc.mockResp, tc.mockRPCErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}
