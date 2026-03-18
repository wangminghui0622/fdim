package convert

import (
	"time"

	"fdim/pkg/model"
	"fdim/protocol/sdkws"
)

// GroupRequestDB2Pb 将数据库群组申请模型转换为 Protocol Buffer 格式
// 注意：sdkws.GroupRequest 需要 UserInfo 和 GroupInfo 对象，而不是直接字段
func GroupRequestDB2Pb(req *model.GroupRequest, userInfo *sdkws.PublicUserInfo, groupInfo *sdkws.GroupInfo) *sdkws.GroupRequest {
	return &sdkws.GroupRequest{
		UserInfo:      userInfo,
		GroupInfo:     groupInfo,
		HandleResult:  req.HandleResult,
		ReqMsg:        req.ReqMsg,
		HandleMsg:     req.HandledMsg,
		ReqTime:       req.ReqTime.UnixMilli(),
		HandleUserID:  req.HandleUserID,
		HandleTime:    req.HandledTime.UnixMilli(),
		Ex:            req.Ex,
		JoinSource:    req.JoinSource,
		InviterUserID: req.InviterUserID,
	}
}

// GroupRequestsDB2Pb 将数据库群组申请模型列表转换为 Protocol Buffer 格式列表
func GroupRequestsDB2Pb(requests []*model.GroupRequest, getUserInfo func(string) *sdkws.PublicUserInfo, getGroupInfo func(string) *sdkws.GroupInfo) []*sdkws.GroupRequest {
	result := make([]*sdkws.GroupRequest, 0, len(requests))
	for _, r := range requests {
		var userInfo *sdkws.PublicUserInfo
		var groupInfo *sdkws.GroupInfo
		if getUserInfo != nil {
			userInfo = getUserInfo(r.UserID)
		}
		if getGroupInfo != nil {
			groupInfo = getGroupInfo(r.GroupID)
		}
		result = append(result, GroupRequestDB2Pb(r, userInfo, groupInfo))
	}
	return result
}

// GroupRequestPb2DB 将 Protocol Buffer 群组申请转换为数据库模型
func GroupRequestPb2DB(req *sdkws.GroupRequest) *model.GroupRequest {
	var userID string
	var groupID string
	if req.UserInfo != nil {
		userID = req.UserInfo.UserID
	}
	if req.GroupInfo != nil {
		groupID = req.GroupInfo.GroupID
	}
	return &model.GroupRequest{
		UserID:        userID,
		GroupID:       groupID,
		HandleResult:  req.HandleResult,
		ReqMsg:        req.ReqMsg,
		HandledMsg:    req.HandleMsg,
		ReqTime:       time.UnixMilli(req.ReqTime),
		HandleUserID:  req.HandleUserID,
		HandledTime:   time.UnixMilli(req.HandleTime),
		JoinSource:    req.JoinSource,
		InviterUserID: req.InviterUserID,
		Ex:            req.Ex,
	}
}
