package nats

import (
	"context"

	"fdim/pkg/constant"
	"fdim/pkg/mcontext"
	"github.com/nats-io/nats.go"
)

// GetMQHeaderWithContext extracts message queue headers from the context.
func GetMQHeaderWithContext(ctx context.Context) (nats.Header, error) {
	operationID, opUserID, platform, connID, err := mcontext.GetCtxInfos(ctx)
	if err != nil {
		return nil, err
	}
	header := make(nats.Header)
	header.Set(constant.OperationID, operationID)
	header.Set(constant.OpUserID, opUserID)
	header.Set(constant.OpUserPlatform, platform)
	header.Set(constant.ConnID, connID)
	return header, nil
}

// GetContextWithMQHeader creates a context from message queue headers.
func GetContextWithMQHeader(header nats.Header) context.Context {
	var values []string
	if header != nil {
		values = append(values, header.Get(constant.OperationID))
		values = append(values, header.Get(constant.OpUserID))
		values = append(values, header.Get(constant.OpUserPlatform))
		values = append(values, header.Get(constant.ConnID))
	}
	return mcontext.WithMustInfoCtx(values)
}
