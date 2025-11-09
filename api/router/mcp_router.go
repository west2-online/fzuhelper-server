package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	mcpgoserver "github.com/mark3labs/mcp-go/server"
	"github.com/west2-online/fzuhelper-server/api/mcp"
)

// registerMCPRouter 注册 MCP 路由，桥接至 Hertz
func registerMCPRouter(r *server.Hertz, proxy *mcp.Proxy) {
	sseServer := mcpgoserver.NewSSEServer(
		proxy.Instance,
		mcpgoserver.WithStaticBasePath("/mcp"),
		mcpgoserver.WithSSEEndpoint("/sse"),
		mcpgoserver.WithMessageEndpoint("/message"),
	)
	r.Any("/mcp/sse", adaptor.HertzHandler(sseServer.SSEHandler()))
	r.Any("/mcp/message", adaptor.HertzHandler(sseServer.MessageHandler()))
}
