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

package service

import (
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/version"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestUploadParams(t *testing.T) {
	type testCase struct {
		name                string                       // 测试用例名称
		mockCheckPwd        bool                         // 模拟 CheckPwd 的返回值
		mockPolicy          string                       // 模拟 GetPolicy 的返回值
		mockAuthorization   string                       // 模拟 SignStr 的返回值
		request             *version.UploadParamsRequest // 请求参数
		expectPolicy        string                       // 期望的 policy 返回值
		expectAuthorization string                       // 期望的 authorization 返回值
		expectError         string                       // 期望的错误信息
	}

	// 测试用例
	testCases := []testCase{
		{
			name:              "ValidPassword",
			mockCheckPwd:      true,
			mockPolicy:        "mockPolicy",
			mockAuthorization: "mockAuthorization",
			request: &version.UploadParamsRequest{
				Password: "validpassword",
			},
			expectPolicy:        "mockPolicy",
			expectAuthorization: "mockAuthorization",
			expectError:         "",
		},
		{
			name:              "InvalidPassword",
			mockCheckPwd:      false,
			mockPolicy:        "",
			mockAuthorization: "",
			request: &version.UploadParamsRequest{
				Password: "invalidpassword",
			},
			expectPolicy:        "",
			expectAuthorization: "",
			expectError:         "[401] authorization failed", // 假设 buildAuthFailedError 返回这个错误信息
		},
	}

	defer mockey.UnPatchAll() // 清理所有mock

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// Mock utils.CheckPwd 方法
			mockey.Mock(utils.CheckPwd).Return(tc.mockCheckPwd).Build()

			// Mock upyun.GetPolicy 方法
			mockey.Mock(upyun.GetPolicy).Return(tc.mockPolicy).Build()

			// Mock upyun.SignStr 方法
			mockey.Mock(upyun.SignStr).Return(tc.mockAuthorization).Build()

			// 初始化 UrlService 实例
			versionService := &VersionService{}

			// 调用方法
			policy, authorization, err := versionService.UploadParams(tc.request)

			if tc.expectError != "" {
				// 如果期望抛错，检查错误信息
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				assert.Equal(t, tc.expectPolicy, policy)
				assert.Equal(t, tc.expectAuthorization, authorization)
			} else {
				// 如果不期望抛错，验证结果
				assert.Nil(t, err)
				assert.Equal(t, tc.expectPolicy, policy)
				assert.Equal(t, tc.expectAuthorization, authorization)
			}
		})
	}
}
