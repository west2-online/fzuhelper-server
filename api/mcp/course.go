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

func GetCourseTool() mcpgoserver.ServerTool {
	return mcpgoserver.ServerTool{
		Tool: mcp.NewTool("get_course_list",
			mcp.WithDescription(
				"Fetch the user's course list for a given academic term. "+
					"Use this when the user asks to view courses/timetable for a specific term, "+
					"or before operations that need the current course roster. "+
					"Returns the course list for the given term."),
			mcp.WithString("term",
				mcp.Description(
					"Academic term code in the form yyyymm. "+
						"Examples: 202401 means 2024 Autumn term, 202402 means 2025 Spring term. "+
						"Optional: defaults to current term")),
			mcp.WithString("user_id",
				mcp.Required(),
				mcp.Description(
					"user_id data comes from the login method response (user_id field).")),
			mcp.WithString("user_cookies",
				mcp.Required(),
				mcp.Description(
					"user_cookies data comes from the login method response (user_cookies field).")),
		),
		Handler: handleGetCourse,
	}
}

func GetDateTool() mcpgoserver.ServerTool {
	return mcpgoserver.ServerTool{
		Tool: mcp.NewTool("get_date",
			mcp.WithDescription(
				"Get the current year, term, week, and date information. "+
					"Use this when other tools need to know the current term and week. "+
					"Use this when get current date information is needed. "+
					"Returns the current year, term, week, and date information."),
			mcp.WithString("user_id",
				mcp.Required(),
				mcp.Description(
					"user_id data comes from the login method response (user_id field).")),
			mcp.WithString("user_cookies",
				mcp.Required(),
				mcp.Description(
					"user_cookies data comes from the login method response (user_cookies field).")),
		),
		Handler: handleGetDate,
	}
}

func handleGetCourse(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	term := request.GetString("term", "")
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
	if term == "" {
		locateDate, err := rpc.GetLocateDateRPC(ctx, course.NewGetLocateDateRequest())
		if err != nil {
			return mcp.NewToolResultError("failed to determine default term: " + err.Error()), nil
		}
		if locateDate == nil || locateDate.Year == "" || locateDate.Week == "" {
			return mcp.NewToolResultError("failed to determine default term: locate date is empty"), nil
		}
		term = locateDate.Year + locateDate.Week // term 默认为当前学期
	}

	courseList, err := rpc.GetCourseListRPC(ctx, &course.CourseListRequest{
		Term:      term,
		IsRefresh: nil,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// 包装成JSON，JSON数组直接返回时不合法的
	resp := map[string]any{
		"term":    term,
		"courses": courseList,
	}

	return mcp.NewToolResultJSON(resp)
}

func handleGetDate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	locateDate, err := rpc.GetLocateDateRPC(ctx, course.NewGetLocateDateRequest())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if locateDate == nil {
		return mcp.NewToolResultError("locate date information is empty"), nil
	}

	resp := map[string]string{
		"year": locateDate.Year,
		"term": locateDate.Term,
		"week": locateDate.Week,
		"date": locateDate.Date,
	}

	return mcp.NewToolResultJSON(resp)
}
