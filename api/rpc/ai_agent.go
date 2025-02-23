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

package rpc

import (
	"context"
	"errors"
	"io"

	"github.com/west2-online/fzuhelper-server/kitex_gen/ai_agent"
	"github.com/west2-online/fzuhelper-server/pkg/base/client/grpc"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func InitAiAgentRPC() {
	c, err := grpc.InitAiAgentRPC()
	if err != nil {
		logger.Fatalf("api.rpc.ai_agent InitUserRPC failed, err is %v", err)
	}
	aiAgentClient = c
}

func InitAiAgentStreamRPC() {
}

func TestRPC(ctx context.Context, req *ai_agent.ChatRequest) (string, error) {
	resp, err := aiAgentClient.Cli.Single(ctx, req)
	if err != nil {
		logger.Errorf("TestRPC: RPC called failed: %v", err.Error())
		return "", errno.InternalServiceError.WithError(err)
	}
	return resp.Answer, nil
}

func StreamChatRPC(ctx context.Context, req *ai_agent.ChatRequest) (chan string, error) {
	s, err := aiAgentClient.Cli.StreamChat(ctx, req)
	if err != nil {
		logger.Errorf("StreamChatRPC: RPC create stream failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	ch := make(chan string)
	go func() {
		defer close(ch) // 确保在 goroutine 退出时关闭 channel

		for {
			resp, err := s.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else {
					logger.Errorf("StreamChatRPC: receive error: %v", err)
				}
				break
			}
			if resp.GetEndOfStream() {
				break
			}
			ch <- resp.GetAnswer()
		}
	}()
	return ch, nil
}
