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

package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
)

// AuthParams 认证参数
type AuthParams struct {
	UserID      string
	UserCookies string
}

// ValidateAuthParams 验证并提取认证参数
// 从 MCP 请求中提取 user_id 和 user_cookies，并进行验证
// 返回 nil 表示验证成功，否则返回错误结果
func ValidateAuthParams(request mcp.CallToolRequest) (*AuthParams, *mcp.CallToolResult) {
	userID := request.GetString("user_id", "")
	userCookies := request.GetString("user_cookies", "")

	if userID == "" {
		return nil, mcp.NewToolResultError("user_id is required")
	}
	if userCookies == "" {
		return nil, mcp.NewToolResultError("user_cookies is required")
	}

	return &AuthParams{
		UserID:      userID,
		UserCookies: userCookies,
	}, nil
}

// WithLoginData 将认证参数添加到 context 中
// 用于统一处理 LoginData 的注入
func WithLoginData(ctx context.Context, params *AuthParams) context.Context {
	return metainfoContext.WithLoginData(ctx, &model.LoginData{
		Id:      params.UserID,
		Cookies: params.UserCookies,
	})
}
