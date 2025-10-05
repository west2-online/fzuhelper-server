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
	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	oaDB "github.com/west2-online/fzuhelper-server/pkg/db/oa"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func makeSuccessReq() *CreateFeedbackReq {
	return &CreateFeedbackReq{
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

func makeSuccessFeedback() *dbmodel.Feedback {
	return &dbmodel.Feedback{
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

func TestOAService_CreateFeedback(t *testing.T) {
	type testCase struct {
		name              string
		req               *CreateFeedbackReq
		mockError         error
		expectingError    bool
		expectingErrorMsg string
	}

	testCases := []testCase{
		{
			name:           "success",
			req:            makeSuccessReq(),
			mockError:      nil,
			expectingError: false,
		},
		{
			name: "missing required fields",
			req: func() *CreateFeedbackReq {
				// 构造“缺少必填”的请求：让 ReportId=0
				r := makeSuccessReq()
				r.ReportId = 0
				return r
			}(),
			mockError:         nil,
			expectingError:    true,
			expectingErrorMsg: "missing required fields",
		},
		{
			name:              "dal error",
			req:               makeSuccessReq(),
			mockError:         errno.InternalServiceError,
			expectingError:    true,
			expectingErrorMsg: "service.CreateFeedback error",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// 初始化
			mockClientSet := &base.ClientSet{
				SFClient: new(utils.Snowflake),
				DBClient: new(db.Database),
			}
			oaService := NewOAService(context.Background(), "", nil, mockClientSet)

			// Mock DB 方法
			mockey.Mock((*oaDB.DBOA).CreateFeedback).To(
				func(_ *oaDB.DBOA, _ context.Context, _ *dbmodel.Feedback) error {
					return tc.mockError
				},
			).Build()

			// 开始测试
			err := oaService.CreateFeedback(tc.req)

			if tc.expectingError {
				assert.Error(t, err, "error should not be nil")
				if tc.expectingErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.expectingErrorMsg)
				}
			} else {
				assert.NoError(t, err, "should be no error")
			}
		})
	}
}

func TestOAService_GetFeedback(t *testing.T) {
	type testCase struct {
		name              string
		id                int64
		mockError         error
		mockFb            *dbmodel.Feedback
		expectingError    bool
		expectingErrorMsg string
		expectedInfo      *dbmodel.Feedback
	}

	testCases := []testCase{
		{
			name:           "success",
			id:             1234567890123456789,
			mockError:      nil,
			mockFb:         makeSuccessFeedback(),
			expectingError: false,
			expectedInfo:   makeSuccessFeedback(),
		},
		{
			name:              "invalid id",
			id:                0,
			mockError:         nil,
			expectingError:    true,
			expectingErrorMsg: "invalid id",
			expectedInfo:      nil,
		},
		{
			name:              "not ok",
			id:                1,
			mockError:         errno.InternalServiceError,
			expectingError:    true,
			expectingErrorMsg: "service.GetFeedback error",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// 初始化
			mockClientSet := &base.ClientSet{
				SFClient: new(utils.Snowflake),
				DBClient: new(db.Database),
			}
			oaService := NewOAService(context.Background(), "", nil, mockClientSet)

			// Mock DAL：GetFeedbackById
			mockey.Mock((*oaDB.DBOA).GetFeedbackById).To(
				func(_ *oaDB.DBOA, _ context.Context, id int64) (bool, *dbmodel.Feedback, error) {
					if tc.mockError == nil {
						return true, tc.mockFb, nil
					}
					return false, nil, tc.mockError
				},
			).Build()

			// 执行
			feedback, err := oaService.GetFeedback(tc.id)

			// 断言
			if tc.expectingError {
				assert.Error(t, err, "error should not be nil")
				if tc.expectingErrorMsg != "" {
					assert.Contains(t, err.Error(), tc.expectingErrorMsg)
				}
				assert.Nil(t, feedback)
			} else {
				assert.NoError(t, err, "should be no error")
				assert.Equal(t, tc.expectedInfo, feedback)
				assert.Equal(t, tc.id, feedback.ReportId)
			}
		})
	}
}
