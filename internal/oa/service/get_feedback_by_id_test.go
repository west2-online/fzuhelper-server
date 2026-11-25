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

func makeFeedbackDetail() *dbmodel.Feedback {
	return &dbmodel.Feedback{
		ReportId:       1234567890123456789,
		StuId:          "102301000",
		Name:           "张三",
		College:        "计算机与大数据学院",
		ContactPhone:   "13800000000",
		ContactQQ:      "123456789",
		ContactEmail:   "123456789@qq.com",
		NetworkEnv:     dbmodel.NetworkWifi,
		IsOnCampus:     true,
		OsName:         "Android",
		OsVersion:      "14",
		Manufacturer:   "Xiaomi",
		DeviceModel:    "Mi 14",
		ProblemDesc:    "闪退",
		Screenshots:    `["https://example.com/a.png"]`,
		AppVersion:     "1.2.3",
		VersionHistory: `[]`,
		NetworkTraces:  `{}`,
		Events:         `[]`,
		UserSettings:   `{}`,
	}
}

func TestGetFeedbackByID(t *testing.T) {
	type testCase struct {
		name         string
		id           int64
		stuID        string
		mockOK       bool
		mockError    error
		mockFeedback *dbmodel.Feedback
		expectError  string
	}

	testCases := []testCase{
		{name: "success", id: 1234567890123456789, stuID: "102301000", mockOK: true, mockFeedback: makeFeedbackDetail()},
		{name: "invalid id", stuID: "102301000", expectError: "invalid id"},
		{name: "missing stu_id", id: 1234567890123456789, expectError: "missing stu_id"},
		{name: "record not found", id: 1234567890123456789, stuID: "102301000", expectError: "feedback not found"},
		{name: "database error", id: 1, stuID: "102301000", mockError: errno.InternalServiceError, expectError: "service.GetFeedbackByID: get feedback failed"},
		{
			name:   "invalid NetworkEnv corrected",
			id:     1234567890123456789,
			stuID:  "102301000",
			mockOK: true,
			mockFeedback: func() *dbmodel.Feedback {
				feedback := makeFeedbackDetail()
				feedback.NetworkEnv = "invalid"
				return feedback
			}(),
		},
		{
			name:   "empty JSON fields corrected",
			id:     1234567890123456789,
			stuID:  "102301000",
			mockOK: true,
			mockFeedback: func() *dbmodel.Feedback {
				feedback := makeFeedbackDetail()
				feedback.Screenshots = ""
				feedback.VersionHistory = ""
				feedback.NetworkTraces = ""
				feedback.Events = ""
				feedback.UserSettings = ""
				return feedback
			}(),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			service := NewOAService(context.Background(), "", nil, &base.ClientSet{
				SFClient: new(utils.Snowflake),
				DBClient: new(db.Database),
			})
			mockey.Mock((*oaDB.DBOA).GetFeedbackById).Return(tc.mockOK, tc.mockFeedback, tc.mockError).Build()

			feedback, err := service.GetFeedbackByID(tc.id, tc.stuID)
			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, feedback)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, feedback)
				assert.Equal(t, tc.id, feedback.ReportId)
			}
		})
	}
}
