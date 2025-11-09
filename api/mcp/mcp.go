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
		//mcpgoserver.WithLogger(logger.Logger) // TODO: 引入我们自己的 logger
	)

	server.AddTools(LoginTool(), CheckSessionTool())

	return NewMCPProxy(server)
}
