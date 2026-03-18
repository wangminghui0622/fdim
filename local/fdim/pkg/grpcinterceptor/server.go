package grpcinterceptor

import (
	"context"

	"fdim/protocol/constant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GrpcServer 返回gRPC服务端拦截器，用于从gRPC metadata中提取用户信息并设置到context
func GrpcServer() grpc.UnaryServerInterceptor {
	return RpcServerInterceptor
}

// RpcServerInterceptor gRPC服务端拦截器，从gRPC metadata中提取用户信息并设置到context
func RpcServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = ExtractContextFromMetadata(ctx)
	return handler(ctx, req)
}

// ExtractContextFromMetadata 从gRPC metadata中提取用户信息并设置到context
func ExtractContextFromMetadata(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	// 提取 operationID
	if values := md.Get(constant.OperationID); len(values) > 0 {
		ctx = context.WithValue(ctx, constant.OperationID, values[0])
	}

	// 提取 opUserID
	if values := md.Get(constant.OpUserID); len(values) > 0 {
		ctx = context.WithValue(ctx, constant.OpUserID, values[0])
	}

	// 提取 platform
	if values := md.Get(constant.OpUserPlatform); len(values) > 0 {
		ctx = context.WithValue(ctx, constant.OpUserPlatform, values[0])
	}

	// 提取 connID
	if values := md.Get(constant.ConnID); len(values) > 0 {
		ctx = context.WithValue(ctx, constant.ConnID, values[0])
	}

	return ctx
}
