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

package llm

import (
	"context"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/west2-online/fzuhelper-server/config"
)

type Client struct {
	client *openai.Client
}

func NewClient(apiKey, endpoint string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("ai: api key is empty")
	}

	cfg := openai.DefaultConfig(apiKey)
	if endpoint != "" {
		cfg.BaseURL = endpoint
	}

	return &Client{
		client: openai.NewClientWithConfig(cfg),
	}, nil
}

func NewClientFromConfig() (*Client, error) {
	if config.AI == nil {
		return nil, errors.New("ai: config.AI is nil")
	}

	return NewClient(config.AI.Key, config.AI.Endpoint)
}

func (c *Client) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ai: failed to create chat completion: %w", err)
	}

	return &resp, nil
}
