package ws

import (
	"encoding/json"

	"fdim/protocol/sdkws"
)

// msgDataJSON is a JSON-friendly representation of sdkws.MsgData
// where Content ([]byte in protobuf) is serialized as a UTF-8 string
// instead of base64, matching the official OpenIM client SDK expectation.
type offlinePushInfoJSON struct {
	Title        string `json:"title,omitempty"`
	Desc         string `json:"desc,omitempty"`
	Ex           string `json:"ex,omitempty"`
	IOSPushSound string `json:"iOSPushSound,omitempty"`
	IOSBadgeCount bool  `json:"iOSBadgeCount,omitempty"`
	SignalInfo   string `json:"signalInfo,omitempty"`
}

type msgDataJSON struct {
	SendID           string                 `json:"sendID,omitempty"`
	RecvID           string                 `json:"recvID,omitempty"`
	GroupID          string                 `json:"groupID,omitempty"`
	ClientMsgID      string                 `json:"clientMsgID,omitempty"`
	ServerMsgID      string                 `json:"serverMsgID,omitempty"`
	SenderPlatformID int32                  `json:"senderPlatformID,omitempty"`
	SenderNickname   string                 `json:"senderNickname,omitempty"`
	SenderFaceURL    string                 `json:"senderFaceURL,omitempty"`
	SessionType      int32                  `json:"sessionType,omitempty"`
	MsgFrom          int32                  `json:"msgFrom,omitempty"`
	ContentType      int32                  `json:"contentType,omitempty"`
	Content          string                 `json:"content,omitempty"`
	Seq              int64                  `json:"seq,omitempty"`
	SendTime         int64                  `json:"sendTime,omitempty"`
	CreateTime       int64                  `json:"createTime,omitempty"`
	Status           int32                  `json:"status,omitempty"`
	IsRead           bool                   `json:"isRead,omitempty"`
	Options          map[string]bool        `json:"options,omitempty"`
	OfflinePushInfo  *offlinePushInfoJSON   `json:"offlinePushInfo,omitempty"`
	AtUserIDList     []string               `json:"atUserIDList,omitempty"`
	AttachedInfo     string                 `json:"attachedInfo,omitempty"`
	Ex               string                 `json:"ex,omitempty"`
}

// MsgDataToJSON converts a protobuf MsgData to JSON bytes with Content as a
// UTF-8 string (not base64), suitable for sending over JSON WebSocket.
func MsgDataToJSON(m *sdkws.MsgData) (json.RawMessage, error) {
	if m == nil {
		return nil, nil
	}
	j := msgDataJSON{
		SendID:           m.SendID,
		RecvID:           m.RecvID,
		GroupID:          m.GroupID,
		ClientMsgID:      m.ClientMsgID,
		ServerMsgID:      m.ServerMsgID,
		SenderPlatformID: m.SenderPlatformID,
		SenderNickname:   m.SenderNickname,
		SenderFaceURL:    m.SenderFaceURL,
		SessionType:      m.SessionType,
		MsgFrom:          m.MsgFrom,
		ContentType:      m.ContentType,
		Content:          string(m.Content),
		Seq:              m.Seq,
		SendTime:         m.SendTime,
		CreateTime:       m.CreateTime,
		Status:           m.Status,
		IsRead:           m.IsRead,
		Options:          m.Options,
		AtUserIDList:     m.AtUserIDList,
		AttachedInfo:     m.AttachedInfo,
		Ex:               m.Ex,
	}
	if m.OfflinePushInfo != nil {
		j.OfflinePushInfo = &offlinePushInfoJSON{
			Title:         m.OfflinePushInfo.Title,
			Desc:          m.OfflinePushInfo.Desc,
			Ex:            m.OfflinePushInfo.Ex,
			IOSPushSound:  m.OfflinePushInfo.IOSPushSound,
			IOSBadgeCount: m.OfflinePushInfo.IOSBadgeCount,
			SignalInfo:    m.OfflinePushInfo.SignalInfo,
		}
	}
	return json.Marshal(j)
}
