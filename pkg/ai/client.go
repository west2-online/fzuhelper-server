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
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/west2-online/fzuhelper-server/config"
)

type Message struct {
	Role    string
	Content string
}

type Client struct {
	client *openai.Client
	model  string
}

func NewClient(apiKey, modelName, endpoint string) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("ai: api key is empty")
	}
	if strings.TrimSpace(modelName) == "" {
		return nil, errors.New("ai: model name is empty")
	}

	cfg := openai.DefaultConfig(apiKey)
	if strings.TrimSpace(endpoint) != "" {
		cfg.BaseURL = endpoint
	}

	return &Client{
		client: openai.NewClientWithConfig(cfg),
		model:  modelName,
	}, nil
}

func NewClientFromConfig() (*Client, error) {
	if config.AI == nil {
		return nil, errors.New("ai: config.AI is nil")
	}

	return NewClient(config.AI.Key, config.AI.ModelName, config.AI.Endpoint)
}

func (c *Client) Chat(ctx context.Context, messages []Message) (string, error) {
	reqMessages, err := buildRequestMessages(messages)
	if err != nil {
		return "", err
	}

	return c.chat(ctx, reqMessages, nil)
}

func (c *Client) chat(
	ctx context.Context,
	reqMessages []openai.ChatCompletionMessage,
	responseFormat *openai.ChatCompletionResponseFormat,
) (string, error) {
	if c == nil || c.client == nil {
		return "", errors.New("ai: client is not initialized")
	}

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:          c.model,
		Messages:       reqMessages,
		ResponseFormat: responseFormat,
	})
	if err != nil {
		return "", fmt.Errorf("ai: create chat completion failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("ai: empty response choices")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	if content == "" {
		return "", errors.New("ai: empty response content")
	}

	return content, nil
}

func buildRequestMessages(messages []Message) ([]openai.ChatCompletionMessage, error) {
	if len(messages) == 0 {
		return nil, errors.New("ai: messages is empty")
	}

	reqMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, message := range messages {
		if strings.TrimSpace(message.Content) == "" {
			continue
		}

		role := strings.TrimSpace(message.Role)
		if role == "" {
			role = openai.ChatMessageRoleUser
		}

		reqMessages = append(reqMessages, openai.ChatCompletionMessage{
			Role:    role,
			Content: message.Content,
		})
	}

	if len(reqMessages) == 0 {
		return nil, errors.New("ai: valid messages is empty")
	}

	return reqMessages, nil
}
