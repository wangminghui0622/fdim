package rpcli

import (
	"fdim/protocol/rtc"
	"google.golang.org/grpc"
)

func NewRtcServiceClient(cc grpc.ClientConnInterface) *RtcServiceClient {
	return &RtcServiceClient{rtc.NewRtcServiceClient(cc)}
}

type RtcServiceClient struct {
	rtc.RtcServiceClient
}
