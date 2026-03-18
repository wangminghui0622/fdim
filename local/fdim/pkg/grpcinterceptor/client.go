package grpcinterceptor

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/protocol/constant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GrpcClient 返回gRPC客户端拦截器选项，用于传递context中的用户信息到gRPC metadata
func GrpcClient() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(RpcClientInterceptor)
}

// RpcClientInterceptor gRPC客户端拦截器，将context中的用户信息传递到gRPC metadata
func RpcClientInterceptor(ctx context.Context, method string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	if ctx == nil {
		return errs.ErrInternalServer.WrapMsg("call rpc request context is nil")
	}
	ctx, err = getRpcContext(ctx)
	if err != nil {
		return err
	}
	return invoker(ctx, method, req, resp, cc, opts...)
}

// getRpcContext 从context中提取用户信息并设置到gRPC metadata
func getRpcContext(ctx context.Context) (context.Context, error) {
	md := metadata.Pairs()
	
	// 处理自定义header
	if keys, _ := ctx.Value(constant.RpcCustomHeader).([]string); len(keys) > 0 {
		for _, key := range keys {
			val, ok := ctx.Value(key).([]string)
			if !ok {
				return nil, errs.ErrInternalServer.WrapMsg(fmt.Sprintf("ctx missing key: %s", key))
			}
			if len(val) == 0 {
				return nil, errs.ErrInternalServer.WrapMsg(fmt.Sprintf("ctx key value is empty: %s", key))
			}
			md.Set(key, val...)
		}
		md.Set(constant.RpcCustomHeader, keys...)
	}
	
	// 设置operationID（可选，如果没有则生成默认值）
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if ok && operationID != "" {
		md.Set(constant.OperationID, operationID)
	} else {
		md.Set(constant.OperationID, "default-op-id")
	}
	
	// 设置opUserID
	opUserID, ok := ctx.Value(constant.OpUserID).(string)
	if ok && opUserID != "" {
		md.Set(constant.OpUserID, opUserID)
	}
	
	// 设置platform
	opUserPlatform, ok := ctx.Value(constant.OpUserPlatform).(string)
	if ok && opUserPlatform != "" {
		md.Set(constant.OpUserPlatform, opUserPlatform)
	}
	
	// 设置connID
	connID, ok := ctx.Value(constant.ConnID).(string)
	if ok && connID != "" {
		md.Set(constant.ConnID, connID)
	}
	
	return metadata.NewOutgoingContext(ctx, md), nil
}
