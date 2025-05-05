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
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestUploadVersion(t *testing.T) {
	type testCase struct {
		name              string                 // 测试用例名称
		mockCheckPwd      bool                   // 模拟 CheckPwd 的返回值
		mockUploadError   error                  // 模拟 URlUploadFile 的错误
		mockMarshalError  error                  // 模拟 JSON Marshal 的错误
		request           *version.UploadRequest // 请求参数
		expectedError     bool                   // 是否期望抛出错误
		expectedErrorInfo string                 // 期望的错误信息
	}

	// 测试用例
	testCases := []testCase{
		{
			name:             "ValidPasswordAndUploadRelease",
			mockCheckPwd:     true,
			mockUploadError:  nil,
			mockMarshalError: nil,
			request: &version.UploadRequest{
				Password: "validpassword",
				Version:  "1.0.0",
				Code:     "633001",
				Url:      "http://example.com/release.apk",
				Feature:  "New features",
				Type:     apkTypeRelease,
			},
			expectedError: false,
		},
		{
			name:             "ValidPasswordAndUploadBeta",
			mockCheckPwd:     true,
			mockUploadError:  nil,
			mockMarshalError: nil,
			request: &version.UploadRequest{
				Password: "validpassword",
				Version:  "1.0.1-beta",
				Code:     "633001",
				Url:      "http://example.com/beta.apk",
				Feature:  "Beta features",
				Type:     apkTypeBeta,
			},
			expectedError: false,
		},
		{
			name:             "ValidPasswordAndUploadAlpha",
			mockCheckPwd:     true,
			mockUploadError:  nil,
			mockMarshalError: nil,
			request: &version.UploadRequest{
				Password: "validpassword",
				Version:  "1.0.2-alpha",
				Code:     "633001",
				Url:      "http://example.com/alpha.apk",
				Feature:  "Alpha features",
				Type:     apkTypeAlpha,
			},
			expectedError: false,
		},
		{
			name:             "InvalidPassword",
			mockCheckPwd:     false,
			mockUploadError:  nil,
			mockMarshalError: nil,
			request: &version.UploadRequest{
				Password: "invalidpassword",
				Version:  "1.0.0",
				Code:     "633001",
				Url:      "http://example.com/release.apk",
				Feature:  "New features",
				Type:     apkTypeRelease,
			},
			expectedError:     true,
			expectedErrorInfo: "[401] authorization failed",
		},
		{
			name:             "InvalidApkType",
			mockCheckPwd:     true,
			mockUploadError:  nil,
			mockMarshalError: nil,
			request: &version.UploadRequest{
				Password: "validpassword",
				Version:  "1.0.0",
				Code:     "633001",
				Url:      "http://example.com/release.apk",
				Feature:  "New features",
				Type:     "invalidType",
			},
			expectedError:     true,
			expectedErrorInfo: errno.ParamError.ErrorMsg,
		},
	}

	defer mockey.UnPatchAll() // 清理所有mock

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// Mock utils.CheckPwd 方法
			mockey.Mock(utils.CheckPwd).To(func(password string) bool {
				return tc.mockCheckPwd
			}).Build()

			// Mock upyun.URlUploadFile 方法
			mockey.Mock(upyun.URlUploadFile).To(func(data []byte, filename string) error {
				return tc.mockUploadError
			}).Build()
			mockey.Mock(upyun.JoinFileName).To(func(filename string) string {
				return filename
			}).Build()

			// 初始化 UrlService 实例
			versionService := &VersionService{}

			// 调用方法
			err := versionService.UploadVersion(tc.request)

			if tc.expectedError {
				// 如果期望抛错，检查错误信息
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrorInfo)
			} else {
				// 如果不期望抛错，验证结果
				assert.Nil(t, err)
			}
		})
	}
}
