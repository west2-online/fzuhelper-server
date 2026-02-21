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
	"testing"

	"github.com/west2-online/fzuhelper-server/config"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("empty key", func(t *testing.T) {
		t.Parallel()
		_, err := NewClient("", "")
		if err == nil {
			t.Fatal("expected error when api key is empty")
		}
	})

	t.Run("valid config", func(t *testing.T) {
		t.Parallel()
		client, err := NewClient("key", "https://api.openai.com/v1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("expected non-nil client")
		}
	})
}

func TestNewClientFromConfig(t *testing.T) {
	t.Parallel()

	original := config.AI
	t.Cleanup(func() {
		config.AI = original
	})

	config.AI = nil
	if _, err := NewClientFromConfig(); err == nil {
		t.Fatal("expected error when config.AI is nil")
	}
}

func TestChatValidation(t *testing.T) {
	t.Parallel()

	var c *Client
	if _, err := c.Chat(context.Background(), "gpt-4o-mini", []Message{{Role: "user", Content: "hi"}}); err == nil {
		t.Fatal("expected error when client is nil")
	}

	client := &Client{}
	if _, err := client.Chat(context.Background(), "gpt-4o-mini", nil); err == nil {
		t.Fatal("expected error when messages is empty")
	}

	if _, err := client.Chat(context.Background(), "", []Message{{Role: "user", Content: "hi"}}); err == nil {
		t.Fatal("expected error when model name is empty")
	}
}
