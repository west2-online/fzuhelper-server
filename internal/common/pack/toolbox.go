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

package pack

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
)

// BuildToolboxConfig 将数据库模型转换为kitex模型
func BuildToolboxConfig(dbConfig *dbmodel.ToolboxConfig) *model.ToolboxConfig {
	kitexConfig := &model.ToolboxConfig{
		ToolId: dbConfig.ToolID,
	}

	// 处理指针字段
	if dbConfig.Visible {
		kitexConfig.Visible = &dbConfig.Visible
	}

	if dbConfig.Name != "" {
		kitexConfig.Name = &dbConfig.Name
	}

	if dbConfig.Icon != "" {
		kitexConfig.Icon = &dbConfig.Icon
	}

	if dbConfig.Type != "" {
		kitexConfig.Type = &dbConfig.Type
	}

	if dbConfig.Message != "" {
		kitexConfig.Message = &dbConfig.Message
	}

	if dbConfig.Extra != "" {
		kitexConfig.Extra = &dbConfig.Extra
	}

	if dbConfig.Platform != "" {
		kitexConfig.Platform = &dbConfig.Platform
	}

	if dbConfig.Version > 0 {
		kitexConfig.Version = &dbConfig.Version
	}

	return kitexConfig
}

// BuildToolboxConfigList 将数据库模型列表转换为kitex模型列表
func BuildToolboxConfigList(dbConfigs []*dbmodel.ToolboxConfig) []*model.ToolboxConfig {
	result := make([]*model.ToolboxConfig, len(dbConfigs))
	for i, dbConfig := range dbConfigs {
		result[i] = BuildToolboxConfig(dbConfig)
	}
	return result
}
