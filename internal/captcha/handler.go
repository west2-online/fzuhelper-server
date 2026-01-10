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

package captcha

import (
	"context"
	"fmt"

	"github.com/west2-online/fzuhelper-server/internal/captcha/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/captcha"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// CaptchaServiceImpl implements the last service interface defined in the IDL.
type CaptchaServiceImpl struct {
	ClientSet *base.ClientSet
}

func NewCaptchaService(clientSet *base.ClientSet) *CaptchaServiceImpl {
	return &CaptchaServiceImpl{
		ClientSet: clientSet,
	}
}

// ValidateCode implements the CaptchaServiceImpl interface.
func (s *CaptchaServiceImpl) ValidateCode(ctx context.Context, req *captcha.ValidateCodeRequest) (resp *captcha.ValidateCodeResponse, err error) {
	resp = new(captcha.ValidateCodeResponse)
	data, err := service.NewCaptchaService(ctx).ValidateCaptcha(&req.Image)
	if err != nil {
		logger.Infof("Captcha.ValidateCode: %v", err)
		return resp, nil
	}
	resp.Data = fmt.Sprint(data)
	return resp, nil
}

// ValidateCodeForAndroid implements the CaptchaServiceImpl interface.
func (s *CaptchaServiceImpl) ValidateCodeForAndroid(ctx context.Context, req *captcha.ValidateCodeForAndroidRequest) (resp *captcha.ValidateCodeForAndroidResponse, err error) { //nolint:lll
	resp = new(captcha.ValidateCodeForAndroidResponse)
	data, err := service.NewCaptchaService(ctx).ValidateCaptcha(&req.ValidateCode)
	resp.Message = fmt.Sprint(data)
	if err != nil {
		logger.Infof("Captcha.ValidateCodeForAndroid: %v", err)
		return resp, nil
	}
	return resp, nil
}
