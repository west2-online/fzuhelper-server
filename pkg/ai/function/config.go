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
	"github.com/sashabaranov/go-openai/jsonschema"

	"github.com/west2-online/fzuhelper-server/pkg/ai/internal"
)

type FunctionConfig struct {
	name             string
	description      string
	instruction      string
	structuredOutput bool
	outputSchema     *jsonschema.Definition
	model            string
	temperature      float32
}

func Name(name string) internal.Option[*FunctionConfig] {
	return internal.NewApplyOption(func(config *FunctionConfig) {
		config.name = name
	})
}

func Description(description string) internal.Option[*FunctionConfig] {
	return internal.NewApplyOption(func(config *FunctionConfig) {
		config.description = description
	})
}

func Instruction(instruction string) internal.Option[*FunctionConfig] {
	return internal.NewApplyOption(func(config *FunctionConfig) {
		config.instruction = instruction
	})
}

func StructuredOutput(structured bool) internal.Option[*FunctionConfig] {
	return internal.NewApplyOption(func(config *FunctionConfig) {
		config.structuredOutput = structured
	})
}

func OutputSchema(schema *jsonschema.Definition) internal.Option[*FunctionConfig] {
	return internal.NewApplyOption(func(config *FunctionConfig) {
		config.outputSchema = schema
	})
}

func Model(model string) internal.Option[*FunctionConfig] {
	return internal.NewApplyOption(func(config *FunctionConfig) {
		config.model = model
	})
}

func Temperature(temp float32) internal.Option[*FunctionConfig] {
	return internal.NewApplyOption(func(config *FunctionConfig) {
		config.temperature = temp
	})
}
