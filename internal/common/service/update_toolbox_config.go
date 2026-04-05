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
	"time"

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *CommonService) PutToolboxConfig(ctx context.Context, req *common.PutToolboxConfigRequest) (*model.ToolboxConfig, error) {
	// 获取请求参数，如果为空则使用默认值
	studentID := ""
	if req.StudentId != nil {
		studentID = *req.StudentId
	}

	platform := ""
	if req.Platform != nil {
		platform = *req.Platform
	}

	version := int64(0)
	if req.Version != nil {
		version = *req.Version
	}

	// 验证管理员密钥
	if !utils.CheckPwd(req.Secret) {
		return nil, errno.AuthError.WithMessage("Common.PutToolboxConfig: invalid admin secret")
	}

	// 验证必填参数
	if req.ToolId == 0 {
		return nil, errno.ParamError.WithMessage("Common.PutToolboxConfig: tool_id cannot be empty")
	}

	// 验证版本号范围（如果提供了版本号）
	if version > MaxVersionNumber {
		return nil, errno.ParamError.WithMessage("Common.PutToolboxConfig: version cannot exceed 9,999,999 (7-digit limit)")
	}
	if version < 0 {
		return nil, errno.ParamError.WithMessage("Common.PutToolboxConfig: version cannot be negative")
	}

	// 构建配置对象，只设置传入的字段
	config := &model.ToolboxConfig{
		ToolID:    req.ToolId,
		StudentID: studentID,
		Platform:  platform,
		Version:   version,
		UpdatedAt: time.Now(),
	}

	// 处理可选字段
	if req.Visible != nil {
		config.Visible = *req.Visible
	}
	if req.Name != nil && *req.Name != "" {
		config.Name = *req.Name
	}
	if req.Icon != nil && *req.Icon != "" {
		config.Icon = *req.Icon
	}
	if req.Type != nil && *req.Type != "" {
		config.Type = *req.Type
	}
	if req.Message != nil {
		config.Message = *req.Message
	}
	if req.Extra != nil && *req.Extra != "" {
		config.Extra = *req.Extra
	}

	result, err := s.db.Toolbox.UpsertToolboxConfig(ctx, config)
	if err != nil {
		return nil, errno.ErrNoWithPreMessage(err, "Common.PutToolboxConfig: Upsert config failed")
	}

	return result, nil
}
