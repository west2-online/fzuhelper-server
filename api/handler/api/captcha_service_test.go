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
	"mime/multipart"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/captcha"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// 测试数据常量
//
//nolint:lll
const (
	ImageBase64 = `data:image/png;base64,Qk2mCAAAAAAAADYAAAAoAAAASAAAAAoAAAABABgAAAAAAHAIAAASCwAAEgsAAAAAAAAAAAAA+vr/+vr/+vr/lgD6lgD6lgD6lgD6+vr/+vr/+vr/+vr/+vr/+vr/ljIAljIAljIAljIA+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYAAJYAAJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIAljIA+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/ljIAljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/AJYAAJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/AACWAACW+vr/AACW+vr/+vr/+vr/lgD6+vr/lgD6lgD6lgD6lgD6lgD6+vr/+vr/+vr/+vr/+vr/ljIAljIAljIAljIA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACWAACWAACW+vr/+vr/+vr/+vr/+vr/`
)

// buildValidateCodeForm 构建验证码验证的 form 数据
func buildValidateCodeForm() (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("image", ImageBase64)
	_ = w.Close()
	return &buf, w.FormDataContentType()
}

// buildValidateCodeForAndroidForm 构建 Android 验证码验证的 form 数据
func buildValidateCodeForAndroidForm() (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("validateCode", ImageBase64)
	_ = w.Close()
	return &buf, w.FormDataContentType()
}

func TestValidateCode(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		buildForm      func() (*bytes.Buffer, string)
		mockRPCError   error
		expectError    bool
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/user/validate-code",
			buildForm:      buildValidateCodeForm,
			mockRPCError:   nil,
			expectError:    false,
			expectContains: `"code":"10000","message":"Success","data":"104"`,
		},
		{
			name:           "invalid_param",
			url:            "/api/v1/user/validate-code",
			buildForm:      buildValidateCodeForAndroidForm,
			mockRPCError:   nil,
			expectError:    true,
			expectContains: `"code":"20001","message":"参数错误`,
		},
		{
			name:           "rpc_error",
			url:            "/api/v1/user/validate-code",
			buildForm:      buildValidateCodeForm,
			mockRPCError:   errno.InternalServiceError,
			expectError:    true,
			expectContains: `"code":"50001","message":"内部服务错误"`,
		},
	}
	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/user/validate-code", ValidateCode)
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.ValidateCodeRPC).To(func(ctx context.Context, req *captcha.ValidateCodeRequest) (string, error) {
				if tc.mockRPCError != nil {
					return "", tc.mockRPCError
				}
				return "104", nil
			}).Build()

			buf, contentType := tc.buildForm()
			result := ut.PerformRequest(router, "POST", tc.url,
				&ut.Body{Body: buf, Len: buf.Len()},
				ut.Header{Key: "Content-Type", Value: contentType})
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestValidateCodeForAndroid(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		buildForm      func() (*bytes.Buffer, string)
		mockRPCError   error
		expectError    bool
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/login/validateCode",
			buildForm:      buildValidateCodeForAndroidForm,
			mockRPCError:   nil,
			expectError:    false,
			expectContains: `"code":"200","message":"104"`,
		},
		{
			name:           "invalid_param",
			url:            "/api/login/validateCode",
			buildForm:      buildValidateCodeForm,
			mockRPCError:   nil,
			expectError:    true,
			expectContains: `"code":"20001","message":"参数错误`,
		},
		{
			name:           "rpc_error",
			url:            "/api/login/validateCode",
			buildForm:      buildValidateCodeForAndroidForm,
			mockRPCError:   errno.InternalServiceError,
			expectError:    true,
			expectContains: `"code":"50001","message":"内部服务错误"`,
		},
	}
	router := route.NewEngine(&config.Options{})
	router.POST("/api/login/validateCode", ValidateCodeForAndroid)
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.ValidateCodeForAndroidRPC).To(func(ctx context.Context, req *captcha.ValidateCodeForAndroidRequest) (string, error) {
				if tc.mockRPCError != nil {
					return "", tc.mockRPCError
				}
				return "104", nil
			}).Build()

			buf, contentType := tc.buildForm()
			result := ut.PerformRequest(router, "POST", tc.url,
				&ut.Body{Body: buf, Len: buf.Len()},
				ut.Header{Key: "Content-Type", Value: contentType})
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}
