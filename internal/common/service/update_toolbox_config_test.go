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

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/admin_secret"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/db/toolbox"
)

func TestPutToolboxConfig(t *testing.T) {
	type testCase struct {
		name              string
		secret            string
		toolID            int64
		studentID         string
		platform          string
		version           int64
		visible           *bool
		toolName          *string
		icon              *string
		toolType          *string
		message           *string
		extra             *string
		mockSecretError   error
		mockUpsertError   error
		mockUpsertResult  *model.ToolboxConfig
		expectingError    bool
		expectingErrorMsg string
	}

	// Helper 函数
	boolPtr := func(b bool) *bool { return &b }
	stringPtr := func(s string) *string { return &s }

	testCases := []testCase{
		{
			name:            "SuccessCase",
			secret:          "valid-secret",
			toolID:          1,
			studentID:       "102301001",
			platform:        "android",
			version:         1,
			visible:         boolPtr(true),
			toolName:        stringPtr("Tool 1"),
			icon:            stringPtr("icon.png"),
			toolType:        stringPtr("type1"),
			message:         stringPtr("msg"),
			extra:           stringPtr("extra"),
			mockSecretError: nil,
			mockUpsertError: nil,
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
			expectingError: false,
		},
		{
			name:            "SuccessCaseWithMinimalFields",
			secret:          "valid-secret",
			toolID:          2,
			studentID:       "",
			platform:        "",
			version:         0,
			visible:         nil,
			toolName:        nil,
			icon:            nil,
			toolType:        nil,
			message:         nil,
			extra:           nil,
			mockSecretError: nil,
			mockUpsertError: nil,
			mockUpsertResult: &model.ToolboxConfig{
				Id:        2,
				ToolID:    2,
				StudentID: "",
				Platform:  "",
				Version:   0,
			},
			expectingError: false,
		},
		{
			name:              "InvalidSecretError",
			secret:            "invalid-secret",
			toolID:            1,
			studentID:         "",
			platform:          "",
			version:           0,
			mockSecretError:   fmt.Errorf("secret validation failed"),
			expectingError:    true,
			expectingErrorMsg: "secret validation failed",
		},
		{
			name:              "MissingToolID",
			secret:            "valid-secret",
			toolID:            0,
			studentID:         "",
			platform:          "",
			version:           0,
			mockSecretError:   nil,
			expectingError:    true,
			expectingErrorMsg: "tool_id cannot be empty",
		},
		{
			name:              "VersionTooLarge",
			secret:            "valid-secret",
			toolID:            1,
			studentID:         "",
			platform:          "",
			version:           MaxVersionNumber + 1,
			mockSecretError:   nil,
			expectingError:    true,
			expectingErrorMsg: "version cannot exceed 9,999,999",
		},
		{
			name:              "NegativeVersion",
			secret:            "valid-secret",
			toolID:            1,
			studentID:         "",
			platform:          "",
			version:           -1,
			mockSecretError:   nil,
			expectingError:    true,
			expectingErrorMsg: "version cannot be negative",
		},
		{
			name:              "UpsertError",
			secret:            "valid-secret",
			toolID:            1,
			studentID:         "",
			platform:          "",
			version:           0,
			mockSecretError:   nil,
			mockUpsertError:   fmt.Errorf("database error"),
			expectingError:    true,
			expectingErrorMsg: "upsert config failed",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient: new(db.Database),
			}

			// Mock ValidateSecret
			mockey.Mock((*admin_secret.DBAdminSecret).ValidateSecret).To(
				func(ctx context.Context, moduleName, secret string) error {
					return tc.mockSecretError
				},
			).Build()

			// Mock UpsertToolboxConfig
			mockey.Mock((*toolbox.DBToolbox).UpsertToolboxConfig).To(
				func(ctx context.Context, config *model.ToolboxConfig) (*model.ToolboxConfig, error) {
					if tc.mockUpsertError != nil {
						return nil, tc.mockUpsertError
					}
					return tc.mockUpsertResult, nil
				},
			).Build()

			commonService := NewCommonService(context.Background(), mockClientSet)
			result, err := commonService.PutToolboxConfig(
				context.Background(),
				tc.secret,
				tc.toolID,
				tc.studentID,
				tc.platform,
				tc.version,
				tc.visible,
				tc.toolName,
				tc.icon,
				tc.toolType,
				tc.message,
				tc.extra,
			)

			if tc.expectingError {
				assert.NotNil(t, err)
				if tc.expectingErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.expectingErrorMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.toolID, result.ToolID)
				if tc.visible != nil {
					assert.Equal(t, *tc.visible, result.Visible)
				}
				if tc.toolName != nil && *tc.toolName != "" {
					assert.Equal(t, *tc.toolName, result.Name)
				}
			}
		})
	}
}
