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
	"fmt"

	"github.com/sashabaranov/go-openai"
	"go.baoshuo.dev/llmfunc"
)

const AutoAdjustCourseInstruction = `
你是一名调课通知解析助手，负责从调课通知的标题和内容中提取调课信息。调课信息包括调整前的日期和调整后的日期（如果当日课程取消则留空）。
请根据以下要求提取信息：
1. 输入的通知中可能包含多条调课信息，请提取所有的调课信息。
2. 日期格式为 YYYY-MM-DD。
3. 对于通知中提到放假的日期，如果没有提到其对应的补课日期则表示课程取消，产生一条调课信息，调整前的日期是放假的日期，调整后的日期留空。
4. 如果通知中提到日期X补上日期Y的课，则原来应该在日期X上课的课程取消，产生两条调课信息：一条是日期X的课程取消（调整后的日期留空），另一条是日期Y的课程调整（调整前的日期是X，调整后的日期是Y）。
5. 输出的结果应该是一个包含多个调课信息的列表，每条调课信息包含调整前的日期和调整后的日期（如果当日课程取消则留空）。
6. 如果无法提取到有效的调课信息，请返回一个空列表。

【示例输入】
# 关于2026年元旦放假课程调整的通知

各学院，各教学单位：

根据党政办有关2026年元旦放假通知的精神，现将放假期间的课程调整如下：
1、1月1日（周四）至1月3日（周六）放假，共3天，全校本科生课程（含通识教育选修课）停课。
2、1月4日（周日）补上1月2日（周五）的课，原1月4日的课程停课。
3、因停课受影响的教学内容，请任课老师自行调整安排。

请各学院、教学单位及时通知相关师生。

教务处
2025年12月29日

【解析结果】
第一条信息提示1月1日到1月3日放假，停课，得到一个1月1日、1月2日、1月3日均取消的记录列表；
第二条信息提示1月4日补上1月2日的课，原1月4日的课程停课，得到一个1月4日取消的记录，将1月2日的取消记录改为调整到1月4日。
综上，得到如示例输出所示的结果。

【示例输出】
[
  {
    "from_date": "2026-01-01",
    "to_date": "",
  },
  {
    "from_date": "2026-01-02",
    "to_date": "2026-01-04",
  },
  {
    "from_date": "2026-01-03",
    "to_date": "",
  },
  {
    "from_date": "2026-01-04",
    "to_date": "",
  }
]
`

const autoAdjustCourseTemperature = 0.2

type AutoAdjustCourseInput struct {
	Title   string `json:"title"   description:"通知标题"`
	Content string `json:"content" description:"通知内容"`
}

func (i AutoAdjustCourseInput) FunctionInput() *llmfunc.FunctionInput {
	return &llmfunc.FunctionInput{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("#%s\n\n%s", i.Title, i.Content),
			},
		},
	}
}

type AutoAdjustCourseItem struct {
	FromDate string `json:"from_date" description:"调整前课程本应上课的日期，格式为 YYYY-MM-DD"`
	ToDate   string `json:"to_date"   description:"调整后的实际上课日期，格式为 YYYY-MM-DD，如果课程取消则留空"`
}

type AutoAdjustCourseOutput struct {
	Items []AutoAdjustCourseItem `json:"items"`
}

func AutoAdjustCourse(ctx context.Context, input AutoAdjustCourseInput) (*AutoAdjustCourseOutput, error) {
	f := NewFunction(
		llmfunc.UnmarshalOutput[AutoAdjustCourseInput, AutoAdjustCourseOutput](),
		llmfunc.Name("auto_adjust_course"),
		llmfunc.Description("解析调课通知提取调课信息"),
		llmfunc.Instruction(AutoAdjustCourseInstruction),
		llmfunc.StructuredOutput(true),
		llmfunc.Model("deepseek-ai/DeepSeek-V3.2"),
		llmfunc.Temperature(autoAdjustCourseTemperature),
	)

	output, err := f.Run(ctx, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to run auto adjust course function: %w", err)
	}

	return output, nil
}
