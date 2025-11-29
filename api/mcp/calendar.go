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
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
)

func GetCalendarTool() mcpgoserver.ServerTool {
	return mcpgoserver.ServerTool{
		Tool: mcp.NewTool("get_calendar",
			mcp.WithDescription(
				"Fetch the academic calendar in ICS format for the current term. "+
					"Use this when the user asks to view the academic calendar, course schedule in calendar format, "+
					"or wants to import course schedule into calendar applications. "+
					"Returns the calendar data in ICS format (base64 encoded)."),
			mcp.WithString("user_id",
				mcp.Required(),
				mcp.Description(
					"user_id data comes from the login method response (user_id field).")),
			mcp.WithString("user_cookies",
				mcp.Required(),
				mcp.Description(
					"user_cookies data comes from the login method response (user_cookies field).")),
		),
		Handler: handleGetCalendar,
	}
}

func handleGetCalendar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	userID := request.GetString("user_id", "")
	userCookies := request.GetString("user_cookies", "")

	if userID == "" {
		return mcp.NewToolResultError("user_id is required"), nil
	}
	if userCookies == "" {
		return mcp.NewToolResultError("user_cookies is required"), nil
	}
	ctx = metainfoContext.WithLoginData(ctx, &model.LoginData{
		Id:      userID,
		Cookies: userCookies,
	})
	//研究生查询课表没有任何返回结果（即便用调试工具改到有课的学期），暂不知晓原因。
	icsData, err := rpc.GetCalendarRPC(ctx, &course.GetCalendarRequest{
		StuId: userID,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(icsData)), nil
}
