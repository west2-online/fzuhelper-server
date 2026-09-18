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
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/cos"
)

func (s *OAService) UploadFeedbackLog(file []byte) (string, error) {
	if len(file) == 0 {
		return "", fmt.Errorf("service.UploadFeedbackLog: 反馈日志不能为空")
	}
	if len(file) > constants.FeedbackLogMaxSize {
		return "", fmt.Errorf("service.UploadFeedbackLog: 反馈日志不能超过 2MB")
	}

	var compressed bytes.Buffer
	if err := compressFeedbackLog(file, &compressed, gzip.DefaultCompression); err != nil {
		return "", err
	}
	if s.cosClient == nil {
		return "", fmt.Errorf("service.UploadFeedbackLog: 反馈文件存储未初始化")
	}

	// 只有全部日志校验通过且 gzip 尾部写入成功后，才允许产生 COS 文件。
	url, remotePath, err := s.cosClient.GenerateFileName(cos.FeedbackLogCategory, cos.FeedbackLogFileExtension)
	if err != nil {
		return "", fmt.Errorf("service.UploadFeedbackLog: 生成反馈日志名称失败: %w", err)
	}
	if err = s.cosClient.Upload(compressed.Bytes(), remotePath); err != nil {
		return "", fmt.Errorf("service.UploadFeedbackLog: 上传反馈日志失败: %w", err)
	}
	return url, nil
}

// compressFeedbackLog 不缓存整份脱敏明文；失败时调用方必须丢弃 dst 中的部分压缩结果。
func compressFeedbackLog(file []byte, dst io.Writer, level int) error {
	writer, err := gzip.NewWriterLevel(dst, level)
	if err != nil {
		return fmt.Errorf("service.UploadFeedbackLog: 压缩反馈日志失败: %w", err)
	}
	if err = redactFeedbackLog(file, writer); err != nil {
		_ = writer.Close()
		return fmt.Errorf("service.UploadFeedbackLog: 反馈日志必须是 UTF-8 JSON Lines: %w", err)
	}
	if err = writer.Close(); err != nil {
		return fmt.Errorf("service.UploadFeedbackLog: 压缩反馈日志失败: %w", err)
	}
	return nil
}
