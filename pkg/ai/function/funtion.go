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

package function

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/west2-online/fzuhelper-server/pkg/ai/internal"
	"github.com/west2-online/fzuhelper-server/pkg/ai/llm"
)

type FunctionInput struct {
	Messages []openai.ChatCompletionMessage `json:"messages"`
}

type FunctionOutput struct {
	FinalAnswer string `json:"final_answer"`
}

func BypassOutput() func(i *FunctionInput, o *FunctionOutput) (*FunctionOutput, error) {
	return func(i *FunctionInput, o *FunctionOutput) (*FunctionOutput, error) {
		return o, nil
	}
}

func UnmarshalOutput[T any, R any]() func(_ *T, o *FunctionOutput) (*R, error) {
	return func(_ *T, o *FunctionOutput) (*R, error) {
		r := new(R)
		return r, json.Unmarshal([]byte(o.FinalAnswer), r)
	}
}

type FunctionInputFormatter interface {
	FunctionInput() *FunctionInput
}

type Function[T any, R any] struct {
	output func(*T, *FunctionOutput) (*R, error)
	config *FunctionConfig
}

func NewFunction[T any, R any](
	handler func(i *T, o *FunctionOutput) (*R, error),
	opts ...internal.Option[*FunctionConfig],
) *Function[T, R] {
	config := &FunctionConfig{}
	for _, opt := range opts {
		opt.Apply(config)
	}
	if config.structuredOutput {
		schema, err := jsonschema.GenerateSchemaForType(*new(R))
		if err != nil {
			panic(err) // 生成不出来肯定是代码写得有问题，直接 panic 吧
		}
		OutputSchema(schema).Apply(config)
	}
	return &Function[T, R]{
		output: handler,
		config: config,
	}
}

func (f *Function[T, R]) Run(ctx context.Context, in *T) (*R, error) {
	var input *FunctionInput
	ok := false

	if formatter, hasFormatter := any(*in).(FunctionInputFormatter); hasFormatter {
		input = formatter.FunctionInput()
		ok = true
	}

	if !ok {
		return nil, fmt.Errorf("input does not implement FunctionInputFormatter")
	}

	client, err := llm.NewClientFromConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create llm client: %w", err)
	}

	messages := make([]openai.ChatCompletionMessage, 0, len(input.Messages)+1)
	if f.config.instruction != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: f.config.instruction,
		})
	}
	messages = append(messages, input.Messages...)

	req := openai.ChatCompletionRequest{
		Model:    f.config.model,
		Messages: messages,
	}

	if f.config.structuredOutput {
		req.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:        f.config.name,
				Description: f.config.description,
				Schema:      f.config.outputSchema,
				Strict:      true,
			},
		}
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from LLM")
	}

	output := &FunctionOutput{
		FinalAnswer: resp.Choices[0].Message.Content,
	}

	return f.output(in, output)
}
