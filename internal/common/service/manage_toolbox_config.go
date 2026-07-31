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

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func validateToolboxAdminSecret(secret string) error {
	if !utils.CheckPwd(secret) {
		return errno.NewErrNo(errno.AuthErrorCode, "invalid admin secret")
	}
	return nil
}

func validateToolboxConfig(config *model.ToolboxConfig) error {
	if config == nil {
		return errno.NewErrNo(errno.ParamErrorCode, "toolbox config cannot be nil")
	}
	if config.ToolID <= 0 {
		return errno.NewErrNo(errno.ParamErrorCode, "tool_id must be positive")
	}
	if config.Version > MaxVersionNumber {
		return errno.NewErrNo(errno.ParamErrorCode, "version cannot exceed 9,999,999 (7-digit limit)")
	}
	if config.Version < 0 {
		return errno.NewErrNo(errno.ParamErrorCode, "version cannot be negative")
	}
	return nil
}

func validateToolboxConfigID(id int64) error {
	if id <= 0 {
		return errno.NewErrNo(errno.ParamErrorCode, "config_id must be positive")
	}
	return nil
}

func (s *CommonService) CreateToolboxConfig(
	ctx context.Context,
	secret string,
	config *model.ToolboxConfig,
) (*model.ToolboxConfig, error) {
	if err := validateToolboxAdminSecret(secret); err != nil {
		return nil, err
	}
	if err := validateToolboxConfig(config); err != nil {
		return nil, err
	}
	if err := s.db.Toolbox.CreateToolboxConfig(ctx, config); err != nil {
		return nil, fmt.Errorf("service.CreateToolboxConfig: %w", err)
	}
	return config, nil
}

func (s *CommonService) GetToolboxConfigByID(
	ctx context.Context,
	secret string,
	id int64,
) (*model.ToolboxConfig, error) {
	if err := validateToolboxAdminSecret(secret); err != nil {
		return nil, err
	}
	if err := validateToolboxConfigID(id); err != nil {
		return nil, err
	}
	config, err := s.db.Toolbox.GetToolboxConfigByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.GetToolboxConfigByID: %w", err)
	}
	return config, nil
}

func (s *CommonService) UpdateToolboxConfig(
	ctx context.Context,
	secret string,
	id int64,
	config *model.ToolboxConfig,
) (*model.ToolboxConfig, error) {
	if err := validateToolboxAdminSecret(secret); err != nil {
		return nil, err
	}
	if err := validateToolboxConfigID(id); err != nil {
		return nil, err
	}
	if err := validateToolboxConfig(config); err != nil {
		return nil, err
	}
	updated, err := s.db.Toolbox.UpdateToolboxConfig(ctx, id, config)
	if err != nil {
		return nil, fmt.Errorf("service.UpdateToolboxConfig: %w", err)
	}
	return updated, nil
}

func (s *CommonService) DeleteToolboxConfig(ctx context.Context, secret string, id int64) error {
	if err := validateToolboxAdminSecret(secret); err != nil {
		return err
	}
	if err := validateToolboxConfigID(id); err != nil {
		return err
	}
	if err := s.db.Toolbox.DeleteToolboxConfig(ctx, id); err != nil {
		return fmt.Errorf("service.DeleteToolboxConfig: %w", err)
	}
	return nil
}
