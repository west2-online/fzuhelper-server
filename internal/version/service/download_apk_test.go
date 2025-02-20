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

	"github.com/west2-online/fzuhelper-server/pkg/upyun"
)

func TestDownloadBetaApk(t *testing.T) {
	type testCase struct {
		name              string // 测试用例名称
		mockJsonBytes     []byte // mock返回的json数据
		mockError         error  // mock返回的错误
		expectedUrl       string // 期望返回的URL
		expectingError    bool   // 是否期望抛出错误
		expectedErrorInfo string // 期望的错误信息
	}

	// 测试用例
	testCases := []testCase{
		{
			name:              "SuccessCase",
			mockJsonBytes:     []byte(`{"Url":"https://example.com/beta.apk"}`),
			mockError:         nil,
			expectedUrl:       "https://example.com/beta.apk",
			expectingError:    false,
			expectedErrorInfo: "",
		},
		{
			name:              "FileNotFound",
			mockJsonBytes:     nil,
			mockError:         fmt.Errorf("file not found"),
			expectedUrl:       "",
			expectingError:    true,
			expectedErrorInfo: "VersionService.DownloadBetaApk error:file not found",
		},
		{
			name:              "UnmarshalError",
			mockJsonBytes:     []byte(`invalid json`),
			mockError:         nil,
			expectedUrl:       "",
			expectingError:    true,
			expectedErrorInfo: `VersionService.DownloadBetaApk error:"Syntax error at index`,
		},
	}

	defer mockey.UnPatchAll() // 清理所有mock

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// Mock upyun.URlGetFile 方法
			mockey.Mock(upyun.URlGetFile).To(func(url string) (*[]byte, error) {
				return &tc.mockJsonBytes, tc.mockError
			}).Build()
			mockey.Mock(upyun.JoinFileName).To(func(filename string) string {
				return filename
			}).Build()

			// 初始化UrlService实例
			urlService := &VersionService{}

			// 调用方法
			result, err := urlService.DownloadBetaApk()

			if tc.expectingError {
				// 如果期望抛错，检查错误信息
				assert.Contains(t, err.Error(), tc.expectedErrorInfo)
			} else {
				// 如果不期望抛错，验证结果
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUrl, result)
			}
		})
	}
}

func TestDownloadReleaseApk(t *testing.T) {
	type testCase struct {
		name              string // 测试用例名称
		mockJsonBytes     []byte // mock返回的json数据
		mockError         error  // mock返回的错误
		expectedUrl       string // 期望返回的URL
		expectingError    bool   // 是否期望抛出错误
		expectedErrorInfo string // 期望的错误信息
	}

	// 测试用例
	testCases := []testCase{
		{
			name:              "SuccessCase",
			mockJsonBytes:     []byte(`{"Url":"https://example.com/release.apk"}`),
			mockError:         nil,
			expectedUrl:       "https://example.com/release.apk",
			expectingError:    false,
			expectedErrorInfo: "",
		},
		{
			name:              "FileNotFound",
			mockJsonBytes:     nil,
			mockError:         fmt.Errorf("file not found"),
			expectedUrl:       "",
			expectingError:    true,
			expectedErrorInfo: "VersionService.DownloadReleaseApk error:file not found",
		},
		{
			name:              "UnmarshalError",
			mockJsonBytes:     []byte(`invalid json`),
			mockError:         nil,
			expectedUrl:       "",
			expectingError:    true,
			expectedErrorInfo: `VersionService.DownloadReleaseApk error:"Syntax error at index`,
		},
	}

	defer mockey.UnPatchAll() // 清理所有mock

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// Mock upyun.URlGetFile 方法
			mockey.Mock(upyun.URlGetFile).To(func(url string) (*[]byte, error) {
				return &tc.mockJsonBytes, tc.mockError
			}).Build()
			mockey.Mock(upyun.JoinFileName).To(func(filename string) string {
				return filename
			}).Build()

			// 初始化UrlService实例
			urlService := &VersionService{}

			// 调用方法
			result, err := urlService.DownloadReleaseApk()

			if tc.expectingError {
				// 如果期望抛错，检查错误信息
				assert.Contains(t, err.Error(), tc.expectedErrorInfo)
			} else {
				// 如果不期望抛错，验证结果
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUrl, result)
			}
		})
	}
}
