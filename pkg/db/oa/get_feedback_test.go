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

package oa

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func makeFeedback() *model.Feedback {
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

func makeFeedbackList() []model.FeedbackListItem {
	items := []model.FeedbackListItem{
		{
			ReportId:    2199023256001,
			Name:        "张三",
			NetworkEnv:  model.NetworkWifi,
			ProblemDesc: "登录后提示网络异常，刷新失败",
			AppVersion:  "2.3.1",
		},
	}
	return items
}

func makeListReq() model.FeedbackListReq {
	req := model.FeedbackListReq{
		StuId: "102301000",
		Limit: 10,
	}
	return req
}

func TestDBOA_GetFeedbackById(t *testing.T) {
	type testCase struct {
		name             string
		inPutId          int64
		mockError        error
		expectingError   bool
		expectedFeedback *model.Feedback
		ErrorMsg         string
	}
	fb := makeFeedback()
	testCases := []testCase{
		{
			name:             "success",
			inPutId:          fb.ReportId,
			mockError:        nil,
			expectingError:   false,
			expectedFeedback: fb,
		},
		{
			name:             "record not found",
			inPutId:          123456789,
			mockError:        gorm.ErrRecordNotFound,
			expectingError:   true,
			expectedFeedback: nil,
		},
		{
			name:             "error",
			inPutId:          0,
			mockError:        gorm.ErrInvalidValue,
			expectingError:   true,
			expectedFeedback: nil,
			ErrorMsg:         "dal.GetFeedbackById error",
		},
	}
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockGormDB := new(gorm.DB)
			mockSnowflake := new(utils.Snowflake)
			mockDBOA := NewDBOA(mockGormDB, mockSnowflake)

			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()
			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}

				if res, ok := dest.(*model.Feedback); ok && tc.expectedFeedback != nil {
					*res = *tc.expectedFeedback
				}
				return mockGormDB
			}).Build()

			_, result, err := mockDBOA.GetFeedbackById(context.Background(), tc.inPutId, "102301000")
			if tc.expectingError {
				if err == nil {
					return
				}
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tc.ErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedFeedback, result)
			}
		})
	}
}

func TestDBOA_GetFeedbackList(t *testing.T) {
	type testCase struct {
		name                  string
		req                   model.FeedbackListReq
		mockError             error
		mockRows              []model.FeedbackListItem
		expectingError        bool
		expectedFeedback      []model.FeedbackListItem
		expectedNextPageToken int64
		expectedOrder         string
		expectedLimit         int
		expectedWhere         map[string]interface{}
		ErrorMsg              string
	}
	desc := true
	asc := false
	fb := makeFeedbackList()
	descRows := []model.FeedbackListItem{
		{ReportId: 299, Name: "A"},
		{ReportId: 298, Name: "B"},
		{ReportId: 297, Name: "C"},
	}
	ascRows := []model.FeedbackListItem{
		{ReportId: 301, Name: "A"},
		{ReportId: 302, Name: "B"},
		{ReportId: 303, Name: "C"},
	}
	testCases := []testCase{
		{
			name:             "success",
			req:              makeListReq(),
			mockRows:         fb,
			expectedFeedback: fb,
			expectedOrder:    "report_id DESC",
			expectedLimit:    11,
			expectedWhere: map[string]interface{}{
				"stu_id = ?": "102301000",
			},
		},
		{
			name: "desc cursor pagination by student",
			req: model.FeedbackListReq{
				StuId:     "102301000",
				Limit:     2,
				PageToken: 300,
				OrderDesc: &desc,
			},
			mockRows:              descRows,
			expectedFeedback:      descRows[:2],
			expectedNextPageToken: 298,
			expectedOrder:         "report_id DESC",
			expectedLimit:         3,
			expectedWhere: map[string]interface{}{
				"stu_id = ?":    "102301000",
				"report_id < ?": int64(300),
			},
		},
		{
			name: "asc cursor pagination by student",
			req: model.FeedbackListReq{
				StuId:     "102301000",
				Limit:     2,
				PageToken: 300,
				OrderDesc: &asc,
			},
			mockRows:              ascRows,
			expectedFeedback:      ascRows[:2],
			expectedNextPageToken: 302,
			expectedOrder:         "report_id ASC",
			expectedLimit:         3,
			expectedWhere: map[string]interface{}{
				"stu_id = ?":    "102301000",
				"report_id > ?": int64(300),
			},
		},
		{
			name:           "error",
			req:            makeListReq(),
			mockError:      gorm.ErrInvalidValue,
			expectingError: true,
			expectedOrder:  "report_id DESC",
			expectedLimit:  11,
			expectedWhere: map[string]interface{}{
				"stu_id = ?": "102301000",
			},
			ErrorMsg: "dal.ListFeedback error",
		},
	}
	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			observedWhere := make(map[string]interface{})
			observedOrder := ""
			observedLimit := 0
			mockGormDB := new(gorm.DB)
			mockSnowflake := new(utils.Snowflake)
			mockDBOA := NewDBOA(mockGormDB, mockSnowflake)
			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Table).To(func(name string, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Select).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				if len(args) == 1 {
					observedWhere[fmt.Sprint(query)] = args[0]
				}
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Order).To(func(value interface{}) *gorm.DB {
				observedOrder = fmt.Sprint(value)
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Limit).To(func(limit int) *gorm.DB {
				observedLimit = limit
				return mockGormDB
			}).Build()

			mockey.Mock((*gorm.DB).Find).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockGormDB.Error = tc.mockError
					return mockGormDB
				}
				if sl, ok := dest.(*[]model.FeedbackListItem); ok {
					*sl = tc.mockRows
				}
				return mockGormDB
			}).Build()

			result, nextPageToken, err := mockDBOA.ListFeedback(context.Background(), tc.req)
			assert.Equal(t, tc.expectedOrder, observedOrder)
			assert.Equal(t, tc.expectedLimit, observedLimit)
			assert.Equal(t, tc.expectedWhere, observedWhere)
			if tc.expectingError {
				if err == nil {
					return
				}
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tc.ErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedFeedback, result)
				assert.Equal(t, tc.expectedNextPageToken, nextPageToken)
			}
		})
	}
}
