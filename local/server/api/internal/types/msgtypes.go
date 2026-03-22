package types

// Msg 相关请求和响应类型
// 注意：这些类型需要与 open-im-server 的 API 接口保持一致

// GetSeq (newest_seq) 相关类型
type GetSeqReq struct {
	UserID         string `json:"userID"`
	ConversationID string `json:"conversationID,optional"`
}

type GetSeqResp struct {
	BaseResp
	MaxSeq int64 `json:"maxSeq"`
}

// SearchMsg 相关类型
type SearchMsgReq struct {
	SendID      string     `json:"sendID,optional"`
	RecvID      string     `json:"recvID,optional"`
	ContentType []int32    `json:"contentType,optional"`
	SendTime    int64      `json:"sendTime,optional"`
	SessionType int32      `json:"sessionType,optional"`
	Pagination  Pagination `json:"pagination"`
}

type SearchMsgResp struct {
	BaseResp
	Total int64         `json:"total"`
	Msgs  []interface{} `json:"msgs"` // 使用 sdkws.MsgData
}

// SendMsg 相关类型
type SendMsgReq struct {
	RecvID  string     `json:"recvID"`
	SendMsg *MsgData   `json:"sendMsg"` // 使用 MsgData 结构
}

// MsgData 消息数据结构
type MsgData struct {
	SendID           string                 `json:"sendID"`
	RecvID           string                 `json:"recvID"`
	GroupID          string                 `json:"groupID"`
	ClientMsgID      string                 `json:"clientMsgID"`
	ServerMsgID      string                 `json:"serverMsgID,optional"`
	SenderPlatformID int32                  `json:"senderPlatformID"`
	SenderNickname   string                 `json:"senderNickname"`
	SenderFaceURL    string                 `json:"senderFaceURL"`
	SessionType      int32                  `json:"sessionType"`
	MsgFrom          int32                  `json:"msgFrom,optional"`
	ContentType      int32                  `json:"contentType"`
	Content          string                 `json:"content"`
	Seq              int64                  `json:"seq,optional"`
	SendTime         int64                  `json:"sendTime,optional"`
	CreateTime       int64                  `json:"createTime,optional"`
	Status           int32                  `json:"status,optional"`
	IsRead           bool                   `json:"isRead,optional"`
	Options          map[string]bool        `json:"options,optional"`
	OfflinePushInfo  map[string]interface{} `json:"offlinePushInfo,optional"`
	AtUserIDList     []string               `json:"atUserIDList,optional"`
	AttachedInfo     string                 `json:"attachedInfo,optional"`
	Ex               string                 `json:"ex,optional"`
	SenderTimeZone   int32                  `json:"senderTimeZone,optional"` // 发送者时区偏移（秒），如 UTC+8 = 28800
}

type SendMsgResp struct {
	BaseResp
	ServerMsgID string                 `json:"serverMsgID"`
	ClientMsgID string                 `json:"clientMsgID"`
	SendTime    int64                  `json:"sendTime"`
	Modify      map[string]interface{} `json:"modify,optional"`
}

// SendBusinessNotification 相关类型
type SendBusinessNotificationReq struct {
	Key              string `json:"key"`
	Data             string `json:"data"`
	SendUserID       string `json:"sendUserID"`
	RecvUserID       string `json:"recvUserID,optional"`
	RecvGroupID      string `json:"recvGroupID,optional"`
	SendMsg          bool   `json:"sendMsg"`
	ReliabilityLevel *int   `json:"reliabilityLevel,optional"`
}

type SendBusinessNotificationResp struct {
	BaseResp
	ServerMsgID string `json:"serverMsgID"`
	ClientMsgID string `json:"clientMsgID"`
	SendTime    int64  `json:"sendTime"`
}

// PullMsgBySeqs 相关类型
type PullMsgBySeqsReq struct {
	UserID         string   `json:"userID"`
	ConversationID string   `json:"conversationID,optional"`
	Seqs           []uint32 `json:"seqs"`
	GroupSeqs      []uint32 `json:"groupSeqs,optional"`
}

type PullMsgBySeqsResp struct {
	BaseResp
	Msgs []interface{} `json:"msgs"` // 使用 sdkws.MsgData
}

