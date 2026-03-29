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

	// ȴmetadataȡûϢעԱID
	combinedInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 1. gRPC metadataȡûϢcontext
		ctx = grpcinterceptor.ExtractContextFromMetadata(ctx)
		// 2. עԱûID
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

	// ӷ
	s.AddUnaryInterceptors(combinedInterceptor)

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

