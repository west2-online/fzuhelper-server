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
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/config"
)

// testInput implements FunctionInputFormatter for testing
type testInput struct {
	Query string `json:"query"`
}

func (t testInput) FunctionInput() *FunctionInput {
	return &FunctionInput{
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: t.Query},
		},
	}
}

// testOutput is a structured output type for testing
type testOutput struct {
	Result string `json:"result"`
}

// notFormatterInput does not implement FunctionInputFormatter
type notFormatterInput struct{}

func TestNewFunction(t *testing.T) {
	t.Parallel()

	t.Run("basic creation", func(t *testing.T) {
		t.Parallel()
		fn := NewFunction(
			BypassOutput(),
			Name("test"),
			Description("test function"),
			Instruction("You are a test assistant"),
			Model("test-model"),
		)
		assert.NotNil(t, fn)
	})

	t.Run("with structured output", func(t *testing.T) {
		t.Parallel()
		fn := NewFunction(
			UnmarshalOutput[testInput, testOutput](),
			Name("test-structured"),
			Description("test structured output"),
			Model("test-model"),
			StructuredOutput(true),
		)
		assert.NotNil(t, fn)
	})
}

func TestFunctionRun(t *testing.T) {
	t.Parallel()

	t.Run("input not implementing FunctionInputFormatter", func(t *testing.T) {
		t.Parallel()

		fn := NewFunction(
			func(_ *notFormatterInput, o *FunctionOutput) (*FunctionOutput, error) {
				return o, nil
			},
			Model("test-model"),
		)

		in := &notFormatterInput{}
		_, err := fn.Run(context.Background(), in)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not implement FunctionInputFormatter")
	})

	t.Run("config.AI nil returns error", func(t *testing.T) {
		t.Parallel()

		original := config.AI
		config.AI = nil
		t.Cleanup(func() { config.AI = original })

		fn := NewFunction(
			UnmarshalOutput[testInput, testOutput](),
			Model("test-model"),
		)

		in := &testInput{Query: "hello"}
		_, err := fn.Run(context.Background(), in)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create llm client")
	})
}

func TestBypassOutput(t *testing.T) {
	t.Parallel()

	handler := BypassOutput()
	input := &FunctionInput{}
	output := &FunctionOutput{FinalAnswer: "test"}

	result, err := handler(input, output)
	assert.NoError(t, err)
	assert.Equal(t, output, result)
}

func TestUnmarshalOutput(t *testing.T) {
	t.Parallel()

	t.Run("valid json", func(t *testing.T) {
		t.Parallel()
		handler := UnmarshalOutput[testInput, testOutput]()
		output := &FunctionOutput{FinalAnswer: `{"result":"hello"}`}

		result, err := handler(nil, output)
		assert.NoError(t, err)
		assert.Equal(t, "hello", result.Result)
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()
		handler := UnmarshalOutput[testInput, testOutput]()
		output := &FunctionOutput{FinalAnswer: "invalid json"}

		_, err := handler(nil, output)
		assert.Error(t, err)
	})
}