// RevokeMsg 相关类型
type RevokeMsgReq struct {
	ConversationID string `json:"conversationID"`
	Seq            int64  `json:"seq"`
	UserID         string `json:"userID"`
}

type RevokeMsgResp struct {
	BaseResp
}

// MarkMsgsAsRead 相关类型
type MarkMsgsAsReadReq struct {
	ConversationID string   `json:"conversationID"`
	Seqs           []int64  `json:"seqs"`
	UserID         string   `json:"userID"`
}

type MarkMsgsAsReadResp struct {
	BaseResp
}

// MarkConversationAsRead 相关类型
type MarkConversationAsReadReq struct {
	ConversationID string `json:"conversationID"`
	UserID         string `json:"userID"`
	HasReadSeq     int64  `json:"hasReadSeq,optional"`
	Seqs           []int64 `json:"seqs,optional"`
}

type MarkConversationAsReadResp struct {
	BaseResp
}

// GetConversationsHasReadAndMaxSeq 相关类型
type GetConversationsHasReadAndMaxSeqReq struct {
	UserID          string   `json:"userID"`
	ConversationIDs []string `json:"conversationIDs"`
}

type GetConversationsHasReadAndMaxSeqResp struct {
	BaseResp
	Seqs map[string]interface{} `json:"seqs"` // 使用 SeqsInfoResp
}

// SetConversationHasReadSeq 相关类型
type SetConversationHasReadSeqReq struct {
	UserID         string `json:"userID"`
	ConversationID string `json:"conversationID"`
	HasReadSeq     int64  `json:"hasReadSeq"`
}

type SetConversationHasReadSeqResp struct {
	BaseResp
}

// ClearConversationsMsg 相关类型
type ClearConversationsMsgReq struct {
	ConversationIDs []string `json:"conversationIDs"`
	UserID          string   `json:"userID"`
}

type ClearConversationsMsgResp struct {
	BaseResp
}

// UserClearAllMsg 相关类型
type UserClearAllMsgReq struct {
	UserID string `json:"userID"`
}

type UserClearAllMsgResp struct {
	BaseResp
}

// DeleteMsgs 相关类型
type DeleteMsgsReq struct {
	ConversationID string   `json:"conversationID"`
	Seqs           []int64  `json:"seqs"`
	UserID         string   `json:"userID"`
}

type DeleteMsgsResp struct {
	BaseResp
}

// DeleteMsgPhysicalBySeq 相关类型
type DeleteMsgPhysicalBySeqReq struct {
	ConversationID string `json:"conversationID"`
	Seq            int64  `json:"seq"`
}

type DeleteMsgPhysicalBySeqResp struct {
	BaseResp
}

// DeleteMsgPhysical 相关类型
type DeleteMsgPhysicalReq struct {
	ConversationID string `json:"conversationID"`
}

type DeleteMsgPhysicalResp struct {
	BaseResp
}

// BatchSendMsg 相关类型
type BatchSendMsgReq struct {
	SendMsg   *MsgData `json:"sendMsg"` // 使用 MsgData 结构
	RecvIDs   []string `json:"recvIDs,optional"`
	IsSendAll bool     `json:"isSendAll,optional"`
}

type BatchSendMsgResp struct {
	BaseResp
	Results  []interface{} `json:"results"`
	FailedIDs []string       `json:"failedIDs"`
}

// SendSimpleMsg 相关类型
type SendSimpleMsgReq struct {
	SendID         string `json:"sendID"`
	Content        string `json:"content"`
	OfflinePushInfo interface{} `json:"offlinePushInfo,optional"`
	Ex             string `json:"ex,optional"`
}

type SendSimpleMsgResp struct {
	BaseResp
	ServerMsgID string `json:"serverMsgID"`
	ClientMsgID string `json:"clientMsgID"`
	SendTime    int64  `json:"sendTime"`
}

// CheckMsgIsSendSuccess 相关类型
type CheckMsgIsSendSuccessReq struct {
	UserID string `json:"userID"`
}

type CheckMsgIsSendSuccessResp struct {
	BaseResp
	Status int32 `json:"status"`
}

// GetServerTime 相关类型
type GetServerTimeReq struct {
}

type GetServerTimeResp struct {
	BaseResp
	ServerTime int64 `json:"serverTime"`
}
