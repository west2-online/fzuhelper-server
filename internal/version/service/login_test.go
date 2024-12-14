package service

import (
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/version"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestLogin(t *testing.T) {
	type testCase struct {
		name             string                // 测试用例名称
		mockCheckPwd     bool                  // 模拟 CheckPwd 的返回值
		request          *version.LoginRequest // 输入的登录请求
		expectedError    bool                  // 是否期望抛出错误
		expectedErrorMsg string                // 期望的错误类型或信息
	}

	testCases := []testCase{
		{
			name:         "ValidPassword",
			mockCheckPwd: true,
			request: &version.LoginRequest{
				Password: "validpassword",
			},
			expectedError: false,
		},
		{
			name:         "InvalidPassword",
			mockCheckPwd: false,
			request: &version.LoginRequest{
				Password: "invalidpassword",
			},
			expectedError:    true,
			expectedErrorMsg: "[401] authorization failed",
		},
	}

	defer mockey.UnPatchAll() // 清理所有mock

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// Mock utils.CheckPwd 方法
			mockey.Mock(utils.CheckPwd).To(func(password string) bool {
				return tc.mockCheckPwd
			}).Build()

			// 初始化 UrlService 实例
			versionService := &VersionService{}

			// 调用方法
			err := versionService.Login(tc.request)

			if tc.expectedError {
				// 如果期望抛错，检查错误信息
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrorMsg)
			} else {
				// 如果不期望抛错，验证结果
				assert.Nil(t, err)
			}
		})
	}
}
