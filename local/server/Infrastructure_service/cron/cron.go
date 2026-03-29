package main

import (
	"context"
	"fdim/Infrastructure_service/cron/internal/config"
	"fdim/Infrastructure_service/cron/internal/handler"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/cron.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建上下?
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建 cron 处理?
	cronHandler, err := handler.NewCronHandler(ctx, &c)
	if err != nil {
		logx.Errorf("Failed to create cron handler: %v", err)
		os.Exit(1)
	}

	// 处理信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logx.Info("Shutting down cron service...")
		cancel()
	}()

	fmt.Printf("Starting cron service...\n")

	// 启动 cron 服务
	if err := cronHandler.Start(ctx); err != nil {
		logx.Errorf("Cron service error: %v", err)
		os.Exit(1)
	}
}
