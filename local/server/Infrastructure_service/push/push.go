package main

import (
	"context"
	"fdim/Infrastructure_service/push/internal/config"
	"fdim/Infrastructure_service/push/internal/handler"
	"fdim/Infrastructure_service/push/internal/server"
	"fdim/Infrastructure_service/push/internal/svc"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	pbpush "fdim/protocol/push"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/push.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.Infof("push config loaded: Nats.Brokers=%v, ToPushTopic=%q, ToPushGroupID=%q, ToOfflinePushTopic=%q, ToOfflinePushGroupID=%q", c.Nats.Brokers, c.Nats.ToPushTopic, c.Nats.ToPushGroupID, c.Nats.ToOfflinePushTopic, c.Nats.ToOfflinePushGroupID)

	// 创建服务上下文
	svcCtx := svc.NewServiceContext(c)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建推送处理器（用于 NATS Consumer）
	pushHandler, err := handler.NewPushHandler(ctx, &c, svcCtx)
	if err != nil {
		logx.Errorf("Failed to create push handler: %v", err)
		os.Exit(1)
	}

	// 处理信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logx.Info("Shutting down push service...")
		pushHandler.Stop(ctx)
		cancel()
	}()

	// 启动 gRPC 服务器
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pbpush.RegisterPushMsgServiceServer(grpcServer, server.NewPushServer(svcCtx))

		if c.Mode == "dev" {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 启动 NATS Consumer
	go func() {
		if err := pushHandler.Start(ctx); err != nil {
			logx.Errorf("Push handler error: %v", err)
		}
	}()

	fmt.Printf("Starting push rpc server at %s...\n", c.ListenOn)
	s.Start()
}
