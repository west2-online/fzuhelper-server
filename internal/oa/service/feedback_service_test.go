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
	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	oaDB "github.com/west2-online/fzuhelper-server/pkg/db/oa"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func makeSuccessReq() *CreateFeedbackReq {
	return &CreateFeedbackReq{
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

func makeSuccessFeedbackListReq() *FeedbackListReq {
	return &FeedbackListReq{
		Name: "张三",
	}
}

func makeSuccessFeedbackList() []dbmodel.FeedbackListItem {
	return []dbmodel.FeedbackListItem{
		{
			ReportId:    2199023256001,
			Name:        "张三",
			NetworkEnv:  dbmodel.NetworkWifi,
			ProblemDesc: "登录后提示网络异常，刷新失败",
			AppVersion:  "2.3.1",
		},
	}
}

func TestCreateFeedback(t *testing.T) {
	type testCase struct {
		name        string
		req         *CreateFeedbackReq
		mockError   error
		mockSFError error
		expectError string
	}

	testCases := []testCase{
		{
			name: "success",
			req:  makeSuccessReq(),
		},
		{
			name: "missing required fields",
			req: func() *CreateFeedbackReq {
				// 构造“缺少必填”的请求：让 ReportId=0
				r := makeSuccessReq()
				r.Name = ""
				return r
			}(),
			mockError:   nil,
			expectError: "missing required fields",
		},
		{
			name:        "dal error",
			req:         makeSuccessReq(),
			mockError:   errno.InternalServiceError,
			expectError: "service.CreateFeedback error",
		},
		{
			name: "invalid NetworkEnv corrected",
			req: func() *CreateFeedbackReq {
				r := makeSuccessReq()
				r.NetworkEnv = "invalid_network"
				return r
			}(),
			mockError: nil,
		},
		{
			name:        "Snowflake NextVal error",
			req:         makeSuccessReq(),
			mockError:   nil,
			mockSFError: errno.InternalServiceError,
			expectError: "generate report_id failed",
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

			// Mock Snowflake
			if tc.mockSFError != nil {
				mockey.Mock((*utils.Snowflake).NextVal).Return(int64(0), tc.mockSFError).Build()
			} else {
				mockey.Mock((*utils.Snowflake).NextVal).Return(int64(1234567890123456789), nil).Build()
			}

			// Mock DB 方法
			mockey.Mock((*oaDB.DBOA).CreateFeedback).Return(tc.mockError).Build()

			// 开始测试
			_, err := oaService.CreateFeedback(tc.req)
			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err, "should be no error")
			}
		})
	}
}

func TestGetFeedbackById(t *testing.T) {
	type testCase struct {
		name        string
		id          int64
		mockOk      bool
		mockError   error
		mockFb      *dbmodel.Feedback
		expectError string
		expectInfo  *dbmodel.Feedback
	}

	testCases := []testCase{
		{
			name:       "success",
			id:         1234567890123456789,
			mockOk:     true,
			mockError:  nil,
			mockFb:     makeSuccessFeedback(),
			expectInfo: makeSuccessFeedback(),
		},
		{
			name:        "invalid id",
			id:          0,
			mockOk:      true,
			mockError:   nil,
			expectError: "invalid id",
			expectInfo:  nil,
		},
		{
			name:        "record not found",
			id:          1234567890123456789,
			mockOk:      false,
			mockError:   nil,
			mockFb:      nil,
			expectError: "service.GetFeedback error",
		},
		{
			name:        "database error",
			id:          1,
			mockOk:      false,
			mockError:   errno.InternalServiceError,
			expectError: "service.GetFeedback error",
		},
		{
			name:      "invalid NetworkEnv corrected",
			id:        1234567890123456789,
			mockOk:    true,
			mockError: nil,
			mockFb: func() *dbmodel.Feedback {
				fb := makeSuccessFeedback()
				fb.NetworkEnv = "invalid"
				return fb
			}(),
		},
		{
			name:      "empty JSON fields corrected",
			id:        1234567890123456789,
			mockOk:    true,
			mockError: nil,
			mockFb: func() *dbmodel.Feedback {
				fb := makeSuccessFeedback()
				fb.Screenshots = ""
				fb.VersionHistory = ""
				fb.NetworkTraces = ""
				fb.Events = ""
				fb.UserSettings = ""
				return fb
			}(),
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
			mockey.Mock((*oaDB.DBOA).GetFeedbackById).Return(tc.mockOk, tc.mockFb, tc.mockError).Build()

			// 执行
			feedback, err := oaService.GetFeedbackById(tc.id)
			// 断言
			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, feedback)
			} else {
				assert.NoError(t, err, "should be no error")
				assert.NotNil(t, feedback)
				assert.Equal(t, tc.id, feedback.ReportId)
			}
		})
	}
}

