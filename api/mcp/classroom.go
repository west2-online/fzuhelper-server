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
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	mcpgoserver "github.com/mark3labs/mcp-go/server"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
)

func GetExamRoomTool() mcpgoserver.ServerTool {
	return mcpgoserver.ServerTool{
		Tool: mcp.NewTool(
			"get_exam_room",
			mcp.WithDescription(
				"Fetch the user's exam room information and schedule. "+
					"Use this when the user asks to view exam locations, exam schedule, or examination room details. "+
					"Returns exam room information including location, time, seat number, and course details.",
			),
			mcp.WithString("term",
				mcp.Description(
					"Academic term code in the form yyyymm. "+
						"Examples: 202401 means 2024 Autumn term, 202402 means 2025 Spring term. "+
						"Optional: defaults to current term",
				)),
		),
		Handler: handleGetExamRoom,
	}
}

func handleGetExamRoom(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return WithAuth(ctx, request, func(ctx context.Context, _ *AuthParams) (*mcp.CallToolResult, error) {
		term := request.GetString("term", "")

		if term == "" {
			locateDate, err := rpc.GetLocateDateRPC(ctx, course.NewGetLocateDateRequest())
			if err != nil {
				return nil, fmt.Errorf("failed to determine default term: %w", err)
			}
			if locateDate == nil || locateDate.Year == "" || locateDate.Term == "" {
				return nil, errors.New("failed to determine default term: locate date is empty")
			}
			term = locateDate.Year + locateDate.Term
		}

		examRooms, err := rpc.GetExamRoomInfoRPC(ctx, &classroom.ExamRoomInfoRequest{
			Term: term,
		})
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultJSON(map[string]any{
			"term":       term,
			"exam_rooms": examRooms,
		})
	})
}
