package types

// Conversation 相关请求和响应类型
// 注意：这些类型需要与 open-im-server 的 API 接口保持一致

// GetAllConversations 相关类型
type GetAllConversationsReq struct {
	OwnerUserID string `json:"ownerUserID"`
}

type GetAllConversationsResp struct {
	BaseResp
	Conversations []interface{} `json:"conversations"` // 使用 sdkws.Conversation
}

// GetSortedConversationList 相关类型
type GetSortedConversationListReq struct {
	UserID         string     `json:"userID,optional"`
	ConversationIDs []string   `json:"conversationIDs,optional"`
	Pagination     Pagination `json:"pagination,optional"`
}

type ConversationElem struct {
	ConversationID string `json:"conversationID"`
	RecvMsgOpt     int32  `json:"recvMsgOpt"`
	UnreadCount    int64  `json:"unreadCount"`
	IsPinned       bool   `json:"isPinned"`
}

type GetSortedConversationListResp struct {
	BaseResp
	ConversationTotal int64           `json:"conversationTotal"`
	UnreadTotal       int64           `json:"unreadTotal"`
	ConversationElems []ConversationElem `json:"conversationElems"`
}

// GetConversation 相关类型
type GetConversationReq struct {
	OwnerUserID     string `json:"ownerUserID"`
	ConversationID  string `json:"conversationID"`
}

type GetConversationResp struct {
	BaseResp
	Conversation interface{} `json:"conversation"` // 使用 sdkws.Conversation
}

// GetConversations 相关类型
type GetConversationsReq struct {
	OwnerUserID    string   `json:"ownerUserID"`
	ConversationIDs []string `json:"conversationIDs"`
}

type GetConversationsResp struct {
	BaseResp
	Conversations []interface{} `json:"conversations"` // 使用 sdkws.Conversation
}

// SetConversations 相关类型
type SetConversationsReq struct {
	OwnerUserID   string        `json:"ownerUserID"`
	Conversations []interface{} `json:"conversations"` // 使用 sdkws.Conversation
}

type SetConversationsResp struct {
	BaseResp
}

// GetFullOwnerConversationIDs 相关类型
type GetFullOwnerConversationIDsReq struct {
	IdHash uint64 `json:"idHash"`
	UserID string `json:"userID"`
}

type GetFullOwnerConversationIDsResp struct {
	BaseResp
	ConversationIDs []string `json:"conversationIDs"`
}

// GetIncrementalConversation 相关类型
type GetIncrementalConversationReq struct {
	UserID    string `json:"userID"`
	VersionID string `json:"versionID,optional"`
	Version   uint64 `json:"version"`
}

type GetIncrementalConversationResp struct {
	BaseResp
	Version        uint64        `json:"version"`
	VersionID      string        `json:"versionID"`
	Full           bool          `json:"full"`
	Delete         []string      `json:"delete"`
	Insert         []interface{} `json:"insert"` // 使用 sdkws.Conversation
	Update         []interface{} `json:"update"` // 使用 sdkws.Conversation
}

// GetOwnerConversation 相关类型
type GetOwnerConversationReq struct {
	UserID     string     `json:"userID"`
	Pagination Pagination `json:"pagination,optional"`
}

type GetOwnerConversationResp struct {
	BaseResp
	Total         int64         `json:"total"`
	Conversations []interface{} `json:"conversations"` // 使用 conversation.Conversation
}

// GetNotNotifyConversationIDs 相关类型
type GetNotNotifyConversationIDsReq struct {
	UserID string `json:"userID"`
}

type GetNotNotifyConversationIDsResp struct {
	BaseResp
	ConversationIDs []string `json:"conversationIDs"`
}

// GetPinnedConversationIDs 相关类型
type GetPinnedConversationIDsReq struct {
	UserID string `json:"userID"`
}

type GetPinnedConversationIDsResp struct {
	BaseResp
	ConversationIDs []string `json:"conversationIDs"`
}

// UpdateConversationsByUser 相关类型
type UpdateConversationsByUserReq struct {
	UserID        string  `json:"userID"`
	Ex            *string `json:"ex,optional"`
}

type UpdateConversationsByUserResp struct {
	BaseResp
}

// DeleteConversations 相关类型
type DeleteConversationsReq struct {
	OwnerUserID     string   `json:"ownerUserID"`
	NeedDeleteTime  int64    `json:"needDeleteTime"`
	ConversationIDs []string `json:"conversationIDs,optional"`
}

type DeleteConversationsResp struct {
	BaseResp
}
