package mw

import (
	"context"
	"fmt"

	"fdim/protocol/constant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// RpcServerInterceptor gRPC服务端拦截器，从metadata中提取用户信息到context
func RpcServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, err := validateMetadata(ctx)
	if err != nil {
		return nil, err
	}
	ctx, err = enrichContextWithMetadata(ctx, md)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func validateMetadata(ctx context.Context) (metadata.MD, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.New(codes.InvalidArgument, "missing metadata").Err()
	}
	return md, nil
}

func enrichContextWithMetadata(ctx context.Context, md metadata.MD) (context.Context, error) {
	if keys := md.Get(constant.RpcCustomHeader); len(keys) > 0 {
		ctx = context.WithValue(ctx, constant.RpcCustomHeader, keys)
		for _, key := range keys {
			values := md.Get(key)
			if len(values) == 0 {
				return nil, status.New(codes.InvalidArgument, fmt.Sprintf("missing metadata key %s", key)).Err()
			}
			ctx = context.WithValue(ctx, key, values)
		}
	}
	if opts := md.Get(constant.OperationID); len(opts) == 1 {
		ctx = context.WithValue(ctx, constant.OperationID, opts[0])
	}
	if opts := md.Get(constant.OpUserID); len(opts) == 1 {
		ctx = context.WithValue(ctx, constant.OpUserID, opts[0])
	}
	if opts := md.Get(constant.OpUserPlatform); len(opts) == 1 {
		ctx = context.WithValue(ctx, constant.OpUserPlatform, opts[0])
	}
	if opts := md.Get(constant.ConnID); len(opts) == 1 {
		ctx = context.WithValue(ctx, constant.ConnID, opts[0])
	}
	return ctx, nil
}

func GrpcServer() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(RpcServerInterceptor)
}
