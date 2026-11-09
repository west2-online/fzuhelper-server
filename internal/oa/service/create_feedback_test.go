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
	"encoding/json"
	"strings"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	oaDB "github.com/west2-online/fzuhelper-server/pkg/db/oa"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func makeScreenshotJSON(count int) string {
	screenshots := make([]string, count)
	for i := range screenshots {
		screenshots[i] = "https://example.com/a.png"
	}
	data, _ := json.Marshal(screenshots)
	return string(data)
}

func makeCreateFeedbackReq() *CreateFeedbackReq {
	return &CreateFeedbackReq{
		StuId:          "102301000",
		Name:           "张三",
		College:        "计算机与大数据学院",
		ContactPhone:   "13800000000",
		ContactQQ:      "123456789",
		ContactEmail:   "123456789@qq.com",
		NetworkEnv:     "wifi",
		IsOnCampus:     true,
		OsName:         "Android",
		OsVersion:      "14",
		Manufacturer:   "Xiaomi",
		DeviceModel:    "Mi 14",
		ProblemDesc:    "闪退",
		Screenshots:    `["https://example.com/a.png"]`,
		AppVersion:     "1.2.3",
		VersionHistory: `["1.2.0","1.2.3"]`,
		NetworkTraces:  `{"url":"/api/login","code":500}`,
		Events:         `[{"event":"page_view"}]`,
		UserSettings:   `{"language":"zh"}`,
	}
}

