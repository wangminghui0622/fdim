package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/config"
	"fdim/rpc/conversation/internal/server"
	"fdim/rpc/conversation/internal/svc"
)

var configFile = flag.String("f", "etc/conversation.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		conversation.RegisterConversationServer(grpcServer, server.NewConversationServer(ctx))

		if c.Mode == "dev" {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
