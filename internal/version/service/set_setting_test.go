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
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/version"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestSetSetting(t *testing.T) {
	type testCase struct {
		name            string                   // 测试用例名称
		mockCheckPwd    bool                     // 模拟 CheckPwd 的返回值
		mockUploadError error                    // 模拟 URlUploadFile 的错误
		request         *version.SetCloudRequest // 输入的请求
		expectError     string                   // 期望的错误信息
	}

	testCases := []testCase{
		{
			name:            "ValidPasswordAndSuccessfulUpload",
			mockCheckPwd:    true,
			mockUploadError: nil,
			request: &version.SetCloudRequest{
				Password: "validpassword",
				Setting:  "{\"key\": \"value\"}",
			},
		},
		{
			name:            "InvalidPassword",
			mockCheckPwd:    false,
			mockUploadError: nil,
			request: &version.SetCloudRequest{
				Password: "invalidpassword",
				Setting:  "{\"key\": \"value\"}",
			},
			expectError: "[401] authorization failed", // 假设 buildAuthFailedError 返回这个错误信息
		},
		{
			name:            "ValidPasswordButUploadFails",
			mockCheckPwd:    true,
			mockUploadError: fmt.Errorf("upload failed"),
			request: &version.SetCloudRequest{
				Password: "validpassword",
				Setting:  "{\"key\": \"value\"}",
			},
			expectError: "upload failed",
		},
	}

	defer mockey.UnPatchAll() // 清理所有mock

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// Mock utils.CheckPwd 方法
			mockey.Mock(utils.CheckPwd).Return(tc.mockCheckPwd).Build()

			// Mock upyun.URlUploadFile 方法
			mockey.Mock(upyun.URlUploadFile).Return(tc.mockUploadError).Build()
			mockey.Mock(upyun.JoinFileName).To(func(filename string) string {
				return filename
			}).Build()

			// 初始化 UrlService 实例
			versionService := &VersionService{}

			// 调用方法
			err := versionService.SetSetting(tc.request)

			if tc.expectError != "" {
				// 如果期望抛错，检查错误信息
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				// 如果不期望抛错，验证结果
				assert.Nil(t, err)
			}
		})
	}
}
