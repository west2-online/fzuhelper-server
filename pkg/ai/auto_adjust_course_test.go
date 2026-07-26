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

package ai

import (
	"context"
	"os"
	"sort"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/sashabaranov/go-openai"
	"go.baoshuo.dev/llmfunc"

	"github.com/west2-online/fzuhelper-server/config"
)

const mockedAutoAdjustCourseResponse = `{
	"items": [
		{"from_date": "2025-10-01", "to_date": ""},
		{"from_date": "2025-10-02", "to_date": ""},
		{"from_date": "2025-10-03", "to_date": ""},
		{"from_date": "2025-10-04", "to_date": ""},
		{"from_date": "2025-10-05", "to_date": ""},
		{"from_date": "2025-10-06", "to_date": ""},
		{"from_date": "2025-10-07", "to_date": "2025-09-28"},
		{"from_date": "2025-10-08", "to_date": "2025-10-11"},
		{"from_date": "2025-09-28", "to_date": ""},
		{"from_date": "2025-10-11", "to_date": ""}
	]
}`

func TestAutoAdjustCourse(t *testing.T) {
	_ = config.InitForTest("common")

	// 这里不方便写死测试用的 API Key，所以需要从环境变量中读取
	config.AI.Endpoint = os.Getenv("AI_ENDPOINT")
	if config.AI.Endpoint == "" {
		config.AI.Endpoint = "https://openrouter.ai/api/v1"
	}
	config.AI.Key = os.Getenv("AI_KEY")
	if config.AI.Key == "" {
		config.AI.Key = "mock-key"
	}
	useMock := config.AI.Key == "mock-key"

	testcases := []struct {
		name     string
		content  string
		expected AutoAdjustCourseOutput
	}{
		{
			name: "关于2025年国庆节、中秋节放假课程调整的通知",
			content: `
各学院、教学单位：

根据党政办有关2025年国庆节、中秋节放假通知的精神，现将放假期间的课程调整如下：

1、10月1日（星期三）至10月8日（星期三）放假，共8天，全校本科生课程（含通识教育选修课）停课。

2、9月28日（星期日）补上10月7日（星期二）的课（2025级按原有既定安排），10月11日（星期六）补上10月8号（星期三）的课，原9月28日和10月11日的课程停课。

3、因停课受影响的教学内容，请任课老师自行调整安排。



请各学院、教学单位及时通知相关师生。



教务处

2025年9月24日
`,
			expected: AutoAdjustCourseOutput{
				Items: []AutoAdjustCourseItem{
					{
						FromDate: "2025-10-01",
						ToDate:   "",
					},
					{
						FromDate: "2025-10-02",
						ToDate:   "",
					},
					{
						FromDate: "2025-10-03",
						ToDate:   "",
					},
					{
						FromDate: "2025-10-04",
						ToDate:   "",
					},
					{
						FromDate: "2025-10-05",
						ToDate:   "",
					},
					{
						FromDate: "2025-10-06",
						ToDate:   "",
					},
					{
						FromDate: "2025-10-07",
						ToDate:   "2025-09-28",
					},
					{
						FromDate: "2025-10-08",
						ToDate:   "2025-10-11",
					},
					{
						FromDate: "2025-09-28",
						ToDate:   "",
					},
					{
						FromDate: "2025-10-11",
						ToDate:   "",
					},
				},
			},
		},
	}

	if useMock {
		defer mockey.UnPatchAll()
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if useMock {
				mockey.Mock((*llmfunc.Client).CreateChatCompletion).To(
					func(_ *llmfunc.Client, _ context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
						t.Logf("mocked CreateChatCompletion called with messages: %+v", req.Messages)
						return &openai.ChatCompletionResponse{
							Choices: []openai.ChatCompletionChoice{
								{
									Message: openai.ChatCompletionMessage{
										Role:    openai.ChatMessageRoleAssistant,
										Content: mockedAutoAdjustCourseResponse,
									},
								},
							},
						}, nil
					},
				).Build()
			}

			result, err := AutoAdjustCourse(context.Background(), AutoAdjustCourseInput{
				Title:   tc.name,
				Content: tc.content,
			})
			if err != nil {
				t.Fatalf("AutoAdjustCourse failed: %v", err)
			}
			if !equalAutoAdjustCourseOutput(result, &tc.expected) {
				t.Errorf("unexpected result:\n got: %+v\n want: %+v", result, tc.expected)
			}
			for i := range result.Items {
				if result.Items[i].FromDate == "" {
					t.Errorf("item %d: from_date is empty", i)
				}

				t.Logf("item %d: from_date=%v, to_date=%v", i, result.Items[i].FromDate, result.Items[i].ToDate)
			}
		})
	}
}

func equalAutoAdjustCourseOutput(a, b *AutoAdjustCourseOutput) bool {
	if len(a.Items) != len(b.Items) {
		return false
	}
	sort.Slice(a.Items, func(i, j int) bool {
		return a.Items[i].FromDate < a.Items[j].FromDate
	})
	sort.Slice(b.Items, func(i, j int) bool {
		return b.Items[i].FromDate < b.Items[j].FromDate
	})
	for i := range a.Items {
		if a.Items[i].FromDate != b.Items[i].FromDate || a.Items[i].ToDate != b.Items[i].ToDate {
			return false
		}
	}
	return true
}
