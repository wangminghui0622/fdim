package server

import (
	"context"

	"google.golang.org/grpc"
)

// GrpcServerLogger 简化的日志拦截?
func GrpcServerLogger() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		return handler(ctx, req)
	})
}
