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
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type StructuredOutputOptions struct {
	SchemaName        string
	SchemaDescription string
	Strict            bool
}

func StructuredOutputs[R any](
	c *Client,
	ctx context.Context,
	messages []Message,
	options StructuredOutputOptions,
) (*R, error) {
	var result R

	reqMessages, err := buildRequestMessages(messages)
	if err != nil {
		return nil, err
	}

	responseFormat, generatedSchema, err := buildResponseFormat[R](options)
	if err != nil {
		return nil, err
	}

	raw, err := c.chat(ctx, reqMessages, responseFormat)
	if err != nil {
		return nil, err
	}

	if err = jsonschema.VerifySchemaAndUnmarshal(*generatedSchema, []byte(raw), &result); err != nil {
		return nil, fmt.Errorf("ai: verify or unmarshal structured output failed: %w", err)
	}

	return &result, nil
}

func buildResponseFormat[R any](
	options StructuredOutputOptions,
) (*openai.ChatCompletionResponseFormat, *jsonschema.Definition, error) {
	targetType := reflect.TypeOf((*R)(nil)).Elem()
	if targetType.Kind() == reflect.Invalid {
		return nil, nil, errors.New("ai: output type is invalid")
	}

	schemaName := strings.TrimSpace(options.SchemaName)
	if schemaName == "" {
		schemaName = "response"
	}

	generatedSchema, err := jsonschema.GenerateSchemaForType(reflect.New(targetType).Elem().Interface())
	if err != nil {
		return nil, nil, fmt.Errorf("ai: generate schema for output failed: %w", err)
	}

	return &openai.ChatCompletionResponseFormat{
		Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
		JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
			Name:        schemaName,
			Description: strings.TrimSpace(options.SchemaDescription),
			Schema:      generatedSchema,
			Strict:      options.Strict,
		},
	}, generatedSchema, nil
}
