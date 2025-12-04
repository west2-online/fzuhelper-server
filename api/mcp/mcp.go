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
	mcpgoserver "github.com/mark3labs/mcp-go/server"

	"github.com/west2-online/fzuhelper-server/config"
)

type Proxy struct {
	Instance *mcpgoserver.MCPServer
}

func NewMCPProxy(instance *mcpgoserver.MCPServer) *Proxy {
	return &Proxy{Instance: instance}
}

func CreateMCPProxy() *Proxy {
	server := mcpgoserver.NewMCPServer(
		config.MCP.Name,
		config.MCP.Version,
		mcpgoserver.WithToolCapabilities(true),
		// mcpgoserver.WithLogger(logger.Logger) // TODO: 引入我们自己的 logger
	)

	server.AddTools(LoginTool(), CheckSessionTool(),
		GetCourseTool(),
		GetDateTool(),
		GetScoresTool(),
		GetGPATool(),
		GetUserInfoTool(),
		GetExamRoomTool(),
		GetNoticesTool(),
		GetCalendarTool(),
	)

	return NewMCPProxy(server)
}
