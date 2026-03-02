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

func TestGetAllCloudSetting(t *testing.T) {
	type testCase struct {
		name               string  // 测试用例名称
		mockSettingJson    *[]byte // mock返回的设置 JSON 数据
		mockError          error   // mock返回的错误
		expectResult       *[]byte // 期望返回的结果
		expectError        string  // 期望的错误信息
		mockCommentedJson  string  // 模拟带注释的 JSON 数据
		mockCommentedError error   // 模拟去掉注释过程的错误
	}

	mockResult := []byte(`{"key": "value"}`)

	// 测试用例
	testCases := []testCase{
		{
			name:               "SuccessCase",
			mockSettingJson:    &mockResult,
			mockError:          nil,
			expectResult:       &mockResult,
			mockCommentedJson:  `{"key": "value"}`,
			mockCommentedError: nil,
		},
		{
			name:               "FileNotFound",
			mockSettingJson:    nil,
			mockError:          fmt.Errorf("file not found"),
			expectResult:       nil,
			expectError:        "VersionService.GetAllCloudSetting error:file not found",
			mockCommentedJson:  "",
			mockCommentedError: nil,
		},
		{
			name:               "RemoveCommentsError",
			mockSettingJson:    &mockResult,
			mockError:          nil,
			expectResult:       nil,
			expectError:        "VersionService.GetAllCloudSetting error:invalid JSON format",
			mockCommentedJson:  "",
			mockCommentedError: fmt.Errorf("invalid JSON format"),
		},
	}

	defer mockey.UnPatchAll() // 清理所有mock

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// Mock upyun.URlGetFile 方法
			mockey.Mock(upyun.URlGetFile).Return(tc.mockSettingJson, tc.mockError).Build()
			mockey.Mock(upyun.JoinFileName).To(func(filename string) string {
				return filename
			}).Build()

			// Mock getJSONWithoutComments 方法
			mockey.Mock(getJSONWithoutComments).Return(tc.mockCommentedJson, tc.mockCommentedError).Build()

			// 初始化UrlService实例
			versionService := &VersionService{}

			// 调用方法
			result, err := versionService.GetAllCloudSetting()

			if tc.expectError != "" {
				// 如果期望抛错，检查错误信息
				assert.ErrorContains(t, err, tc.expectError)
				assert.Equal(t, tc.expectResult, result)
			} else {
				// 如果不期望抛错，验证结果
				assert.Nil(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}
