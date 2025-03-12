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
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/db/version"
)

func TestGetDump(t *testing.T) {
	testCases := []struct {
		name           string
		mockReturn     []*model.Visit
		mockError      error
		expectedResult string
		expectedErr    string
	}{
		{
			name: "success case",
			mockReturn: []*model.Visit{
				{Id: 1, Date: "2025-01-01", Visits: 100},
				{Id: 2, Date: "2025-01-02", Visits: 200},
			},
			mockError:      nil,
			expectedResult: `{"2025-01-01":100,"2025-01-02":200}`,
			expectedErr:    "",
		},
		{
			name:           "error case",
			mockReturn:     nil,
			mockError:      fmt.Errorf("database error"),
			expectedResult: "",
			expectedErr:    "GetDump: get version list error: database error",
		},
	}
	defer mockey.UnPatchAll() // 清理所有mock

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*version.DBVersion).GetVersionList).Return(tc.mockReturn, tc.mockError).Build()
			mockClientSet := new(base.ClientSet)
			mockClientSet.DBClient = new(db.Database)
			versionService := NewVersionService(context.Background(), mockClientSet)
			result, err := versionService.GetDump()

			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
