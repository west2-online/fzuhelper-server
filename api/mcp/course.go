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
				mcp.Required(),
				mcp.Description(
					"Academic term code in the form yyyymm. "+
						"Examples: 202401 means 2024 Autumn term, 202402 means 2025 Spring term.")),
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
func handleGetCourse(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	term := request.GetString("term", "")
	if term == "" {
		return mcp.NewToolResultError("term is required"), nil
	}
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
