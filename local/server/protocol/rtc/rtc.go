package rtc

import "google.golang.org/grpc"

// RtcServiceClient is the client API for RtcService service.
type RtcServiceClient interface {
}

type rtcServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRtcServiceClient(cc grpc.ClientConnInterface) RtcServiceClient {
	return &rtcServiceClient{cc}
}
