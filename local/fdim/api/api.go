package main

import (
	"flag"
	"fmt"
	"net/http"

	"fdim/api/internal/config"
	"fdim/api/internal/middleware"
	"fdim/api/internal/router"
	"fdim/api/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	fmt.Printf("Config loaded: UserRpc.Etcd.Hosts=%v, UserRpc.Etcd.Key=%s\n", c.UserRpc.Etcd.Hosts, c.UserRpc.Etcd.Key)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(
		func(header http.Header) {
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			header.Set("Access-Control-Allow-Headers", "Content-Type, operationID, token, Authorization")
		}, nil, "*"))
	defer server.Stop()

	// 添加认证中间件
	server.Use(middleware.AuthMiddleware(c.Secret))

	// 注册路由
	router.RegisterHandlers(server, ctx)

	fmt.Printf("Starting api server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
