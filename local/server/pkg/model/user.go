package model

import "time"

// User ûģ
type User struct {
	UserID           string    `bson:"user_id"`
	Nickname         string    `bson:"nickname"`
	FaceURL          string    `bson:"face_url"`
	Ex               string    `bson:"ex"`
	AppMangerLevel   int32     `bson:"app_manger_level"`
	GlobalRecvMsgOpt int32     `bson:"global_recv_msg_opt"`
	CreateTime       time.Time `bson:"create_time"`
}
