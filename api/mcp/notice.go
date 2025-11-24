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
	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
)

func GetNoticesTool() mcpgoserver.ServerTool {
	return mcpgoserver.ServerTool{
		Tool: mcp.NewTool("get_notices",
			mcp.WithDescription(
				"Fetch notices and announcements from the educational administration office. "+
					"Use this when the user asks to view official notices, announcements, or news from the academic affairs office. "+
					"Returns a list of notices with pagination support."),
			mcp.WithString("user_id",
				mcp.Required(),
				mcp.Description(
					"user_id data comes from the login method response (user_id field).")),
			mcp.WithString("user_cookies",
				mcp.Required(),
				mcp.Description(
					"user_cookies data comes from the login method response (user_cookies field).")),
			mcp.WithNumber("page",
				mcp.Description(
					"Page number for pagination. Optional: defaults to 1")),
			mcp.WithNumber("page_size",
				mcp.Description(
					"Number of notices per page. Optional: defaults to 10")),
		),
		Handler: handleGetNotices,
	}
}

func handleGetNotices(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	userID := request.GetString("user_id", "")
	userCookies := request.GetString("user_cookies", "")
	page := int64(request.GetInt("page", 1))
	pageSize := int64(request.GetInt("page_size", 10))

	if userID == "" {
		return mcp.NewToolResultError("user_id is required"), nil
	}
	if userCookies == "" {
		return mcp.NewToolResultError("user_cookies is required"), nil
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	ctx = metainfoContext.WithLoginData(ctx, &model.LoginData{
		Id:      userID,
		Cookies: userCookies,
	})

	notices, total, err := rpc.GetNoticesRPC(ctx, &common.NoticeRequest{
		PageNum: page,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(map[string]any{
		"notices":   notices,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
