package types

// 添加收藏请求
type AddFavoriteReq struct {
	ClientMsgID    string `json:"clientMsgID"`
	ServerMsgID    string `json:"serverMsgID,optional"`
	ConversationID string `json:"conversationID,optional"`
	SendID         string `json:"sendID"`
	RecvID         string `json:"recvID,optional"`
	GroupID        string `json:"groupID,optional"`
	SenderNickname string `json:"senderNickname,optional"`
	SenderFaceURL  string `json:"senderFaceURL,optional"`
	SessionType    int32  `json:"sessionType"`
	ContentType    int32  `json:"contentType"`
	Content        string `json:"content"`
	SendTime       int64  `json:"sendTime"`
	Ex             string `json:"ex,optional"`
}

type AddFavoriteResp struct {
	BaseResp
	FavoriteID string `json:"favoriteID"`
}

// 删除收藏请求
type DeleteFavoriteReq struct {
	FavoriteID string `json:"favoriteID"`
}

type DeleteFavoriteResp struct {
	BaseResp
}

// 获取收藏列表请求
type GetFavoriteListReq struct {
	Pagination Pagination `json:"pagination"`
}

type GetFavoriteListResp struct {
	BaseResp
	Total     int64         `json:"total"`
	Favorites []FavoriteMsg `json:"favorites"`
}

// 收藏消息
type FavoriteMsg struct {
	FavoriteID     string `json:"favoriteID"`
	ClientMsgID    string `json:"clientMsgID"`
	ServerMsgID    string `json:"serverMsgID"`
	ConversationID string `json:"conversationID"`
	SendID         string `json:"sendID"`
	RecvID         string `json:"recvID"`
	GroupID        string `json:"groupID"`
	SenderNickname string `json:"senderNickname"`
	SenderFaceURL  string `json:"senderFaceURL"`
	SessionType    int32  `json:"sessionType"`
	ContentType    int32  `json:"contentType"`
	Content        string `json:"content"`
	SendTime       int64  `json:"sendTime"`
	CreateTime     int64  `json:"createTime"`
	Ex             string `json:"ex"`
}
