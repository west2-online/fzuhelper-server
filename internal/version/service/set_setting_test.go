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
		name              string                   // 测试用例名称
		mockCheckPwd      bool                     // 模拟 CheckPwd 的返回值
		mockUploadError   error                    // 模拟 URlUploadFile 的错误
		request           *version.SetCloudRequest // 输入的请求
		expectedError     bool                     // 是否期望抛出错误
		expectedErrorInfo string                   // 期望的错误信息
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
			expectedError: false,
		},
		{
			name:            "InvalidPassword",
			mockCheckPwd:    false,
			mockUploadError: nil,
			request: &version.SetCloudRequest{
				Password: "invalidpassword",
				Setting:  "{\"key\": \"value\"}",
			},
			expectedError:     true,
			expectedErrorInfo: "[401] authorization failed", // 假设 buildAuthFailedError 返回这个错误信息
		},
		{
			name:            "ValidPasswordButUploadFails",
			mockCheckPwd:    true,
			mockUploadError: fmt.Errorf("upload failed"),
			request: &version.SetCloudRequest{
				Password: "validpassword",
				Setting:  "{\"key\": \"value\"}",
			},
			expectedError:     true,
			expectedErrorInfo: "upload failed",
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
			err := versionService.SetSetting(tc.request)

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
