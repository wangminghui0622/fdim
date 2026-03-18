package model

import "time"

// Group 群组模型
type Group struct {
	GroupID                string    `bson:"group_id"`
	GroupName              string    `bson:"group_name"`
	Notification           string    `bson:"notification"`
	Introduction           string    `bson:"introduction"`
	FaceURL                string    `bson:"face_url"`
	CreateTime             time.Time `bson:"create_time"`
	Ex                     string    `bson:"ex"`
	Status                 int32     `bson:"status"`
	CreatorUserID          string    `bson:"creator_user_id"`
	GroupType              int32     `bson:"group_type"`
	NeedVerification       int32     `bson:"need_verification"`
	LookMemberInfo         int32     `bson:"look_member_info"`
	ApplyMemberFriend      int32     `bson:"apply_member_friend"`
	NotificationUpdateTime time.Time `bson:"notification_update_time"`
	NotificationUserID     string    `bson:"notification_user_id"`
}

// GroupMember 群成员模�?
type GroupMember struct {
	GroupID        string    `bson:"group_id"`
	UserID         string    `bson:"user_id"`
	Nickname       string    `bson:"nickname"`
	FaceURL        string    `bson:"face_url"`
	RoleLevel      int32     `bson:"role_level"`
	JoinTime       time.Time `bson:"join_time"`
	JoinSource     int32     `bson:"join_source"`
	InviterUserID  string    `bson:"inviter_user_id"`
	OperatorUserID string    `bson:"operator_user_id"`
	MuteEndTime    time.Time `bson:"mute_end_time"`
	Ex             string    `bson:"ex"`
}

// GroupRequest 群组申请模型
type GroupRequest struct {
	UserID        string    `bson:"user_id"`
	GroupID       string    `bson:"group_id"`
	HandleResult  int32     `bson:"handle_result"`
	ReqMsg        string    `bson:"req_msg"`
	HandledMsg    string    `bson:"handled_msg"`
	ReqTime       time.Time `bson:"req_time"`
	HandleUserID  string    `bson:"handle_user_id"`
	HandledTime   time.Time `bson:"handled_time"`
	JoinSource    int32     `bson:"join_source"`
	InviterUserID string    `bson:"inviter_user_id"`
	Ex            string    `bson:"ex"`
}
