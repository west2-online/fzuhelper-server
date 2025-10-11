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
	"context"
	"fmt"
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (s *CommonService) PutToolboxConfig(ctx context.Context, secret string, toolID int64, studentID,
	platform string, version int64, visible *bool, name, icon, toolType, message, extra *string,
) (*model.ToolboxConfig, error) {
	// 验证管理员密钥
	if err := s.db.AdminSecret.ValidateSecret(ctx, "toolbox", secret); err != nil {
		return nil, err
	}

	// 验证必填参数
	if toolID == 0 {
		return nil, errno.NewErrNo(errno.ParamErrorCode, "tool_id cannot be empty")
	}

	// 验证版本号范围（如果提供了版本号）
	if version > MaxVersionNumber {
		return nil, errno.NewErrNo(errno.ParamErrorCode, "version cannot exceed 9,999,999 (7-digit limit)")
	}
	if version < 0 {
		return nil, errno.NewErrNo(errno.ParamErrorCode, "version cannot be negative")
	}

	// 构建配置对象，只设置传入的字段
	config := &model.ToolboxConfig{
		ToolID:    toolID,
		StudentID: studentID,
		Platform:  platform,
		Version:   version,
		UpdatedAt: time.Now(),
	}

	// 处理可选字段
	if visible != nil {
		config.Visible = *visible
	}
	if name != nil && *name != "" {
		config.Name = *name
	}
	if icon != nil && *icon != "" {
		config.Icon = *icon
	}
	if toolType != nil && *toolType != "" {
		config.Type = *toolType
	}
	if message != nil {
		config.Message = *message
	}
	if extra != nil && *extra != "" {
		config.Extra = *extra
	}

	result, err := s.db.Toolbox.UpsertToolboxConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("service.PutToolboxConfig: upsert config failed: %w", err)
	}

	return result, nil
}
