package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcServerPanicCapture() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "panic: %v", r)
				fmt.Printf("[PANIC] rpc server panic: %v, method: %s\n", r, info.FullMethod)
			}
		}()
		return handler(ctx, req)
	})
}
