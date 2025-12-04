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
	mcpgoserver "github.com/mark3labs/mcp-go/server"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
)

func GetUserInfoTool() mcpgoserver.ServerTool {
	return mcpgoserver.ServerTool{
		Tool: mcp.NewTool("get_user_info",
			mcp.WithDescription(
				"Fetch the user's personal information from the educational system. "+
					"Use this when the user asks to view their profile, personal details, or student information. "+
					"Returns user information including name, student ID, major, and other personal details."),
			mcp.WithString("user_id",
				mcp.Required(),
				mcp.Description(
					"user_id data comes from the login method response (user_id field).")),
			mcp.WithString("user_cookies",
				mcp.Required(),
				mcp.Description(
					"user_cookies data comes from the login method response (user_cookies field).")),
		),
		Handler: handleGetUserInfo,
	}
}

func handleGetUserInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 验证认证参数
	auth, errResult := ValidateAuthParams(request)
	if errResult != nil {
		return errResult, nil
	}
	ctx = WithLoginData(ctx, auth)

	userInfo, err := rpc.GetUserInfoRPC(ctx, &user.GetUserInfoRequest{})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(map[string]any{
		"user_info": userInfo,
	})
}
