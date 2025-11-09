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
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	mcpgoserver "github.com/mark3labs/mcp-go/server"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

type IdentifierData struct {
	ID      string `json:"user_id"`
	Cookies string `json:"user_cookies"`
}

func LoginTool() mcpgoserver.ServerTool {
	return mcpgoserver.ServerTool{
		Tool: mcp.NewTool("login",
			mcp.WithDescription("Use this tool when the user wants to log into the educational system. "+
				"Call this when: user mentions logging in, needs to authenticate, "+
				"or when other tools fail due to no active session. "+
				"If JWCH_STUDENT_ID and JWCH_PASSWORD environment variables are set, "+
				"this happens automatically on startup. Returns success message on successful login."),
			mcp.WithString("student_id",
				mcp.Required(),
				mcp.Description("Student ID for authentication (optional if FZUHELPER_STUDENT_ID env var is set)"),
			),
			mcp.WithString("password",
				mcp.Required(),
				mcp.Description("Password for authentication (optional if FZUHELPER_STUDENT_PASSWORD env var is set)"),
			),
			mcp.WithString("student_type",
				mcp.Required(),
				mcp.Description("StudentType for authentication (optional if FZUHELPER_STUDENT_TYPE env var is set)"),
			),
		),
		Handler: handleLogin,
	}
}

func CheckSessionTool() mcpgoserver.ServerTool {
	return mcpgoserver.ServerTool{
		Tool: mcp.NewTool("check_session",
			mcp.WithDescription("Use this tool to verify if the current login session is still valid. "+
				"Call this when: user asks about connection status, before performing operations after a long idle "+
				"period, or to troubleshoot authentication issues. Returns session validity status."),
			mcp.WithString("user_id",
				mcp.Required(),
				mcp.Description("user_id data comes from login method response, (user_cookies field)"),
			),
			mcp.WithString("user_cookies",
				mcp.Required(),
				mcp.Description("user_cookies data comes from login method response, (user_cookies field)"),
			),
		),
		Handler: handleCheckSession,
	}
}

func handleLogin(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	studentID := request.GetString("student_id", "")
	password := request.GetString("password", "")
	studentType := request.GetString("student_type", "")
	if studentID == "" {
		return mcp.NewToolResultError("student_id is required (provide as parameter or set JWCH_STUDENT_ID environment variable)"), nil
	}
	if password == "" {
		return mcp.NewToolResultError("password is required (provide as parameter or set JWCH_PASSWORD environment variable)"), nil
	}

	var id, cookies string
	var err error
	switch studentType {
	case "2":
		id, cookies, err = rpc.GetLoginDataForYJSYRPC(ctx, &user.GetLoginDataForYJSYRequest{
			Id:       studentID,
			Password: password,
		})
	default:
		id, cookies, err = rpc.GetLoginDataRPC(ctx, &user.GetLoginDataRequest{
			Id:       studentID,
			Password: password,
		})
	}
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Internal RPC request failed: %v", err)), nil
	}

	return mcp.NewToolResultJSON(IdentifierData{
		ID:      id,
		Cookies: cookies,
	})
}

func handleCheckSession(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	userID := request.GetString("user_id", "")
	userCookies := request.GetString("user_cookies", "")
	if userID == "" {
		return mcp.NewToolResultError("user_id is required"), nil
	}
	if userCookies == "" {
		return mcp.NewToolResultError("user_cookies is required"), nil
	}

	err := jwch.NewStudent().WithLoginData(userID, utils.ParseCookies(userCookies)).CheckSession()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("CheckSession failed: %v", err)), nil
	}
	return mcp.NewToolResultText("Session alive"), nil
}
