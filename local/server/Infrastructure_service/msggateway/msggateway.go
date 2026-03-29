package main

import (
	"context"
	"fdim/Infrastructure_service/msggateway/internal/config"
	"fdim/Infrastructure_service/msggateway/internal/server"
	"fdim/Infrastructure_service/msggateway/internal/svc"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"fdim/protocol/msggateway"
)

var configFile = flag.String("f", "etc/msggateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 调试：打印配
	fmt.Printf("Config loaded: LongConnServer.Ports=%v\n", c.LongConnServer.Ports)

	svcCtx := svc.NewServiceContext(c)

	// 启动 gRPC 服务
	grpcServer := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		msggateway.RegisterMsgGatewayServer(grpcServer, server.NewMsgGatewayServer(svcCtx))

		if c.Mode == "dev" {
			reflection.Register(grpcServer)
		}
	})
	defer grpcServer.Stop()

	// 创建上下
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动 WebSocket 服务
	if len(c.LongConnServer.Ports) > 0 {
		port := c.LongConnServer.Ports[0]
		fmt.Printf("Launching WebSocket server on port %d...\n", port)
		go func() {
			fmt.Printf("WebSocket goroutine started, port %d\n", port)
			if err := svcCtx.WsServer.Run(ctx, port); err != nil {
				logx.Errorf("WebSocket server error: %v", err)
				fmt.Printf("WebSocket server error: %v\n", err)
			}
		}()
	} else {
		fmt.Println("No WebSocket ports configured!")
	}

	// 处理信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logx.Info("Shutting down...")
		cancel()
	}()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	if len(c.LongConnServer.Ports) > 0 {
		fmt.Printf("Starting WebSocket server on port %d...\n", c.LongConnServer.Ports[0])
	}
	grpcServer.Start()
}
