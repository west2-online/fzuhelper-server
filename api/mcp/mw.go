package mcp

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	mcpgoserver "github.com/mark3labs/mcp-go/server"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
)

// 登录上下文注入中间件
func loginInjector() mcpgoserver.ToolHandlerMiddleware {
	return func(next mcpgoserver.ToolHandlerFunc) mcpgoserver.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			userID := req.GetString("user_id", "")
			userCookies := req.GetString("user_cookies", "")

			ctx = metainfoContext.WithLoginData(ctx, &model.LoginData{
				Id:      userID,
				Cookies: userCookies,
			})

			return next(ctx, req)
		}
	}
}
