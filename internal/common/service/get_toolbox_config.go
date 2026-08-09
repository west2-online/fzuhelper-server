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
	"sort"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

const (
	// MaxVersionNumber 版本号最大值，支持7位数字
	MaxVersionNumber = 9999999
	// StudentIDMatchBit 学号匹配位位置（第25位）
	StudentIDMatchBit = 25
)

// MatchResult 匹配结果
type MatchResult struct {
	Config     *model.ToolboxConfig
	MatchScore int // 匹配分数，分数越高优先级越高
}

// getMatchScore 计算配置的匹配分数，分数越高优先级越高
// 优先级：学号 > 版本 > 平台，匹配的维度越多分数越高
func getMatchScore(config *model.ToolboxConfig, studentID string, platform string, version int64) int {
	// 使用位运算来精确表示匹配情况，确保优先级绝对正确
	// 格式：学号匹配(1位) + 版本匹配程度(24位) + 平台匹配(1位)
	// 版本号最大支持7位数字(9,999,999)，需要24位二进制表示

	var matchBits int64 = 0

	// 学号匹配检查（最高优先级，占第25位）
	if config.StudentID != nil && *config.StudentID != "" {
		if studentID == "" || *config.StudentID != studentID {
			return -1
		}
		matchBits |= (1 << StudentIDMatchBit) // 设置第25位，表示学号匹配
	}

	// 版本匹配检查（第二优先级，占第1-24位）
	if config.Version != nil && *config.Version > 0 {
		if version <= 0 || *config.Version > version {
			return -1
		}
		// 版本匹配，版本号越高分数越高
		// 占第1-24位，最多支持版本号到9,999,999（7位数字）
		versionScore := *config.Version
		if versionScore > MaxVersionNumber {
			versionScore = MaxVersionNumber // 限制最大值为7位数字
		}
		matchBits |= (versionScore << 1)
	}

	// 平台匹配检查（最低优先级，占第0位）
	if config.Platform != nil && *config.Platform != "" {
		if platform == "" || *config.Platform != platform {
			return -1
		}
		matchBits |= 1 // 设置第0位，表示平台匹配
	}

	return int(matchBits)
}

func (s *CommonService) GetToolboxConfig(ctx context.Context, studentID string, platform string, version int64) ([]*model.ToolboxConfig, error) {
	// 获取数据库中所有的工具箱配置
	allConfigs, err := s.db.Toolbox.GetToolboxConfigs(ctx)
	if err != nil {
		return nil, err
	}

	// 按ToolID分组，每个工具找到最高匹配分数的配置
	toolBestMatch := make(map[int64]*MatchResult)

	for _, config := range allConfigs {
		matchScore := getMatchScore(config, studentID, platform, version)

		// 跳过不匹配的配置
		if matchScore < 0 {
			continue
		}
		toolID := config.ToolID
		currentBest, exists := toolBestMatch[toolID]

		if !exists || matchScore > currentBest.MatchScore ||
			(matchScore == currentBest.MatchScore && config.Id > currentBest.Config.Id) {
			// 如果是新工具或找到更高匹配分数的配置，则更新
			toolBestMatch[toolID] = &MatchResult{
				Config:     config,
				MatchScore: matchScore,
			}
		}
	}

	// 按 tool_id 排序，确保响应顺序稳定。
	toolIDs := make([]int64, 0, len(toolBestMatch))
	for toolID := range toolBestMatch {
		toolIDs = append(toolIDs, toolID)
	}
	sort.Slice(toolIDs, func(i, j int) bool {
		return toolIDs[i] < toolIDs[j]
	})

	result := make([]*model.ToolboxConfig, 0, len(toolIDs))
	for _, toolID := range toolIDs {
		result = append(result, toolBestMatch[toolID].Config)
	}

	return result, nil
}
