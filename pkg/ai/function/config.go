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
