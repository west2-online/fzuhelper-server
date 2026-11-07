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
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func GetCalendarTool() mcpgoserver.ServerTool {
	return mcpgoserver.ServerTool{
		Tool: mcp.NewTool(
			"get_calendar",
			mcp.WithDescription(
				"Fetch the academic calendar in ICS format for the current term. "+
					"Use this when the user asks to view the academic calendar, course schedule in calendar format, "+
					"or wants to import course schedule into calendar applications. "+
					"Returns the calendar data in ICS format (base64 encoded).",
			),
		),
		Handler: handleGetCalendar,
	}
}

func handleGetCalendar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return WithAuth(ctx, request, func(ctx context.Context, auth *AuthParams) (*mcp.CallToolResult, error) {
		stuId := auth.UserID
		if !utils.IsGraduate(auth.UserID) {
			stuId = utils.RemoveUndergraduatePrefix(auth.UserID)
		}
		icsData, err := rpc.GetCalendarRPC(ctx, &course.GetCalendarRequest{
			StuId: stuId,
		})
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(string(icsData)), nil
	})
}
