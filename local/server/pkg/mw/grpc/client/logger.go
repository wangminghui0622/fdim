package client

import (
	"context"

	"google.golang.org/grpc"
)

// GrpcClientLogger 简化的日志拦截?
func GrpcClientLogger() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	})
}
