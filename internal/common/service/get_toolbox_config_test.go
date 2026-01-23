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
		name          string
		config        *model.ToolboxConfig
		studentID     string
		platform      string
		version       int64
		expectedScore int
	}

	testCases := []testCase{
		{
			name: "AllMatch",
			config: &model.ToolboxConfig{
				StudentID: "102301001",
				Platform:  "android",
				Version:   1,
			},
			studentID:     "102301001",
			platform:      "android",
			version:       2,
			expectedScore: (1 << 25) | (1 << 1) | 1, // 学号匹配 + 版本匹配 + 平台匹配
		},
		{
			name: "StudentIDMatchOnly",
			config: &model.ToolboxConfig{
				StudentID: "102301001",
				Platform:  "ios",
				Version:   0,
			},
			studentID:     "102301001",
			platform:      "android",
			version:       1,
			expectedScore: -1, // 平台不匹配
		},
		{
			name: "VersionAndPlatformMatch",
			config: &model.ToolboxConfig{
				StudentID: "",
				Platform:  "android",
				Version:   1,
			},
			studentID:     "102301002",
			platform:      "android",
			version:       2,
			expectedScore: (1 << 1) | 1, // 版本匹配 + 平台匹配
		},
		{
			name: "OnlyVersionMatch",
			config: &model.ToolboxConfig{
				StudentID: "",
				Platform:  "ios",
				Version:   1,
			},
			studentID:     "",
			platform:      "android",
			version:       2,
			expectedScore: -1, // 平台不匹配
		},
		{
			name: "NoConstraints",
			config: &model.ToolboxConfig{
				StudentID: "",
				Platform:  "",
				Version:   0,
			},
			studentID:     "",
			platform:      "android",
			version:       2,
			expectedScore: 0, // 无任何限制的配置
		},
		{
			name: "VersionNotSatisfied",
			config: &model.ToolboxConfig{
				StudentID: "",
				Platform:  "android",
				Version:   3,
			},
			studentID:     "",
			platform:      "android",
			version:       2,
			expectedScore: -1, // 版本要求高于客户端版本
		},
		{
			name: "VersionExceedsMaxLimit",
			config: &model.ToolboxConfig{
				StudentID: "",
				Platform:  "android",
				Version:   99999999, // 超过 MaxVersionNumber (9999999)
			},
			studentID:     "",
			platform:      "android",
			version:       100000000,
			expectedScore: (MaxVersionNumber << 1) | 1, // 被限制到 MaxVersionNumber，然后计算匹配分数
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			score := getMatchScore(tc.config, tc.studentID, tc.platform, tc.version)
			assert.Equal(t, tc.expectedScore, score)
		})
	}
}

func TestGetToolboxConfig(t *testing.T) {
	type testCase struct {
		name              string
		studentID         string
		platform          string
		version           int64
		mockDBResult      []*model.ToolboxConfig
		mockDBError       error
		expectedResult    []*model.ToolboxConfig
		expectingError    bool
		expectingErrorMsg string
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
			name:           "SuccessCaseWithOneTool",
			studentID:      "102301001",
			platform:       "android",
			version:        1,
			mockDBResult:   singleTool,
			mockDBError:    nil,
			expectedResult: singleTool,
			expectingError: false,
		},
		{
			name:              "DBGetError",
			studentID:         "102301001",
			platform:          "android",
			version:           1,
			mockDBResult:      nil,
			mockDBError:       fmt.Errorf("database error"),
			expectedResult:    nil,
			expectingError:    true,
			expectingErrorMsg: "database error",
		},
		{
			name:           "MultipleToolsSelectBestMatch",
			studentID:      "102301001",
			platform:       "android",
			version:        2,
			mockDBResult:   multiTools,
			mockDBError:    nil,
			expectedResult: []*model.ToolboxConfig{multiTools[0], multiTools[3]},
			expectingError: false,
		},
		{
			name:           "NoConstraintsConfig",
			studentID:      "",
			platform:       "",
			version:        0,
			mockDBResult:   singleTool,
			mockDBError:    nil,
			expectedResult: singleTool,
			expectingError: false,
		},
		{
			name:           "FilterNotMatchingConfigs",
			studentID:      "102301002",
			platform:       "ios",
			version:        2,
			mockDBResult:   multiTools,
			mockDBError:    nil,
			expectedResult: []*model.ToolboxConfig{},
			expectingError: false,
		},
		{
			name:           "EmptyConfigResult",
			studentID:      "102301001",
			platform:       "android",
			version:        1,
			mockDBResult:   []*model.ToolboxConfig{},
			mockDBError:    nil,
			expectedResult: []*model.ToolboxConfig{},
			expectingError: false,
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient: new(db.Database),
			}

			// Mock DB GetToolboxConfigs
			mockey.Mock((*toolbox.DBToolbox).GetToolboxConfigs).To(
				func(ctx context.Context) ([]*model.ToolboxConfig, error) {
					if tc.mockDBError != nil {
						return nil, tc.mockDBError
					}
					return tc.mockDBResult, nil
				},
			).Build()

			commonService := NewCommonService(context.Background(), mockClientSet)
			result, err := commonService.GetToolboxConfig(context.Background(), tc.studentID, tc.platform, tc.version)

			if tc.expectingError {
				assert.NotNil(t, err)
				if tc.expectingErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.expectingErrorMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				// 验证返回结果的工具数量
				assert.Len(t, result, len(tc.expectedResult))
				// 验证返回的配置内容
				for _, expectedConfig := range tc.expectedResult {
					found := false
					for _, resultConfig := range result {
						if resultConfig.Id == expectedConfig.Id {
							assert.Equal(t, expectedConfig, resultConfig)
							found = true
							break
						}
					}
					assert.True(t, found, "expected config with id %d not found in result", expectedConfig.Id)
				}
			}
		})
	}
}
