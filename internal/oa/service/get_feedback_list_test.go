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

func makeFeedbackListReq() *FeedbackListReq {
	return &FeedbackListReq{StuId: "102301000"}
}

func makeFeedbackListItems() []dbmodel.FeedbackListItem {
	return []dbmodel.FeedbackListItem{{
		ReportId:    2199023256001,
		Name:        "张三",
		NetworkEnv:  dbmodel.NetworkWifi,
		ProblemDesc: "登录后提示网络异常",
		AppVersion:  "2.3.1",
	}}
}

func TestGetFeedbackList(t *testing.T) {
	type testCase struct {
		name        string
		req         *FeedbackListReq
		mockError   error
		mockItems   []dbmodel.FeedbackListItem
		expectError string
		expectItems []dbmodel.FeedbackListItem
	}

	items := makeFeedbackListItems()
	testCases := []testCase{
		{name: "success", req: makeFeedbackListReq(), mockItems: items, expectItems: items},
		{name: "request is nil", expectError: "request is nil"},
		{name: "limit defaults to 20", req: makeFeedbackListReq(), mockItems: items, expectItems: items},
		{
			name: "limit over max defaults to 20",
			req: func() *FeedbackListReq {
				req := makeFeedbackListReq()
				req.Limit = 200
				return req
			}(),
			mockItems:   items,
			expectItems: items,
		},
		{
			name: "stu_id trimmed",
			req: func() *FeedbackListReq {
				req := makeFeedbackListReq()
				req.StuId = "  102301000  "
				return req
			}(),
			mockItems:   items,
			expectItems: items,
		},
		{name: "missing stu_id", req: &FeedbackListReq{}, expectError: "missing stu_id"},
		{name: "empty result returns empty array", req: makeFeedbackListReq(), expectItems: []dbmodel.FeedbackListItem{}},
		{
			name: "ascending order",
			req: func() *FeedbackListReq {
				req := makeFeedbackListReq()
				orderDesc := false
				req.OrderDesc = &orderDesc
				return req
			}(),
			mockItems:   items,
			expectItems: items,
		},
		{name: "database error", req: makeFeedbackListReq(), mockError: errno.InternalServiceError, expectError: "service.GetFeedbackList: list feedback failed"},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			service := NewOAService(context.Background(), "", nil, &base.ClientSet{
				SFClient: new(utils.Snowflake),
				DBClient: new(db.Database),
			})
			mockey.Mock((*oaDB.DBOA).ListFeedback).Return(tc.mockItems, int64(0), tc.mockError).Build()

			feedback, _, err := service.GetFeedbackList(tc.req)
			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, feedback)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectItems, feedback)
			}
		})
	}
}
