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
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
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
	// 验证认证参数
	auth, errResult := ValidateAuthParams(request)
	if errResult != nil {
		return errResult, nil
	}
	ctx = WithLoginData(ctx, auth)

	term := request.GetString("term", "")

	// 如果没有指定学期，获取当前学期
	if term == "" {
		locateDate, err := rpc.GetLocateDateRPC(ctx, course.NewGetLocateDateRequest())
		if err != nil {
			return mcp.NewToolResultError("failed to determine default term: " + err.Error()), nil
		}
		if locateDate == nil || locateDate.Year == "" || locateDate.Term == "" {
			return mcp.NewToolResultError("failed to determine default term: locate date is empty"), nil
		}
		term = locateDate.Year + locateDate.Term
	}

	// 为研究生转换学期格式：本科生格式 202501 -> 研究生格式 2025-2026-1
	if utils.IsGraduate(auth.UserID) && len(term) == 6 {
		yjsTerm, err := utils.TransformSemester(term)
		if err != nil {
			return mcp.NewToolResultError("failed to transform semester: " + err.Error()), nil
		}
		term = yjsTerm
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
	// 验证认证参数
	auth, errResult := ValidateAuthParams(request)
	if errResult != nil {
		return errResult, nil
	}
	ctx = WithLoginData(ctx, auth)
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

	// 通过 GetTermsListRPC 获取学期信息
	termList, err := rpc.GetTermsListRPC(ctx, &common.TermListRequest{})
	if err != nil {
		return mcp.NewToolResultError("failed to get term list: " + err.Error()), nil
	}
	if termList == nil || termList.CurrentTerm == nil {
		return mcp.NewToolResultError("term list is empty"), nil
	}

	currentTerm := *termList.CurrentTerm
	if utils.IsGraduate(auth.UserID) {
		// 研究生格式：2025-2026-1
		yjsTerm, err := utils.TransformSemester(currentTerm)
		if err != nil {
			return mcp.NewToolResultError("failed to transform semester: " + err.Error()), nil
		}
		resp["term_formatted"] = yjsTerm
	} else {
		// 本科生格式：202501
		resp["term_formatted"] = currentTerm
	}

	return mcp.NewToolResultJSON(resp)
}
