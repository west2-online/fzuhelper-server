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
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/version"
)

// 辅助函数
func ptrStr(s string) *string { return &s }
func ptrBool(b bool) *bool    { return &b }

func TestLogin(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v2/url/login?password=testpass",
			mockErr:        nil,
			expectContains: `"code":"10000"`,
		},
		{
			name:           "param error - missing password",
			url:            "/api/v2/url/login",
			mockErr:        nil,
			expectContains: `"code":"20001"`,
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/login?password=wrongpass",
			mockErr:        errors.New("authentication failed"),
			expectContains: `"code":"50001"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v2/url/login", Login)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.LoginRPC).To(func(ctx context.Context, req *version.LoginRequest) error {
				return tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPost, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestUploadVersion(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v2/url/api/upload?version=1.0&code=1&url=http://test&feature=test&type=release&password=pass&force=false",
			mockErr:        nil,
			expectContains: `"code":"10000"`,
		},
		{
			name:           "param error - missing required fields",
			url:            "/api/v2/url/api/upload",
			mockErr:        nil,
			expectContains: `"code":"20001"`,
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/api/upload?version=1.0&code=1&url=http://test&feature=test&type=release&password=pass&force=false",
			mockErr:        errors.New("upload failed"),
			expectContains: `"code":"50001"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v2/url/api/upload", UploadVersion)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.UploadVersionRPC).To(func(ctx context.Context, req *version.UploadRequest) error {
				return tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPost, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestUploadParams(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockPolicy     *string
		mockAuth       *string
		mockErr        error
		expectContains string
	}

	policy := "test_policy"
	auth := "test_auth"

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v2/url/api/uploadparams?password=testpass",
			mockPolicy:     &policy,
			mockAuth:       &auth,
			mockErr:        nil,
			expectContains: `"policy":`,
		},
		{
			name:           "param error - missing password",
			url:            "/api/v2/url/api/uploadparams",
			mockPolicy:     nil,
			mockAuth:       nil,
			mockErr:        nil,
			expectContains: `"code":"20001"`,
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/api/uploadparams?password=wrongpass",
			mockPolicy:     nil,
			mockAuth:       nil,
			mockErr:        errors.New("params error"),
			expectContains: `"code":"50001"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v2/url/api/uploadparams", UploadParams)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.UploadParamsRPC).To(func(ctx context.Context, req *version.UploadParamsRequest) (*string, *string, error) {
				return tc.mockPolicy, tc.mockAuth, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPost, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestDownloadReleaseApk(t *testing.T) {
	type testCase struct {
		name       string
		url        string
		mockResp   *string
		mockErr    error
		expectCode int
	}

	releaseApkUrl := "http://example.com/release.apk"

	testCases := []testCase{
		{
			name:       "success - redirect",
			url:        "/api/v2/url/release.apk",
			mockResp:   &releaseApkUrl,
			mockErr:    nil,
			expectCode: consts.StatusFound,
		},
		{
			name:       "rpc error",
			url:        "/api/v2/url/release.apk",
			mockResp:   nil,
			mockErr:    errors.New("download failed"),
			expectCode: consts.StatusOK,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/url/release.apk", DownloadReleaseApk)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.DownloadReleaseApkRPC).To(func(ctx context.Context, req *version.DownloadReleaseApkRequest) (*string, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, tc.expectCode, res.Result().StatusCode())
		})
	}
}

func TestDownloadBetaApk(t *testing.T) {
	type testCase struct {
		name       string
		url        string
		mockResp   *string
		mockErr    error
		expectCode int
	}

	betaApkUrl := "http://example.com/beta.apk"

	testCases := []testCase{
		{
			name:       "success - redirect",
			url:        "/api/v2/url/beta.apk",
			mockResp:   &betaApkUrl,
			mockErr:    nil,
			expectCode: consts.StatusFound,
		},
		{
			name:       "rpc error",
			url:        "/api/v2/url/beta.apk",
			mockResp:   nil,
			mockErr:    errors.New("download failed"),
			expectCode: consts.StatusOK,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/url/beta.apk", DownloadBetaApk)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.DownloadBetaApkRPC).To(func(ctx context.Context, req *version.DownloadBetaApkRequest) (*string, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, tc.expectCode, res.Result().StatusCode())
		})
	}
}

func TestGetReleaseVersion(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *version.GetReleaseVersionResponse
		mockErr        error
		expectContains string
	}

	releaseVersion := &version.GetReleaseVersionResponse{
		Version: ptrStr("1.0.0"),
		Url:     ptrStr("http://example.com/app.apk"),
		Code:    ptrStr("1"),
		Feature: ptrStr("new feature"),
		Force:   ptrBool(true),
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v2/url/version.json",
			mockResp:       releaseVersion,
			mockErr:        nil,
			expectContains: `"version":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/version.json",
			mockResp:       nil,
			mockErr:        errors.New("version error"),
			expectContains: `"code":"50001"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/url/version.json", GetReleaseVersion)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetReleaseVersionRPC).To(func(ctx context.Context, req *version.GetReleaseVersionRequest) (*version.GetReleaseVersionResponse, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetBetaVersion(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *version.GetBetaVersionResponse
		mockErr        error
		expectContains string
	}

	betaVersion := &version.GetBetaVersionResponse{
		Version: ptrStr("2.0.0-beta"),
		Url:     ptrStr("http://example.com/beta.apk"),
		Code:    ptrStr("2"),
		Feature: ptrStr("beta feature"),
		Force:   ptrBool(false),
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v2/url/versionbeta.json",
			mockResp:       betaVersion,
			mockErr:        nil,
			expectContains: `"version":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/versionbeta.json",
			mockResp:       nil,
			mockErr:        errors.New("version error"),
			expectContains: `"code":"50001"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/url/versionbeta.json", GetBetaVersion)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetBetaVersionRPC).To(func(ctx context.Context, req *version.GetBetaVersionRequest) (*version.GetBetaVersionResponse, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetSetting(t *testing.T) {
	type testCase struct {
		name       string
		url        string
		mockResp   *version.GetSettingResponse
		mockErr    error
		expectCode int
	}

	settingResp := &version.GetSettingResponse{
		CloudSetting: []byte(`{"key":"value"}`),
	}

	testCases := []testCase{
		{
			name:       "success",
			url:        "/api/v2/url/settings.php?account=test&version=1.0&beta=false&phone=android&isLogin=true&loginType=ldap",
			mockResp:   settingResp,
			mockErr:    nil,
			expectCode: consts.StatusOK,
		},
		{
			name:       "rpc error",
			url:        "/api/v2/url/settings.php?account=test&version=1.0&beta=false&phone=android&isLogin=true&loginType=ldap",
			mockResp:   nil,
			mockErr:    errors.New("setting error"),
			expectCode: consts.StatusOK,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/url/settings.php", GetSetting)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetSettingRPC).To(func(ctx context.Context, req *version.GetSettingRequest) (*version.GetSettingResponse, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, tc.expectCode, res.Result().StatusCode())
		})
	}
}

func TestGetTest(t *testing.T) {
	type testCase struct {
		name       string
		url        string
		mockResp   *version.GetTestResponse
		mockErr    error
		expectCode int
	}

	testResp := &version.GetTestResponse{
		CloudSetting: []byte(`{"test":"data"}`),
	}

	testCases := []testCase{
		{
			name:       "success",
			url:        "/api/v2/url/test?account=test&version=1.0&beta=false&phone=android&isLogin=true&loginType=ldap&setting=test",
			mockResp:   testResp,
			mockErr:    nil,
			expectCode: consts.StatusOK,
		},
		{
			name:       "rpc error",
			url:        "/api/v2/url/test?account=test&version=1.0&beta=false&phone=android&isLogin=true&loginType=ldap&setting=test",
			mockResp:   nil,
			mockErr:    errors.New("test error"),
			expectCode: consts.StatusOK,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v2/url/test", GetTest)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetTestRPC).To(func(ctx context.Context, req *version.GetTestRequest) (*version.GetTestResponse, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPost, tc.url, nil)
			assert.Equal(t, tc.expectCode, res.Result().StatusCode())
		})
	}
}

func TestGetCloud(t *testing.T) {
	type testCase struct {
		name       string
		url        string
		mockResp   *version.GetCloudResponse
		mockErr    error
		expectCode int
	}

	cloudResp := &version.GetCloudResponse{
		CloudSetting: []byte(`{"cloud":"config"}`),
	}

	testCases := []testCase{
		{
			name:       "success",
			url:        "/api/v2/url/getcloud",
			mockResp:   cloudResp,
			mockErr:    nil,
			expectCode: consts.StatusOK,
		},
		{
			name:       "rpc error",
			url:        "/api/v2/url/getcloud",
			mockResp:   nil,
			mockErr:    errors.New("cloud error"),
			expectCode: consts.StatusOK,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/url/getcloud", GetCloud)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetCloudRPC).To(func(ctx context.Context, req *version.GetCloudRequest) (*version.GetCloudResponse, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, tc.expectCode, res.Result().StatusCode())
		})
	}
}

func TestSetCloud(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v2/url/setcloud?password=testpass&setting=test_setting",
			mockErr:        nil,
			expectContains: `"code":"10000"`,
		},
		{
			name:           "param error - missing password",
			url:            "/api/v2/url/setcloud",
			mockErr:        nil,
			expectContains: `"code":"20001"`,
		},
		{
			name:           "rpc error",
			url:            "/api/v2/url/setcloud?password=wrongpass&setting=test_setting",
			mockErr:        errors.New("set cloud failed"),
			expectContains: `"code":"50001"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v2/url/setcloud", SetCloud)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.SetCloudRPC).To(func(ctx context.Context, req *version.SetCloudRequest) error {
				return tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPost, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetDump(t *testing.T) {
	type testCase struct {
		name       string
		url        string
		mockResp   *version.GetDumpResponse
		mockErr    error
		expectCode int
	}

	dumpData := &version.GetDumpResponse{
		Data: `{"dump":"data"}`,
	}

	testCases := []testCase{
		{
			name:       "success",
			url:        "/api/v2/url/dump",
			mockResp:   dumpData,
			mockErr:    nil,
			expectCode: consts.StatusOK,
		},
		{
			name:       "rpc error",
			url:        "/api/v2/url/dump",
			mockResp:   nil,
			mockErr:    errors.New("dump error"),
			expectCode: consts.StatusOK,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/url/dump", GetDump)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetDumpRPC).To(func(ctx context.Context, req *version.GetDumpRequest) (*version.GetDumpResponse, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, tc.expectCode, res.Result().StatusCode())
		})
	}
}

func TestAndroidGetVersion(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *version.AndroidGetVersionResponse
		mockErr        error
		expectContains string
	}

	androidVersionResp := &version.AndroidGetVersionResponse{
		Release: &model.Version{
			VersionCode: ptrStr("1.0.0"),
			Url:         ptrStr("http://example.com/release.apk"),
			Changelog:   ptrStr("release feature"),
			Force:       ptrBool(true),
		},
		Beta: &model.Version{
			VersionCode: ptrStr("2.0.0-beta"),
			Url:         ptrStr("http://example.com/beta.apk"),
			Changelog:   ptrStr("beta feature"),
			Force:       ptrBool(false),
		},
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v2/version/android",
			mockResp:       androidVersionResp,
			mockErr:        nil,
			expectContains: `"release":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v2/version/android",
			mockResp:       nil,
			mockErr:        errors.New("android version error"),
			expectContains: `"code":"50001"`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v2/version/android", AndroidGetVersion)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.AndroidVersionRPC).To(func(ctx context.Context, req *version.AndroidGetVersioneRequest) (*version.AndroidGetVersionResponse, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}
