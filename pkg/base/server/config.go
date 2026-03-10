package server

import (
	"net"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/remote/codec/thrift"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// config.go RPC 服务器配置，配置应当只在 cmd 包中调用

// AssembleCommonServerConfig 组装通用 RPC 服务器配置
func AssembleCommonServerConfig(serviceName string, addr net.Addr, r registry.Registry) []server.Option {
	return append(commonServerConfig(),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
		// server.WithMuxTransport(), // 使用 Mux 传输会和流式传输冲突，这里考虑禁用流式传输
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
	)
}

// commonServerConfig 返回通用的服务端配置
func commonServerConfig() []server.Option {

	return []server.Option{
		server.WithLimit(&limit.Option{
			MaxConnections: constants.MaxConnections, // 最大连接数
			MaxQPS:         constants.MaxQPS,         // 最大 QPS
		}),
		server.WithPayloadCodec(thrift.NewThriftCodecWithConfig(thrift.FrugalReadWrite)), // 使用 Frugal 进行解编码
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),                          // 使用 TTHeader 进行元数据传输
	}
}
