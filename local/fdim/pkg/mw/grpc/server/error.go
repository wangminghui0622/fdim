package server

import (
	"context"
	"errors"

	"fdim/pkg/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcServerErrorConvert() grpc.ServerOption {
	type grpcError interface {
		error
		GRPCStatus() *status.Status
	}
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return
		}
		var grpcErr grpcError
		if errors.As(err, &grpcErr) {
			return
		}
		err = codeErrorToGrpcError(ctx, getCodeError(err))
		return
	})
}

func getCodeError(err error) errs.CodeError {
	return errs.ErrInternalServer.WithDetail(err.Error())
}

func codeErrorToGrpcError(ctx context.Context, codeErr errs.CodeError) error {
	grpcStatus := status.New(codes.Code(codeErr.Code()), codeErr.Msg())
	return grpcStatus.Err()
}
