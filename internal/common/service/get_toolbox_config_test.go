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
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/db/toolbox"
)

func TestGetMatchScore(t *testing.T) {
	type testCase struct {
		name        string
		config      *model.ToolboxConfig
		studentID   string
		platform    string
		version     int64
		expectScore int
	}

	testCases := []testCase{
		{
			name: "AllMatch",
			config: &model.ToolboxConfig{
				StudentID: "102301001",
				Platform:  "android",
				Version:   1,
			},
			studentID:   "102301001",
			platform:    "android",
			version:     2,
			expectScore: (1 << 25) | (1 << 1) | 1, // 学号匹配 + 版本匹配 + 平台匹配
		},
		{
			name: "StudentIDMatchOnly",
			config: &model.ToolboxConfig{
				StudentID: "102301001",
				Platform:  "ios",
			},
			studentID:   "102301001",
			platform:    "android",
			version:     1,
			expectScore: -1, // 平台不匹配
		},
		{
			name: "VersionAndPlatformMatch",
			config: &model.ToolboxConfig{
				Platform: "android",
				Version:  1,
			},
			studentID:   "102301002",
			platform:    "android",
			version:     2,
			expectScore: (1 << 1) | 1, // 版本匹配 + 平台匹配
		},
		{
			name: "OnlyVersionMatch",
			config: &model.ToolboxConfig{
				Platform: "ios",
				Version:  1,
			},
			platform:    "android",
			version:     2,
			expectScore: -1, // 平台不匹配
		},
		{
			name:        "NoConstraints",
			config:      &model.ToolboxConfig{},
			platform:    "android",
			version:     2,
			expectScore: 0, // 无任何限制的配置
		},
		{
			name: "VersionNotSatisfied",
			config: &model.ToolboxConfig{
				Platform: "android",
				Version:  3,
			},
			platform:    "android",
			version:     2,
			expectScore: -1, // 版本要求高于客户端版本
		},
		{
			name: "VersionExceedsMaxLimit",
			config: &model.ToolboxConfig{
				Platform: "android",
				Version:  99999999, // 超过 MaxVersionNumber (9999999)
			},
			platform:    "android",
			version:     100000000,
			expectScore: (MaxVersionNumber << 1) | 1, // 被限制到 MaxVersionNumber，然后计算匹配分数
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			score := getMatchScore(tc.config, tc.studentID, tc.platform, tc.version)
			assert.Equal(t, tc.expectScore, score)
		})
	}
}

func TestGetToolboxConfig(t *testing.T) {
	type testCase struct {
		name         string
		studentID    string
		platform     string
		version      int64
		mockDBResult []*model.ToolboxConfig
		mockDBError  error
		expectResult []*model.ToolboxConfig
		expectError  string
	}

	// 辅助函数：创建工具配置
	newToolConfig := func(id, toolID int64, studentID, platform, name string, version int64) *model.ToolboxConfig {
		return &model.ToolboxConfig{
			Id:        id,
			ToolID:    toolID,
			Visible:   true,
			Name:      name,
			Icon:      fmt.Sprintf("icon%d.png", toolID),
			Type:      fmt.Sprintf("type%d", toolID),
			StudentID: studentID,
			Platform:  platform,
			Version:   version,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	// 测试数据集
	singleTool := []*model.ToolboxConfig{
		newToolConfig(1, 1, "102301001", "android", "Tool 1", 1),
	}

	multiTools := []*model.ToolboxConfig{
		newToolConfig(1, 1, "102301001", "android", "Tool 1 - StudentID", 1), // 学号匹配
		newToolConfig(2, 1, "", "android", "Tool 1 - Platform", 1),           // 只平台匹配
		newToolConfig(3, 2, "", "android", "Tool 2 - Version 1", 1),          // 版本1
		newToolConfig(4, 2, "", "android", "Tool 2 - Version 2", 2),          // 版本2
	}

	testCases := []testCase{
		{
			name:         "SuccessCaseWithOneTool",
			studentID:    "102301001",
			platform:     "android",
			version:      1,
			mockDBResult: singleTool,
			expectResult: singleTool,
		},
		{
			name:        "DBGetError",
			studentID:   "102301001",
			platform:    "android",
			version:     1,
			mockDBError: assert.AnError,
			expectError: "assert.AnError",
		},
		{
			name:         "MultipleToolsSelectBestMatch",
			studentID:    "102301001",
			platform:     "android",
			version:      2,
			mockDBResult: multiTools,
			expectResult: []*model.ToolboxConfig{multiTools[0], multiTools[3]},
		},
		{
			name:         "NoConstraintsConfig",
			mockDBResult: singleTool,
			expectResult: singleTool,
		},
		{
			name:         "FilterNotMatchingConfigs",
			studentID:    "102301002",
			platform:     "ios",
			version:      2,
			mockDBResult: multiTools,
			expectResult: []*model.ToolboxConfig{},
		},
		{
			name:         "EmptyConfigResult",
			studentID:    "102301001",
			platform:     "android",
			version:      1,
			mockDBResult: []*model.ToolboxConfig{},
			expectResult: []*model.ToolboxConfig{},
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient: new(db.Database),
			}

			// Mock DB GetToolboxConfigs
			mockey.Mock((*toolbox.DBToolbox).GetToolboxConfigs).Return(tc.mockDBResult, tc.mockDBError).Build()

			commonService := NewCommonService(context.Background(), mockClientSet)
			result, err := commonService.GetToolboxConfig(context.Background(), tc.studentID, tc.platform, tc.version)

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.Nil(t, err)
				assert.ElementsMatch(t, tc.expectResult, result)
			}
		})
	}
}
