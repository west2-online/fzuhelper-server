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
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
)

func TestGetJobFair(t *testing.T) {
	type testCase struct {
		name                string
		month               string
		responseStatus      int
		responseBody        string
		requestError        error
		expectRequestBody   string
		expectEvents        []*model.JobFairEvent
		expectErrorContains string
	}

	testCases := []testCase{
		{
			name:              "successfully_normalizes_source_events",
			month:             "2026-09",
			responseStatus:    200,
			expectRequestBody: "dateday=2026%2F09",
			responseBody: `{
				"success": true,
				"zhaopinhui_keynoteList": [
					{
						"id": "lecture-id",
						"title": "A&amp;B&mdash;2027届宣讲会",
						"place": null,
						"start_time": "20260904",
						"time": "19:00",
						"zphval": "3"
					},
					{
						"id": "job-fair-id",
						"title": "秋季招聘会",
						"place": "  旗山校区  ",
						"start_time": "20260905",
						"time": "09:30",
						"zphval": "1"
					}
				]
			}`,
			expectEvents: []*model.JobFairEvent{
				{
					Id:        "lecture-id",
					Title:     "A&B—2027届宣讲会",
					Place:     "",
					Time:      "19:00",
					StartsAt:  "2026-09-04T19:00:00+08:00",
					DateKey:   "2026-09-04",
					DetailUrl: "http://fjrclh.fzu.edu.cn/cms/xjhdetail.html?id=lecture-id",
				},
				{
					Id:        "job-fair-id",
					Title:     "秋季招聘会",
					Place:     "旗山校区",
					Time:      "09:30",
					StartsAt:  "2026-09-05T09:30:00+08:00",
					DateKey:   "2026-09-05",
					DetailUrl: "http://fjrclh.fzu.edu.cn/cms/zphdetail.html?id=job-fair-id",
				},
			},
		},
		{
			name:                "rejects_invalid_month",
			month:               "2026/09",
			expectErrorContains: "invalid month",
		},
		{
			name:                "returns_request_error",
			month:               "2026-09",
			requestError:        errors.New("connection refused"),
			expectRequestBody:   "dateday=2026%2F09",
			expectErrorContains: "request source failed",
		},
		{
			name:                "rejects_non_success_status",
			month:               "2026-09",
			responseStatus:      503,
			expectRequestBody:   "dateday=2026%2F09",
			expectErrorContains: "unexpected status code 503",
		},
		{
			name:                "rejects_invalid_json",
			month:               "2026-09",
			responseStatus:      200,
			responseBody:        "not-json",
			expectRequestBody:   "dateday=2026%2F09",
			expectErrorContains: "unmarshal response failed",
		},
		{
			name:                "rejects_unsuccessful_source_response",
			month:               "2026-09",
			responseStatus:      200,
			responseBody:        `{"success":false,"zhaopinhui_keynoteList":[]}`,
			expectRequestBody:   "dateday=2026%2F09",
			expectErrorContains: "source response was unsuccessful",
		},
		{
			name:           "rejects_invalid_event_time",
			month:          "2026-09",
			responseStatus: 200,
			responseBody: `{
				"success": true,
				"zhaopinhui_keynoteList": [{
					"id": "bad-time",
					"title": "无效时间",
					"place": "",
					"start_time": "20260999",
					"time": "19:00",
					"zphval": "3"
				}]
			}`,
			expectRequestBody:   "dateday=2026%2F09",
			expectErrorContains: "parse event time failed",
		},
	}

	httpClient, err := client.NewClient()
	require.NoError(t, err)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*client.Client).Do).To(
				func(_ *client.Client, _ context.Context, req *protocol.Request, resp *protocol.Response) error {
					assert.Equal(t, tc.expectRequestBody, string(req.Body()))
					if tc.requestError != nil {
						return tc.requestError
					}
					resp.SetStatusCode(tc.responseStatus)
					resp.SetBodyString(tc.responseBody)
					return nil
				},
			).Build()

			commonService := NewCommonService(
				context.Background(),
				&base.ClientSet{HzClient: httpClient},
				new(taskqueue.BaseTaskQueue),
			)
			events, err := commonService.GetJobFair(tc.month)

			if tc.expectErrorContains != "" {
				assert.ErrorContains(t, err, tc.expectErrorContains)
				assert.Nil(t, events)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectEvents, events)
		})
	}
}
