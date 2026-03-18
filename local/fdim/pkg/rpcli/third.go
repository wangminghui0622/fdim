package rpcli

import (
	"fdim/protocol/third"
	"google.golang.org/grpc"
)

func NewThirdClient(cc grpc.ClientConnInterface) *ThirdClient {
	return &ThirdClient{third.NewThirdClient(cc)}
}

type ThirdClient struct {
	third.ThirdClient
}
