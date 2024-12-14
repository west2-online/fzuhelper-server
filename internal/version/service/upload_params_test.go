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
		name                  string                       // 测试用例名称
		mockCheckPwd          bool                         // 模拟 CheckPwd 的返回值
		mockPolicy            string                       // 模拟 GetPolicy 的返回值
		mockAuthorization     string                       // 模拟 SignStr 的返回值
		request               *version.UploadParamsRequest // 请求参数
		expectedPolicy        string                       // 期望的 policy 返回值
		expectedAuthorization string                       // 期望的 authorization 返回值
		expectedError         bool                         // 是否期望抛出错误
		expectedErrorInfo     string                       // 期望的错误信息
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
			expectedPolicy:        "mockPolicy",
			expectedAuthorization: "mockAuthorization",
			expectedError:         false,
		},
		{
			name:              "InvalidPassword",
			mockCheckPwd:      false,
			mockPolicy:        "",
			mockAuthorization: "",
			request: &version.UploadParamsRequest{
				Password: "invalidpassword",
			},
			expectedPolicy:        "",
			expectedAuthorization: "",
			expectedError:         true,
			expectedErrorInfo:     "[401] authorization failed", // 假设 buildAuthFailedError 返回这个错误信息
		},
	}

	defer mockey.UnPatchAll() // 清理所有mock

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// Mock utils.CheckPwd 方法
			mockey.Mock(utils.CheckPwd).To(func(password string) bool {
				return tc.mockCheckPwd
			}).Build()

			// Mock upyun.GetPolicy 方法
			mockey.Mock(upyun.GetPolicy).To(func() string {
				return tc.mockPolicy
			}).Build()

			// Mock upyun.SignStr 方法
			mockey.Mock(upyun.SignStr).To(func(policy string) string {
				return tc.mockAuthorization
			}).Build()

			// 初始化 UrlService 实例
			versionService := &VersionService{}

			// 调用方法
			policy, authorization, err := versionService.UploadParams(tc.request)

			if tc.expectedError {
				// 如果期望抛错，检查错误信息
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrorInfo)
				assert.Equal(t, tc.expectedPolicy, policy)
				assert.Equal(t, tc.expectedAuthorization, authorization)
			} else {
				// 如果不期望抛错，验证结果
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedPolicy, policy)
				assert.Equal(t, tc.expectedAuthorization, authorization)
			}
		})
	}
}
