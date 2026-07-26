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
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/db/toolbox"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestGetToolboxConfigList(t *testing.T) {
	type testCase struct {
		name           string
		secret         string
		pageNum        int64
		pageSize       int64
		mockCheckPwd   bool
		mockDBResult   []*model.ToolboxConfig
		mockDBTotal    int64
		mockDBError    error
		expectPageNum  int
		expectPageSize int
		expectError    string
	}

	configs := []*model.ToolboxConfig{
		{Id: 2, ToolID: 1, StudentID: "102300217", Platform: "android", Version: 2},
		{Id: 1, ToolID: 1, Platform: "ios", Version: 1},
	}

	testCases := []testCase{
		{
			name:           "success",
			secret:         "secret",
			pageNum:        2,
			pageSize:       2,
			mockCheckPwd:   true,
			mockDBResult:   configs,
			mockDBTotal:    3,
			expectPageNum:  2,
			expectPageSize: 2,
		},
		{
			name:         "invalid_secret",
			secret:       "wrong",
			pageNum:      1,
			pageSize:     20,
			mockCheckPwd: false,
			expectError:  "invalid admin secret",
		},
		{
			name:           "default_page",
			secret:         "secret",
			pageNum:        0,
			pageSize:       101,
			mockCheckPwd:   true,
			mockDBResult:   []*model.ToolboxConfig{},
			mockDBTotal:    0,
			expectPageNum:  defaultToolboxConfigPageNum,
			expectPageSize: defaultToolboxConfigPageSize,
		},
		{
			name:           "nil_result_to_empty_slice",
			secret:         "secret",
			pageNum:        1,
			pageSize:       20,
			mockCheckPwd:   true,
			mockDBResult:   nil,
			mockDBTotal:    0,
			expectPageNum:  1,
			expectPageSize: 20,
		},
		{
			name:           "db_error",
			secret:         "secret",
			pageNum:        1,
			pageSize:       20,
			mockCheckPwd:   true,
			mockDBError:    assert.AnError,
			expectPageNum:  1,
			expectPageSize: 20,
			expectError:    "assert.AnError",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient: new(db.Database),
			}

			mockey.Mock(utils.CheckPwd).Return(tc.mockCheckPwd).Build()
			mockey.Mock((*toolbox.DBToolbox).ListToolboxConfigs).To(
				func(ctx context.Context, pageNum, pageSize int) ([]*model.ToolboxConfig, int64, error) {
					assert.Equal(t, tc.expectPageNum, pageNum)
					assert.Equal(t, tc.expectPageSize, pageSize)
					return tc.mockDBResult, tc.mockDBTotal, tc.mockDBError
				},
			).Build()

			commonService := NewCommonService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			result, total, err := commonService.GetToolboxConfigList(context.Background(), tc.secret, tc.pageNum, tc.pageSize)

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, result)
				assert.Zero(t, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.mockDBTotal, total)
				if tc.mockDBResult == nil {
					assert.Empty(t, result)
				} else {
					assert.Equal(t, tc.mockDBResult, result)
				}
			}
		})
	}
}
