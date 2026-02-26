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
	"github.com/west2-online/fzuhelper-server/pkg/captcha"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

const (
	// maxImageSize 图片 base64 字符串的最大长度(1MB)
	// 正常验证码图片 base64 编码后通常只有几 KB，1MB 是一个安全的上限
	maxImageSize = 1 << 20 // 1MB = 1048576 bytes
)

func (s *CaptchaService) ValidateCaptcha(reqImageData *string) (int, error) {
	if *reqImageData == "" {
		return 0, errno.ParamError.WithMessage("request image data is empty")
	}
	if len(*reqImageData) > maxImageSize {
		return 0, errno.ParamError.WithMessage("request image data is too large")
	}
	return captcha.ValidateLoginCode(*reqImageData)
}
