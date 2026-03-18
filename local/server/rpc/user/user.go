package main

import (
	"context"
	"flag"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/grpcinterceptor"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/config"
	"fdim/rpc/user/internal/server"
	"fdim/rpc/user/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)

	// 创建组合拦截器：先从metadata提取用户信息，再注入管理员ID
	combinedInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 1. 从gRPC metadata中提取用户信息到context
		ctx = grpcinterceptor.ExtractContextFromMetadata(ctx)
		// 2. 注入管理员用户ID
		ctx = authverify.WithIMAdminUserIDs(ctx, c.AdminUserIDs)
		return handler(ctx, req)
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(svcCtx))
		user.RegisterFriendServer(grpcServer, server.NewFriendServer(svcCtx))
		user.RegisterGroupServer(grpcServer, server.NewGroupServer(svcCtx))

		if c.Mode == "dev" {
			reflection.Register(grpcServer)
		}
	})

	// 添加服务端拦截器
	s.AddUnaryInterceptors(combinedInterceptor)

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

