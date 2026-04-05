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

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/db/toolbox"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestPutToolboxConfig(t *testing.T) {
	type testCase struct {
		name             string
		req              *common.PutToolboxConfigRequest
		mockCheckPwd     bool
		mockUpsertError  error
		mockUpsertResult *model.ToolboxConfig
		expectError      string
	}

	// Helper 函数
	boolPtr := func(b bool) *bool { return &b }
	stringPtr := func(s string) *string { return &s }
	int64Ptr := func(i int64) *int64 { return &i }

	testCases := []testCase{
		{
			name:         "SuccessCase",
			mockCheckPwd: true,
			req: &common.PutToolboxConfigRequest{
				Secret:    "valid_secret",
				ToolId:    1,
				StudentId: stringPtr("102301001"),
				Platform:  stringPtr("android"),
				Version:   int64Ptr(1),
				Visible:   boolPtr(true),
				Name:      stringPtr("Tool 1"),
				Icon:      stringPtr("icon.png"),
				Type:      stringPtr("type1"),
				Message:   stringPtr("msg"),
				Extra:     stringPtr("extra"),
			},
			mockUpsertResult: &model.ToolboxConfig{
				Id:        1,
				ToolID:    1,
				Visible:   true,
				Name:      "Tool 1",
				Icon:      "icon.png",
				Type:      "type1",
				Message:   "msg",
				Extra:     "extra",
				StudentID: "102301001",
				Platform:  "android",
				Version:   1,
			},
		},
		{
			name:         "SuccessCaseWithMinimalFields",
			mockCheckPwd: true,
			req: &common.PutToolboxConfigRequest{
				Secret: "valid_secret",
				ToolId: 2,
			},
			mockUpsertResult: &model.ToolboxConfig{
				Id:     2,
				ToolID: 2,
			},
		},
		{
			name:         "InvalidSecretError",
			mockCheckPwd: false,
			req: &common.PutToolboxConfigRequest{
				Secret: "invalid_secret",
				ToolId: 1,
			},
			expectError: "invalid admin secret",
		},
		{
			name:         "MissingToolID",
			mockCheckPwd: true,
			req: &common.PutToolboxConfigRequest{
				Secret: "valid_secret",
			},
			expectError: "tool_id cannot be empty",
		},
		{
			name:         "VersionTooLarge",
			mockCheckPwd: true,
			req: &common.PutToolboxConfigRequest{
				Secret:  "valid_secret",
				ToolId:  1,
				Version: int64Ptr(MaxVersionNumber + 1),
			},
			expectError: "version cannot exceed 9,999,999",
		},
		{
			name:         "NegativeVersion",
			mockCheckPwd: true,
			req: &common.PutToolboxConfigRequest{
				Secret:  "valid_secret",
				ToolId:  1,
				Version: int64Ptr(-1),
			},
			expectError: "version cannot be negative",
		},
		{
			name:         "UpsertError",
			mockCheckPwd: true,
			req: &common.PutToolboxConfigRequest{
				Secret: "valid_secret",
				ToolId: 1,
			},
			mockUpsertError: assert.AnError,
			expectError:     "Upsert config failed",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient: new(db.Database),
			}

			mockey.Mock(utils.CheckPwd).Return(tc.mockCheckPwd).Build()

			// Mock UpsertToolboxConfig
			mockey.Mock((*toolbox.DBToolbox).UpsertToolboxConfig).Return(tc.mockUpsertResult, tc.mockUpsertError).Build()

			commonService := NewCommonService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			result, err := commonService.PutToolboxConfig(context.Background(), tc.req)

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.req.ToolId, result.ToolID)
				if tc.req.Visible != nil {
					assert.Equal(t, *tc.req.Visible, result.Visible)
				}
				if tc.req.Name != nil && *tc.req.Name != "" {
					assert.Equal(t, *tc.req.Name, result.Name)
				}
			}
		})
	}
}
