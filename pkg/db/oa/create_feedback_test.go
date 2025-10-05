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

// filename: pkg/db/dal/feedback/feedback_test.go
package oa

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func makeSuccessFeedback() *model.Feedback {
	return &model.Feedback{
		ReportId:     1234567890123456789,
		StuId:        "102301000",
		Name:         "张三",
		College:      "计算机与大数据学院",
		ContactPhone: "13800000000",
		ContactQQ:    "123456789",
		ContactEmail: "123456789@qq.com",

		NetworkEnv:   "wifi",
		IsOnCampus:   true,
		OsName:       "Android",
		OsVersion:    "14",
		Manufacturer: "Xiaomi",
		DeviceModel:  "Mi 14",

		ProblemDesc:    "闪退",
		Screenshots:    `["https://example.com/a.png"]`,
		AppVersion:     "1.2.3",
		VersionHistory: `["1.2.0","1.2.1","1.2.3"]`,

		NetworkTraces: `[{"url":"/api/login","code":500}]`,
		Events:        `[{"event":"page_view","path":"/login"}]`,
		UserSettings:  `{"language":"zh"}`,
	}
}

func TestDBFeedback_CreateFeedback(t *testing.T) {
	type testCase struct {
		name           string
		input          *model.Feedback
		mockError      error
		expectingError bool
		ErrorMsg       string
	}

	okFeedback := makeSuccessFeedback()
	testCases := []testCase{
		{
			name:           "success",
			input:          okFeedback,
			mockError:      nil,
			expectingError: false,
		},
		{
			name:           "error",
			input:          nil,
			mockError:      gorm.ErrInvalidValue,
			expectingError: true,
			ErrorMsg:       "dal.CreateFeedback error",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			mockSnowflake := new(utils.Snowflake)
			mockDBFeedback := NewDBOA(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB { return mockGormDB }).Build()
			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB { return mockGormDB }).Build()

			mockey.Mock((*gorm.DB).Create).To(func(value interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
				} else {
					mockGormDB.Error = nil
				}
				return mockGormDB
			}).Build()

			err := mockDBFeedback.CreateFeedback(context.Background(), tc.input)
			if tc.expectingError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.ErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
