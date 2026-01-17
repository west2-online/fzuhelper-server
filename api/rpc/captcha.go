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

package rpc

import (
	"context"

	"github.com/west2-online/fzuhelper-server/kitex_gen/captcha"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitCaptchaRPC() {
	c, err := client.InitCaptchaRPC()
	if err != nil {
		logger.Fatalf("api.rpc.captcha InitCaptchaRPC failed, err is %v", err)
	}
	captchaClient = *c
}

func ValidateCodeRPC(ctx context.Context, req *captcha.ValidateCodeRequest) (string, error) {
	resp, err := captchaClient.ValidateCode(ctx, req)
	if err != nil {
		logger.Errorf("ValidateCodeRPC: RPC called failed: %v", err.Error())
		return "", errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return "", errno.BizError.WithMessage("验证码验证失败: " + resp.Base.Msg)
	}
	return resp.Data, nil
}

func ValidateCodeForAndroidRPC(ctx context.Context, req *captcha.ValidateCodeForAndroidRequest) (string, error) {
	resp, err := captchaClient.ValidateCodeForAndroid(ctx, req)
	if err != nil {
		logger.Errorf("ValidateCodeForAndroidRPC: RPC called failed: %v", err.Error())
		return "", errno.InternalServiceError.WithError(err)
	}
	if resp.Code != "200" {
		return "", errno.BizError.WithMessage("验证码验证失败: " + resp.Message)
	}
	return resp.Message, nil
}
