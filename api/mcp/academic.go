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
	"github.com/west2-online/fzuhelper-server/kitex_gen/academic"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func GetScoresTool() mcpgoserver.ServerTool {
	return mcpgoserver.ServerTool{
		Tool: mcp.NewTool("get_scores",
			mcp.WithDescription(
				"Fetch the user's score records for a given academic term. "+
					"Use this when the user asks to view grades/scores for a specific term. "+
					"Returns the score list for the given term."),
			mcp.WithString("user_id",
				mcp.Required(),
				mcp.Description(
					"user_id data comes from the login method response (user_id field).")),
			mcp.WithString("user_cookies",
				mcp.Required(),
				mcp.Description(
					"user_cookies data comes from the login method response (user_cookies field).")),
		),
		Handler: handleGetScores,
	}
}

func GetGPATool() mcpgoserver.ServerTool {
	return mcpgoserver.ServerTool{
		Tool: mcp.NewTool("get_gpa",
			mcp.WithDescription(
				"Fetch the user's GPA (Grade Point Average) information for a given academic term. "+
					"Use this when the user asks to view GPA or grade point information. "+
					"Returns the GPA data including overall and term-specific GPA. "+
					"Note: This feature is only available for undergraduate students. Graduate students should use get_scores instead."),
			mcp.WithString("user_id",
				mcp.Required(),
				mcp.Description(
					"user_id data comes from the login method response (user_id field).")),
			mcp.WithString("user_cookies",
				mcp.Required(),
				mcp.Description(
					"user_cookies data comes from the login method response (user_cookies field).")),
		),
		Handler: handleGetGPA,
	}
}

func handleGetScores(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 验证认证参数
	auth, errResult := ValidateAuthParams(request)
	if errResult != nil {
		return errResult, nil
	}
	ctx = WithLoginData(ctx, auth)

	scores, err := rpc.GetScoresRPC(ctx, &academic.GetScoresRequest{})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resp := map[string]any{
		"scores": scores,
	}

	return mcp.NewToolResultJSON(resp)
}

func handleGetGPA(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 验证认证参数
	auth, errResult := ValidateAuthParams(request)
	if errResult != nil {
		return errResult, nil
	}

	// 研究生系统不支持 GPA 查询
	if utils.IsGraduate(auth.UserID) {
		return mcp.NewToolResultError("GPA query is not supported for graduate students. The graduate student system does not provide GPA information. Please use get_scores to view your grades instead."), nil
	}

	ctx = WithLoginData(ctx, auth)

	gpa, err := rpc.GetGPARPC(ctx, &academic.GetGPARequest{})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(map[string]any{
		"gpa": gpa,
	})
}
