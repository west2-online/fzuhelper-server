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
	"net/http"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/cos"
)

func feedbackImageSuffix(file []byte) (string, bool) {
	switch http.DetectContentType(file) {
	case "image/jpeg":
		return "jpg", true
	case "image/png":
		return "png", true
	default:
		return "", false
	}
}

func (s *OAService) UploadFeedbackScreenshot(file []byte) (string, error) {
	if len(file) == 0 {
		return "", fmt.Errorf("service.UploadFeedbackScreenshot: 反馈截图不能为空")
	}
	if len(file) > constants.FeedbackScreenshotMaxSize {
		return "", fmt.Errorf("service.UploadFeedbackScreenshot: 反馈截图不能超过 5MB")
	}

	suffix, ok := feedbackImageSuffix(file)
	if !ok {
		return "", fmt.Errorf("service.UploadFeedbackScreenshot: 反馈截图仅支持 JPEG 和 PNG")
	}
	if s.cosClient == nil {
		return "", fmt.Errorf("service.UploadFeedbackScreenshot: 反馈文件存储未初始化")
	}

	url, remotePath, err := s.cosClient.GenerateFileName(cos.FeedbackImageCategory, suffix)
	if err != nil {
		return "", fmt.Errorf("service.UploadFeedbackScreenshot: 生成反馈截图名称失败: %w", err)
	}
	if err = s.cosClient.Upload(file, remotePath); err != nil {
		return "", fmt.Errorf("service.UploadFeedbackScreenshot: 上传反馈截图失败: %w", err)
	}
	return url, nil
}
