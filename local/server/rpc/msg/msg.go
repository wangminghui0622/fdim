package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"fdim/pkg/notification"
	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/config"
	"fdim/rpc/msg/internal/server"
	"fdim/rpc/msg/internal/svc"
)

var configFile = flag.String("f", "etc/msg.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	msgServer := server.NewMsgServer(ctx)
	// 注入内部 SendMsg 函数引用（用于 HasReadReceipt 等通知，与官方一致）
	ctx.SendMsgFunc = msgServer.SendMsg
	// 初始化通知发送器（与官方一致）
	ctx.NotificationSender = notification.NewNotificationSender(&c.NotificationConfig, notification.WithLocalSendMsg(msgServer.SendMsg))

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		msg.RegisterMsgServer(grpcServer, msgServer)

		if c.Mode == "dev" {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
