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

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

const (
	defaultToolboxConfigPageNum  = 1
	defaultToolboxConfigPageSize = 20
	maxToolboxConfigPageSize     = 100
)

func normalizeToolboxConfigListPage(pageNum, pageSize int64) (int, int, error) {
	if pageNum <= 0 {
		pageNum = defaultToolboxConfigPageNum
	}
	if pageSize <= 0 || pageSize > maxToolboxConfigPageSize {
		pageSize = defaultToolboxConfigPageSize
	}

	maxInt := int64(^uint(0) >> 1)
	if pageNum-1 > maxInt/pageSize {
		return 0, 0, errno.NewErrNo(errno.ParamErrorCode, "page offset is too large")
	}

	return int(pageNum), int(pageSize), nil
}

func (s *CommonService) GetToolboxConfigList(ctx context.Context, secret string, pageNum, pageSize int64) ([]*model.ToolboxConfig, int64, error) {
	if !utils.CheckPwd(secret) {
		return nil, 0, errno.NewErrNo(errno.AuthErrorCode, "invalid admin secret")
	}

	normalizedPageNum, normalizedPageSize, err := normalizeToolboxConfigListPage(pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}

	configs, total, err := s.db.Toolbox.ListToolboxConfigs(ctx, normalizedPageNum, normalizedPageSize)
	if err != nil {
		return nil, 0, err
	}
	if configs == nil {
		configs = make([]*model.ToolboxConfig, 0)
	}

	return configs, total, nil
}
