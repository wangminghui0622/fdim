package server

import (
	"context"

	"google.golang.org/grpc"
)

// GrpcServerRequestValidate 简化的验证拦截器（不做验证，直接通过）
func GrpcServerRequestValidate() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		return handler(ctx, req)
	})
}
