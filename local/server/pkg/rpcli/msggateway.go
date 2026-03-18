package rpcli

import (
	"fdim/protocol/msggateway"
	"google.golang.org/grpc"
)

func NewMsgGatewayClient(cc grpc.ClientConnInterface) *MsgGatewayClient {
	return &MsgGatewayClient{msggateway.NewMsgGatewayClient(cc)}
}

type MsgGatewayClient struct {
	msggateway.MsgGatewayClient
}
