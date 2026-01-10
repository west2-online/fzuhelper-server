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

	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/captcha"
	captchapkg "github.com/west2-online/fzuhelper-server/pkg/captcha"
)

func init() {
	// 测试前初始化验证码模板
	if err := captchapkg.Init(); err != nil {
		panic(err)
	}
}

// //nolint:lll
const validDataURL = "data:image/png;base64,Qk2mCAAAAAAAADYAAAAoAAAASAAAAAoAAAABABgAAAAAAHAIAAASCwAAEgsAAAAAAAAAAAAA+vr/+vr/+vr/lgD6lgD6lgD6lgD6+vr/+vr/+vr/+vr/+vr/+vr/ljIAljIAljIAljIA+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYAAJYAAJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6lgD6+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/ljIAljIA+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/+vr/+vr/ljIA+vr/+vr/ljIAljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/AACW+vr/+vr/+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/lgD6+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/ljIA+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4+vr/+vr/+vr/lih4+vr/+vr/+vr/+vr/+vr/AJYAAJYAAJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACW+vr/+vr/AACWAACW+vr/AACW+vr/+vr/+vr/lgD6+vr/lgD6lgD6lgD6lgD6lgD6+vr/+vr/+vr/+vr/+vr/ljIAljIAljIAljIA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/lih4lih4lih4lih4lih4lih4+vr/+vr/+vr/+vr/+vr/+vr/AJYA+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/+vr/AACWAACWAACWAACW+vr/+vr/+vr/+vr/+vr/"

func TestValidateCode(t *testing.T) {
	type testCase struct {
		name        string                       // 测试用例名称
		request     *captcha.ValidateCodeRequest // 请求参数
		expectData  int                          // 验证码结果
		expectError bool                         // 是否期望抛出错误
	}

	testCases := []testCase{
		{
			name:        "success",
			request:     &captcha.ValidateCodeRequest{Image: validDataURL},
			expectData:  104,
			expectError: false,
		},
		{
			name:        "invalid_base64",
			request:     &captcha.ValidateCodeRequest{Image: "not-base64"},
			expectData:  0,
			expectError: true,
		},
		{
			name:        "malformed_image",
			request:     &captcha.ValidateCodeRequest{Image: "data:image/png;base64,YWJjZA=="},
			expectData:  0,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			captchaService := &CaptchaService{}
			data, err := captchaService.ValidateCaptcha(&tc.request.Image)
			if tc.expectError {
				// 如果期望抛错，检查错误信息
				assert.Error(t, err)
			} else {
				// 如果不期望抛错，验证结果
				assert.Nil(t, err)
				assert.Equal(t, tc.expectData, data)
			}
		})
	}
}
