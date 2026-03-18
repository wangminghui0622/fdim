package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Friend 好友模型
type Friend struct {
	ID             primitive.ObjectID `bson:"_id"`
	OwnerUserID    string             `bson:"owner_user_id"`
	FriendUserID   string             `bson:"friend_user_id"`
	Remark         string             `bson:"remark"`
	CreateTime     time.Time          `bson:"create_time"`
	AddSource      int32              `bson:"add_source"`
	OperatorUserID string             `bson:"operator_user_id"`
	Ex             string             `bson:"ex"`
	IsPinned       bool               `bson:"is_pinned"`
}

// FriendRequest 好友申请模型
type FriendRequest struct {
	FromUserID    string    `bson:"from_user_id"`
	ToUserID      string    `bson:"to_user_id"`
	HandleResult  int32     `bson:"handle_result"`
	ReqMsg        string    `bson:"req_msg"`
	CreateTime    time.Time `bson:"create_time"`
	HandlerUserID string    `bson:"handler_user_id"`
	HandleMsg     string    `bson:"handle_msg"`
	HandleTime    time.Time `bson:"handle_time"`
	Ex            string    `bson:"ex"`
}
