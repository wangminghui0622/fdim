package model

// MsgDoc 消息文档模型（对?MongoDB ?stream_msg 集合?
// 参?FDIM 官方实现，时间字段使?int64 存储毫秒时间?
type MsgDoc struct {
	ConversationID   string                 `bson:"conversation_id"`
	Seq              int64                  `bson:"seq"`
	SendID           string                 `bson:"send_id"`
	RecvID           string                 `bson:"recv_id"`
	GroupID          string                 `bson:"group_id"`
	ClientMsgID      string                 `bson:"client_msg_id"`
	ServerMsgID      string                 `bson:"server_msg_id"`
	SenderPlatformID int32                  `bson:"sender_platform_id"`
	SenderNickname   string                 `bson:"sender_nickname"`
	SenderFaceURL    string                 `bson:"sender_face_url"`
	SessionType      int32                  `bson:"session_type"`
	MsgFrom          int32                  `bson:"msg_from"`
	ContentType      int32                  `bson:"content_type"`
	Content          []byte                 `bson:"content"`
	CreateTime       int64                  `bson:"create_time"`     // 毫秒时间?
	SendTime         int64                  `bson:"send_time"`       // 毫秒时间?
	Status           int32                  `bson:"status"`
	Options          map[string]bool        `bson:"options"`
	OfflinePush      map[string]interface{} `bson:"offline_push"`
	AtUserIDs        []string               `bson:"at_user_ids"`
	AttachedInfo     string                 `bson:"attached_info"`
	Ex               string                 `bson:"ex"`
	// 阅后即焚相关字段
	BurnDuration int32 `bson:"burn_duration,omitempty"` // 阅后即焚时长（秒），0表示不启?
	IsRead       bool  `bson:"is_read"`                 // 是否已读
	ReadTime     int64 `bson:"read_time,omitempty"`     // 阅读时间（毫秒时间戳?表示未读?
	BurnTime     int64 `bson:"burn_time,omitempty"`     // 销毁时间（毫秒时间戳，0表示未设置）
	IsEncrypted  bool  `bson:"is_encrypted"`            // 是否加密
}
