package main

import (
	"context"
	"fdim/Infrastructure_service/msgtransfer/internal/config"
	"fdim/Infrastructure_service/msgtransfer/internal/handler"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v3"
)

var configFile = flag.String("f", "etc/msgtransfer.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	
	// 打印配置文件路径
	fmt.Printf("Loading config from: %s\n", *configFile)
	
	// 直接使用 yaml 解析
	data, readErr := os.ReadFile(*configFile)
	if readErr != nil {
		logx.Errorf("Failed to read config file %s: %v", *configFile, readErr)
		os.Exit(1)
	}
	fmt.Printf("Config file content length: %d bytes\n", len(data))
	
	if yamlErr := yaml.Unmarshal(data, &c); yamlErr != nil {
		logx.Errorf("Failed to parse config file: %v", yamlErr)
		os.Exit(1)
	}

	// 打印配置信息用于调试
	fmt.Printf("Config loaded: Name=%s, Mode=%s\n", c.Name, c.Mode)
	fmt.Printf("MongoDB URI: %s, Database: %s\n", c.MongoDB.URI, c.MongoDB.Database)
	fmt.Printf("NATS Brokers: %v\n", c.Nats.Brokers)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建消息传输处理器
	msgTransferHandler, err := handler.NewMsgTransferHandler(ctx, &c)
	if err != nil {
		logx.Errorf("Failed to create msg transfer handler: %v", err)
		os.Exit(1)
	}

	// 启动消息传输服务
	fmt.Printf("Starting msgtransfer service...\n")
	if err := msgTransferHandler.Start(ctx); err != nil {
		logx.Errorf("MsgTransfer service error: %v", err)
		os.Exit(1)
	}

	// 阻塞等待终止信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	logx.Infof("Received signal %v, shutting down msgtransfer...", sig)

	// 优雅关闭
	cancel()
	if err := msgTransferHandler.Stop(ctx); err != nil {
		logx.Errorf("MsgTransfer stop error: %v", err)
	}
	logx.Info("msgtransfer exited")
}
