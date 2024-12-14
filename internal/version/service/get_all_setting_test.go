package service

import (
	"fmt"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCloudSetting(t *testing.T) {
	type testCase struct {
		name               string  // 测试用例名称
		mockSettingJson    *[]byte // mock返回的设置 JSON 数据
		mockError          error   // mock返回的错误
		expectedResult     *[]byte // 期望返回的结果
		expectingError     bool    // 是否期望抛出错误
		expectedErrorInfo  string  // 期望的错误信息
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
			expectedResult:     &mockResult,
			expectingError:     false,
			expectedErrorInfo:  "",
			mockCommentedJson:  `{"key": "value"}`,
			mockCommentedError: nil,
		},
		{
			name:               "FileNotFound",
			mockSettingJson:    nil,
			mockError:          fmt.Errorf("file not found"),
			expectedResult:     nil,
			expectingError:     true,
			expectedErrorInfo:  "VersionService.GetAllCloudSetting error:file not found",
			mockCommentedJson:  "",
			mockCommentedError: nil,
		},
		{
			name:               "RemoveCommentsError",
			mockSettingJson:    &mockResult,
			mockError:          nil,
			expectedResult:     nil,
			expectingError:     true,
			expectedErrorInfo:  "VersionService.GetAllCloudSetting error:invalid JSON format",
			mockCommentedJson:  "",
			mockCommentedError: fmt.Errorf("invalid JSON format"),
		},
	}

	defer mockey.UnPatchAll() // 清理所有mock

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// Mock upyun.URlGetFile 方法
			mockey.Mock(upyun.URlGetFile).To(func(filename string) (*[]byte, error) {
				return tc.mockSettingJson, tc.mockError
			}).Build()
			mockey.Mock(upyun.JoinFileName).To(func(filename string) string {
				return filename
			}).Build()

			// Mock getJSONWithoutComments 方法
			mockey.Mock(getJSONWithoutComments).To(func(json string) (string, error) {
				if tc.mockCommentedError != nil {
					return "", tc.mockCommentedError
				}
				return tc.mockCommentedJson, nil
			}).Build()

			// 初始化UrlService实例
			versionService := &VersionService{}

			// 调用方法
			result, err := versionService.GetAllCloudSetting()

			if tc.expectingError {
				// 如果期望抛错，检查错误信息
				assert.NotNil(t, err)
				assert.EqualError(t, err, tc.expectedErrorInfo)
				assert.Equal(t, tc.expectedResult, result)
			} else {
				// 如果不期望抛错，验证结果
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
