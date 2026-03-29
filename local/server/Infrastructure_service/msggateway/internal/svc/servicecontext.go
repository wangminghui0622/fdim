package svc

import (
	"fdim/Infrastructure_service/msggateway/internal/config"
	"fdim/Infrastructure_service/msggateway/internal/ws"
	"fdim/pkg/cache"
	"fdim/protocol/auth"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	AuthClient auth.AuthClient
	WsServer   *ws.WsServer
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始?Auth RPC 客户?
	var authClient auth.AuthClient
	if c.AuthRpc.Etcd.Key != "" || c.AuthRpc.Target != "" {
		authClient = auth.NewAuthClient(zrpc.MustNewClient(c.AuthRpc).Conn())
	}

	// 初始?Redis 客户端（用于在线状态缓存等?
	redisClient := cache.NewRedisClient(c.Cache)

	// 初始?WebSocket 服务?
	wsServer := ws.NewWsServer(&c, authClient, redisClient)

	return &ServiceContext{
		Config:     c,
		AuthClient: authClient,
		WsServer:   wsServer,
	}
}
