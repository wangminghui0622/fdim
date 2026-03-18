package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Cache cache.CacheConf
	// LongConnServer WebSocket 长连接服务器配置
	LongConnServer struct {
		Ports               []int
		WebsocketMaxConnNum int64
		WebsocketTimeout    int
		WebsocketMaxMsgLen  int
	}
	
	// RPC 客户端配置
	AuthRpc zrpc.RpcClientConf
}