func TestCreateFeedback(t *testing.T) {
	type testCase struct {
		name        string
		req         *CreateFeedbackReq
		mockDBError error
		mockSFError error
		expectError string
	}

	testCases := []testCase{
		{name: "success", req: makeCreateFeedbackReq()},
		{
			name: "missing required fields",
			req: func() *CreateFeedbackReq {
				req := makeCreateFeedbackReq()
				req.Name = ""
				return req
			}(),
			expectError: "missing required field name",
		},
		{
			name: "required field only contains whitespace",
			req: func() *CreateFeedbackReq {
				req := makeCreateFeedbackReq()
				req.ProblemDesc = "   "
				return req
			}(),
			expectError: "missing required field problem_desc",
		},
		{
			name: "name is too long",
			req: func() *CreateFeedbackReq {
				req := makeCreateFeedbackReq()
				req.Name = strings.Repeat("名", 31)
				return req
			}(),
			expectError: "field name cannot exceed 30 characters",
		},
		{name: "dal error", req: makeCreateFeedbackReq(), mockDBError: errno.InternalServiceError, expectError: "service.CreateFeedback: create feedback failed"},
		{
			name: "invalid NetworkEnv corrected",
			req: func() *CreateFeedbackReq {
				req := makeCreateFeedbackReq()
				req.NetworkEnv = "invalid_network"
				return req
			}(),
		},
		{
			name: "nine screenshots",
			req: func() *CreateFeedbackReq {
				req := makeCreateFeedbackReq()
				req.Screenshots = makeScreenshotJSON(9)
				return req
			}(),
		},
		{
			name: "too many screenshots",
			req: func() *CreateFeedbackReq {
				req := makeCreateFeedbackReq()
				req.Screenshots = makeScreenshotJSON(10)
				return req
			}(),
			expectError: "screenshots cannot contain more than 9 items",
		},
		{
			name: "screenshots is not an array",
			req: func() *CreateFeedbackReq {
				req := makeCreateFeedbackReq()
				req.Screenshots = `{"url":"https://example.com/a.png"}`
				return req
			}(),
			expectError: "screenshots must be a JSON string array",
		},
		{
			name: "screenshot URL is invalid",
			req: func() *CreateFeedbackReq {
				req := makeCreateFeedbackReq()
				req.Screenshots = `["/feedback/img/a.png"]`
				return req
			}(),
			expectError: "screenshot URL must be an absolute HTTP(S) URL",
		},
		{
			name: "network traces is not an object",
			req: func() *CreateFeedbackReq {
				req := makeCreateFeedbackReq()
				req.NetworkTraces = `[]`
				return req
			}(),
			expectError: "network_traces must be a JSON object",
		},
		{name: "Snowflake error", req: makeCreateFeedbackReq(), mockSFError: errno.InternalServiceError, expectError: "generate report_id failed"},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			service := NewOAService(context.Background(), "", nil, &base.ClientSet{
				SFClient: new(utils.Snowflake),
				DBClient: new(db.Database),
			})
			if tc.mockSFError != nil {
				mockey.Mock((*utils.Snowflake).NextVal).Return(int64(0), tc.mockSFError).Build()
			} else {
				mockey.Mock((*utils.Snowflake).NextVal).Return(int64(1234567890123456789), nil).Build()
			}
			mockey.Mock((*oaDB.DBOA).CreateFeedback).Return(tc.mockDBError).Build()

			_, err := service.CreateFeedback(tc.req)
			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFeedbackContactValidation(t *testing.T) {
	testCases := []struct {
		name, phone, qq, email, wantError string
	}{
		{name: "phone only", phone: "13800000000"},
		{name: "QQ only", qq: "123456789"},
		{name: "email only", email: "user+feedback@example.com"},
		{name: "multiple contacts", phone: "13800000000", qq: "123456789", email: "user@example.com"},
		{name: "trim whitespace", phone: " 13800000000 ", qq: " 123456789 ", email: " user@example.com "},
		{name: "no carrier prefix restriction", phone: "01234567890"},
		{name: "all empty", wantError: "at least one contact method"},
		{name: "all whitespace", phone: " ", qq: "\t", email: "\n", wantError: "at least one contact method"},
		{name: "phone too short", phone: "1380000000", wantError: "exactly 11 digits"},
		{name: "phone too long", phone: "138000000000", wantError: "exactly 11 digits"},
		{name: "phone contains letters", phone: "1380000000a", wantError: "exactly 11 digits"},
		{name: "phone contains space", phone: "138 0000000", wantError: "exactly 11 digits"},
		{name: "phone uses full width digits", phone: "１３８００００００００", wantError: "exactly 11 digits"},
		{name: "email missing at", email: "example.com", wantError: "valid email address"},
		{name: "email missing domain", email: "user@", wantError: "valid email address"},
		{name: "email multiple at", email: "user@@example.com", wantError: "valid email address"},
		{name: "email contains space", email: "user name@example.com", wantError: "valid email address"},
		{name: "email with display name", email: "User <user@example.com>", wantError: "valid email address"},
		{name: "valid QQ does not bypass invalid phone", phone: "invalid", qq: "123456789", wantError: "exactly 11 digits"},
		{name: "valid phone does not bypass invalid email", phone: "13800000000", email: "invalid", wantError: "valid email address"},
		{name: "QQ too long", qq: strings.Repeat("1", 33), wantError: "field contact_qq cannot exceed 32"},
		{name: "email too long", email: strings.Repeat("a", 129), wantError: "field contact_email cannot exceed 128"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := makeCreateFeedbackReq()
			req.ContactPhone, req.ContactQQ, req.ContactEmail = tc.phone, tc.qq, tc.email
			err := normalizeAndValidateFeedbackFields(req)
			if tc.wantError != "" {
				assert.ErrorContains(t, err, tc.wantError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(tc.phone), req.ContactPhone)
			assert.Equal(t, strings.TrimSpace(tc.qq), req.ContactQQ)
			assert.Equal(t, strings.TrimSpace(tc.email), req.ContactEmail)
		})
	}

	// 不提供数据库和单号依赖，验证非法联系方式会在任何写入前被拒绝。
	req := makeCreateFeedbackReq()
	req.ContactPhone, req.ContactQQ, req.ContactEmail = "", "", ""
	_, err := new(OAService).CreateFeedback(req)
	assert.ErrorContains(t, err, "at least one contact method")
}

func TestNormalizeFeedbackJSONObject(t *testing.T) {
	normalized, err := normalizeFeedbackJSONObject("", "network_traces")
	assert.NoError(t, err)
	assert.Equal(t, "{}", normalized)

	normalized, err = normalizeFeedbackJSONObject(` { "log_url": "https://example.com/log.gz" } `, "network_traces")
	assert.NoError(t, err)
	assert.JSONEq(t, `{"log_url":"https://example.com/log.gz"}`, normalized)

	_, err = normalizeFeedbackJSONObject(`"not-an-object"`, "network_traces")
	assert.ErrorContains(t, err, "network_traces must be a JSON object")
}
