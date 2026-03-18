package model

import (
	"time"
)

// FavoriteMsg 消息收藏模型
type FavoriteMsg struct {
	ID             string    `bson:"_id"`
	UserID         string    `bson:"user_id"`
	ClientMsgID    string    `bson:"client_msg_id"`
	ServerMsgID    string    `bson:"server_msg_id"`
	ConversationID string    `bson:"conversation_id"`
	SendID         string    `bson:"send_id"`
	RecvID         string    `bson:"recv_id"`
	GroupID        string    `bson:"group_id"`
	SenderNickname string    `bson:"sender_nickname"`
	SenderFaceURL  string    `bson:"sender_face_url"`
	SessionType    int32     `bson:"session_type"`
	ContentType    int32     `bson:"content_type"`
	Content        string    `bson:"content"`
	SendTime       int64     `bson:"send_time"`
	CreateTime     time.Time `bson:"create_time"`
	Ex             string    `bson:"ex"`
}

// FavoriteMsgModel 收藏消息数据库操作接口
type FavoriteMsgModel interface {
	Create(msg *FavoriteMsg) error
	Delete(userID, favoriteID string) error
	GetByUserID(userID string, pageNumber, showNumber int32) ([]*FavoriteMsg, int64, error)
	GetByID(userID, favoriteID string) (*FavoriteMsg, error)
	DeleteByClientMsgID(userID, clientMsgID string) error
}
