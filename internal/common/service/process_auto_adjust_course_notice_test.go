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

	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	rpcmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/ai"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/jwch"
)

func TestProcessAutoAdjustCourseNotice(t *testing.T) {
	type testCase struct {
		name            string
		info            *jwch.NoticeInfo
		noticeDetailErr error
		aiResult        *ai.AutoAdjustCourseOutput
		aiErr           error
		createResp      *course.CreateAdjustCourseResponse
		createErr       error
		expectError     string
		expectRPCCalled bool
		expectItems     []*course.CreateAdjustCourseItem
	}

	mockNoticeInfo := &jwch.NoticeInfo{
		Title:    "关于课程调整的通知",
		WbTreeId: "1036",
		WbNewsId: "12345",
	}
	mockNoticeDetail := &jwch.NoticeDetail{
		Content: "课程调整内容",
	}
	successResp := &course.CreateAdjustCourseResponse{
		Base:    &rpcmodel.BaseResp{Code: errno.SuccessCode, Msg: "ok"},
		Created: new(int64),
	}
	errorResp := &course.CreateAdjustCourseResponse{
		Base: &rpcmodel.BaseResp{Code: errno.InternalServiceErrorCode, Msg: "create failed"},
	}

	testCases := []testCase{
		{
			name: "not a course adjust notice",
			info: &jwch.NoticeInfo{Title: "其他通知"},
		},
		{
			name:            "get notice detail error",
			info:            mockNoticeInfo,
			noticeDetailErr: assert.AnError,
			expectError:     "failed to get notice detail",
		},
		{
			name:        "ai error",
			info:        mockNoticeInfo,
			aiErr:       assert.AnError,
			expectError: "failed to auto adjust course",
		},
		{
			name:     "ai extracted nothing, rpc not called",
			info:     mockNoticeInfo,
			aiResult: &ai.AutoAdjustCourseOutput{Items: []ai.AutoAdjustCourseItem{}},
		},
		{
			name:            "create adjust course rpc error",
			info:            mockNoticeInfo,
			aiResult:        &ai.AutoAdjustCourseOutput{Items: []ai.AutoAdjustCourseItem{{FromDate: "2025-05-01"}}},
			createErr:       assert.AnError,
			expectError:     "failed to create adjust course",
			expectRPCCalled: true,
		},
		{
			name:            "create adjust course resp error",
			info:            mockNoticeInfo,
			aiResult:        &ai.AutoAdjustCourseOutput{Items: []ai.AutoAdjustCourseItem{{FromDate: "2025-05-01"}}},
			createResp:      errorResp,
			expectError:     "create adjust course resp error",
			expectRPCCalled: true,
		},
		{
			name:            "success with to_date set",
			info:            mockNoticeInfo,
			aiResult:        &ai.AutoAdjustCourseOutput{Items: []ai.AutoAdjustCourseItem{{FromDate: "2025-05-01", ToDate: "2025-05-08"}}},
			createResp:      successResp,
			expectRPCCalled: true,
			expectItems: []*course.CreateAdjustCourseItem{
				{FromDate: "2025-05-01", ToDate: new("2025-05-08")},
			},
		},
		{
			name:            "success with to_date empty (course canceled)",
			info:            mockNoticeInfo,
			aiResult:        &ai.AutoAdjustCourseOutput{Items: []ai.AutoAdjustCourseItem{{FromDate: "2025-05-01", ToDate: ""}}},
			createResp:      successResp,
			expectRPCCalled: true,
			expectItems: []*course.CreateAdjustCourseItem{
				{FromDate: "2025-05-01", ToDate: nil},
			},
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			courseClient := &mockCourseClient{createResp: tc.createResp, createErr: tc.createErr}
			mockClientSet := &base.ClientSet{
				DBClient:     new(db.Database),
				CourseClient: courseClient,
			}

			mockey.Mock((*jwch.Student).GetNoticeDetail).Return(mockNoticeDetail, tc.noticeDetailErr).Build()
			mockey.Mock(ai.AutoAdjustCourse).Return(tc.aiResult, tc.aiErr).Build()

			commonService := NewCommonService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			err := commonService.ProcessAutoAdjustCourseNotice(tc.info)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectRPCCalled, courseClient.createCalled)
			if tc.expectItems != nil {
				assert.Equal(t, tc.expectItems, courseClient.createReq.GetItems())
			}
		})
	}
}
