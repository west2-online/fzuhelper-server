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
	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/db/version"
)

func TestGetDump(t *testing.T) {
	type testCase struct {
		name            string         // 测试用例名称
		mockReturn      []*model.Visit // mock返回的访问数据列表
		mockError       error          // mock返回的错误
		mockMarshalErr  error          // mock序列化时返回的错误
		expectResult    string         // 期望的结果字符串
		expectError     string         // 期望的错误信息
		needMarshalMock bool           // 是否需要mock序列化错误
	}

	testCases := []testCase{
		{
			name: "success case",
			mockReturn: []*model.Visit{
				{Id: 1, Date: "2025-01-01", Visits: 100},
				{Id: 2, Date: "2025-01-02", Visits: 200},
			},
			mockError:       nil,
			mockMarshalErr:  nil,
			expectResult:    `{"2025-01-01":100,"2025-01-02":200}`,
			expectError:     "",
			needMarshalMock: false,
		},
		{
			name:            "error case",
			mockReturn:      nil,
			mockError:       fmt.Errorf("database error"),
			mockMarshalErr:  nil,
			expectResult:    "",
			expectError:     "GetDump: get version list error: database error",
			needMarshalMock: false,
		},
		{
			name: "marshal error case",
			mockReturn: []*model.Visit{
				{Id: 1, Date: "2025-01-01", Visits: 100},
			},
			mockError:       nil,
			mockMarshalErr:  fmt.Errorf("marshal failed"),
			expectResult:    "",
			expectError:     "GetDump: marshal error:",
			needMarshalMock: true,
		},
	}

	defer mockey.UnPatchAll() // 清理所有mock

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient: new(db.Database),
			}

			mockey.Mock((*version.DBVersion).GetVersionList).Return(tc.mockReturn, tc.mockError).Build()
			if tc.needMarshalMock {
				mockey.Mock(sonic.Marshal).Return(nil, tc.mockMarshalErr).Build()
			}

			versionService := NewVersionService(context.Background(), mockClientSet)
			result, err := versionService.GetDump()
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, tc.expectResult, result)
			}
		})
	}
}
