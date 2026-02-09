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
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/db/notice"
	"github.com/west2-online/jwch"
)

func TestGetNotice(t *testing.T) {
	type testCase struct {
		name          string
		pageNum       int
		mockDBResult  []model.Notice
		mockDBError   error
		mockJwchTotal int
		mockJwchError error
		expectList    []model.Notice
		expectTotal   int
		expectError   string
	}

	// 准备 mock 数据
	mockNotices := []model.Notice{
		{Id: 1, Title: "Notice 1", URL: "https://example.com/1", PublishedAt: "2024-01-01", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Id: 2, Title: "Notice 2", URL: "https://example.com/2", PublishedAt: "2024-01-02", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	testCases := []testCase{
		{
			name:          "SuccessCase",
			pageNum:       1,
			mockDBResult:  mockNotices,
			mockJwchTotal: 10,
			expectList:    mockNotices,
			expectTotal:   10,
		},
		{
			name:        "DBGetError",
			pageNum:     1,
			mockDBError: assert.AnError,
			expectError: "CommonService.GetNotice get notice from database",
		},
		{
			name:          "JwchGetError",
			pageNum:       1,
			mockDBResult:  mockNotices,
			mockJwchError: assert.AnError,
			expectError:   "dal.GetNoticeByPage error",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				DBClient: new(db.Database),
			}

			// Mock DB GetNoticeByPage
			mockey.Mock((*notice.DBNotice).GetNoticeByPage).Return(tc.mockDBResult, tc.mockDBError).Build()
			// Mock jwch GetNoticeInfo
			mockey.Mock((*jwch.Student).GetNoticeInfo).Return(nil, tc.mockJwchTotal, tc.mockJwchError).Build()

			commonService := NewCommonService(context.Background(), mockClientSet)
			list, total, err := commonService.GetNotice(tc.pageNum)

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
				assert.Equal(t, 0, total)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectList, list)
				assert.Equal(t, tc.expectTotal, total)
			}
		})
	}
}