func TestGetFeedbackList(t *testing.T) {
	type testCase struct {
		name        string
		req         *FeedbackListReq
		mockError   error
		mockFb      []dbmodel.FeedbackListItem
		expectError string
		expectInfo  []dbmodel.FeedbackListItem
	}

	testCases := []testCase{
		{
			name:       "success",
			req:        makeSuccessFeedbackListReq(),
			mockError:  nil,
			mockFb:     makeSuccessFeedbackList(),
			expectInfo: makeSuccessFeedbackList(),
		},
		{
			name:        "request is nil",
			req:         nil,
			mockError:   nil,
			expectError: "request is nil",
			expectInfo:  nil,
		},
		{
			name: "limit too small - corrected to 20",
			req: func() *FeedbackListReq {
				r := makeSuccessFeedbackListReq()
				r.Limit = 0
				return r
			}(),
			mockError:  nil,
			mockFb:     makeSuccessFeedbackList(),
			expectInfo: makeSuccessFeedbackList(),
		},
		{
			name: "limit too large - corrected to 20",
			req: func() *FeedbackListReq {
				r := makeSuccessFeedbackListReq()
				r.Limit = 200
				return r
			}(),
			mockError:  nil,
			mockFb:     makeSuccessFeedbackList(),
			expectInfo: makeSuccessFeedbackList(),
		},
		{
			name: "fields with whitespace - trimmed",
			req: func() *FeedbackListReq {
				r := makeSuccessFeedbackListReq()
				r.Name = "  张三  "
				r.StuId = "  102301000  "
				r.OsName = "  Android  "
				r.AppVersion = "  1.2.3  "
				r.ProblemDesc = "  问题描述  "
				r.NetworkEnv = "  wifi  "
				return r
			}(),
			mockError:  nil,
			mockFb:     makeSuccessFeedbackList(),
			expectInfo: makeSuccessFeedbackList(),
		},
		{
			name: "invalid time range - end before begin",
			req: func() *FeedbackListReq {
				r := makeSuccessFeedbackListReq()
				begin := time.Now()
				end := begin.Add(-24 * time.Hour)
				r.BeginTime = &begin
				r.EndTime = &end
				return r
			}(),
			mockError:   nil,
			mockFb:      nil,
			expectError: "invalid time range",
			expectInfo:  nil,
		},
		{
			name: "invalid NetworkEnv - corrected",
			req: func() *FeedbackListReq {
				r := makeSuccessFeedbackListReq()
				r.NetworkEnv = "invalid-network"
				return r
			}(),
			mockError:  nil,
			mockFb:     makeSuccessFeedbackList(),
			expectInfo: makeSuccessFeedbackList(),
		},
		{
			name:       "empty result list - returns empty array",
			req:        makeSuccessFeedbackListReq(),
			mockError:  nil,
			mockFb:     nil,
			expectInfo: []dbmodel.FeedbackListItem{},
		},
		{
			name: "OrderDesc nil - defaults to true",
			req: func() *FeedbackListReq {
				r := makeSuccessFeedbackListReq()
				r.OrderDesc = nil
				return r
			}(),
			mockError:  nil,
			mockFb:     makeSuccessFeedbackList(),
			expectInfo: makeSuccessFeedbackList(),
		},
		{
			name: "OrderDesc false - used directly",
			req: func() *FeedbackListReq {
				r := makeSuccessFeedbackListReq()
				orderDescFalse := false
				r.OrderDesc = &orderDescFalse
				return r
			}(),
			mockError:  nil,
			mockFb:     makeSuccessFeedbackList(),
			expectInfo: makeSuccessFeedbackList(),
		},
		{
			name:        "database error",
			req:         makeSuccessFeedbackListReq(),
			mockError:   errno.InternalServiceError,
			expectError: "list feedback error",
			expectInfo:  nil,
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

			// Mock DAL：GetFeedbackList
			mockey.Mock((*oaDB.DBOA).ListFeedback).Return(tc.mockFb, 0, tc.mockError).Build()

			// 执行
			feedback, _, err := oaService.GetFeedbackList(tc.req)
			// 断言
			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, feedback)
			} else {
				assert.NoError(t, err, "should be no error")
				assert.Equal(t, tc.expectInfo, feedback)
			}
		})
	}
}
