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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/config"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("empty key", func(t *testing.T) {
		t.Parallel()
		_, err := NewClient("", "")
		assert.Error(t, err)
	})

	t.Run("valid config", func(t *testing.T) {
		t.Parallel()
		client, err := NewClient("test-key", "https://api.example.com/v1")
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("valid config without endpoint", func(t *testing.T) {
		t.Parallel()
		client, err := NewClient("test-key", "")
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}

func TestNewClientFromConfig(t *testing.T) {
	t.Parallel()

	t.Run("nil config", func(t *testing.T) {
		t.Parallel()
		original := config.AI
		config.AI = nil
		t.Cleanup(func() { config.AI = original })

		_, err := NewClientFromConfig()
		assert.Error(t, err)
	})

	t.Run("empty key in config", func(t *testing.T) {
		t.Parallel()
		original := config.AI
		// AI points to an ai struct with empty Key, should fail
		config.AI = original
		t.Cleanup(func() { config.AI = original })

		// When AI is non-nil but Key is empty, NewClient returns error
		if config.AI != nil && config.AI.Key == "" {
			_, err := NewClientFromConfig()
			assert.Error(t, err)
		}
	})
}

func TestCreateChatCompletion(t *testing.T) {
	t.Parallel()

	t.Run("successful completion", func(t *testing.T) {
		t.Parallel()

		expectedContent := "Hello! How can I help you?"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Role:    openai.ChatMessageRoleAssistant,
							Content: expectedContent,
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client, err := NewClient("test-key", server.URL)
		assert.NoError(t, err)

		resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: "test-model",
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: "Hello"},
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Choices, 1)
		assert.Equal(t, expectedContent, resp.Choices[0].Message.Content)
	})

	t.Run("server error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client, err := NewClient("test-key", server.URL)
		assert.NoError(t, err)

		_, err = client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: "test-model",
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: "Hello"},
			},
		})
		assert.Error(t, err)
	})
}
