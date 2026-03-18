package ws

import "encoding/json"

// WebSocket protocol constants matching official OpenIM v3

// Client → Server request identifiers
const (
	WSGetNewestSeq        = 1001
	WSSendMsg             = 1002
	WSPullMsgBySeq        = 1003
	WSSendSignalMsg       = 1004
	WSSendBinaryMsg       = 1005
	WSSetBackgroundStatus = 1006
)

// Server → Client push identifiers
const (
	WSPushMsg       = 2001
	WSKickOnlineMsg = 2002
	WsLogoutMsg     = 2003
)

// WsReq is the WebSocket request frame (client → server)
type WsReq struct {
	ReqIdentifier int32            `json:"reqIdentifier"`
	Token         string           `json:"token"`
	SendID        string           `json:"sendID"`
	OperationID   string           `json:"operationID"`
	MsgIncr       string           `json:"msgIncr"`
	Data          json.RawMessage  `json:"data,omitempty"`
}

// WsResp is the WebSocket response frame (server → client)
type WsResp struct {
	ReqIdentifier int32            `json:"reqIdentifier"`
	ErrCode       int32            `json:"errCode"`
	ErrMsg        string           `json:"errMsg"`
	OperationID   string           `json:"operationID"`
	MsgIncr       string           `json:"msgIncr"`
	Data          json.RawMessage  `json:"data,omitempty"`
}
