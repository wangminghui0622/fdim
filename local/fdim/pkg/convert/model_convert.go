// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package convert

import (
	"fdim/pkg/model"
	storagemodel "fdim/pkg/storage/model"
	"fdim/pkg/util/datautil"
	"fdim/protocol/sdkws"
)

// ModelUserDB2Pb converts pkg/model.User to sdkws.UserInfo
func ModelUserDB2Pb(user *model.User) *sdkws.UserInfo {
	if user == nil {
		return nil
	}
	return &sdkws.UserInfo{
		UserID:           user.UserID,
		Nickname:         user.Nickname,
		FaceURL:          user.FaceURL,
		Ex:               user.Ex,
		CreateTime:       user.CreateTime.UnixMilli(),
		AppMangerLevel:   user.AppMangerLevel,
		GlobalRecvMsgOpt: user.GlobalRecvMsgOpt,
	}
}

// ModelUsersDB2Pb converts []*model.User to []*sdkws.UserInfo
func ModelUsersDB2Pb(users []*model.User) []*sdkws.UserInfo {
	return datautil.Slice(users, ModelUserDB2Pb)
}

// ModelGroupDB2Pb converts pkg/model.Group to sdkws.GroupInfo
func ModelGroupDB2Pb(group *model.Group) *sdkws.GroupInfo {
	if group == nil {
		return nil
	}
	return &sdkws.GroupInfo{
		GroupID:                group.GroupID,
		GroupName:              group.GroupName,
		Notification:           group.Notification,
		Introduction:           group.Introduction,
		FaceURL:                group.FaceURL,
		CreateTime:             group.CreateTime.UnixMilli(),
		Ex:                     group.Ex,
		Status:                 group.Status,
		CreatorUserID:          group.CreatorUserID,
		GroupType:              group.GroupType,
		NeedVerification:       group.NeedVerification,
		LookMemberInfo:         group.LookMemberInfo,
		ApplyMemberFriend:      group.ApplyMemberFriend,
		NotificationUpdateTime: group.NotificationUpdateTime.UnixMilli(),
		NotificationUserID:     group.NotificationUserID,
	}
}

// ModelGroupMemberDB2Pb converts pkg/model.GroupMember to sdkws.GroupMemberFullInfo
func ModelGroupMemberDB2Pb(member *model.GroupMember) *sdkws.GroupMemberFullInfo {
	if member == nil {
		return nil
	}
	return &sdkws.GroupMemberFullInfo{
		GroupID:        member.GroupID,
		UserID:         member.UserID,
		RoleLevel:      member.RoleLevel,
		JoinTime:       member.JoinTime.UnixMilli(),
		Nickname:       member.Nickname,
		FaceURL:        member.FaceURL,
		JoinSource:     member.JoinSource,
		OperatorUserID: member.OperatorUserID,
		Ex:             member.Ex,
		MuteEndTime:    member.MuteEndTime.UnixMilli(),
		InviterUserID:  member.InviterUserID,
	}
}

// ModelFriendRequestDB2Pb converts pkg/model.FriendRequest to sdkws.FriendRequest
func ModelFriendRequestDB2Pb(req *model.FriendRequest) *sdkws.FriendRequest {
	if req == nil {
		return nil
	}
	return &sdkws.FriendRequest{
		FromUserID:    req.FromUserID,
		ToUserID:      req.ToUserID,
		HandleResult:  req.HandleResult,
		ReqMsg:        req.ReqMsg,
		CreateTime:    req.CreateTime.UnixMilli(),
		HandlerUserID: req.HandlerUserID,
		HandleMsg:     req.HandleMsg,
		HandleTime:    req.HandleTime.UnixMilli(),
		Ex:            req.Ex,
	}
}

// ModelFriendDB2Pb converts pkg/model.Friend to sdkws.FriendInfo
func ModelFriendDB2Pb(friend *model.Friend, user *model.User) *sdkws.FriendInfo {
	if friend == nil {
		return nil
	}
	info := &sdkws.FriendInfo{
		OwnerUserID:    friend.OwnerUserID,
		Remark:         friend.Remark,
		CreateTime:     friend.CreateTime.UnixMilli(),
		AddSource:      friend.AddSource,
		OperatorUserID: friend.OperatorUserID,
		Ex:             friend.Ex,
		IsPinned:       friend.IsPinned,
	}
	if user != nil {
		info.FriendUser = ModelUserDB2Pb(user)
	}
	return info
}

// ModelBlackDB2Pb converts pkg/model.Black to sdkws.BlackInfo
func ModelBlackDB2Pb(black *model.Black, user *model.User) *sdkws.BlackInfo {
	if black == nil {
		return nil
	}
	info := &sdkws.BlackInfo{
		OwnerUserID:    black.OwnerUserID,
		CreateTime:     black.CreateTime.UnixMilli(),
		AddSource:      black.AddSource,
		OperatorUserID: black.OperatorUserID,
		Ex:             black.Ex,
	}
	if user != nil {
		info.BlackUserInfo = &sdkws.PublicUserInfo{
			UserID:   user.UserID,
			Nickname: user.Nickname,
			FaceURL:  user.FaceURL,
			Ex:       user.Ex,
		}
	}
	return info
}

// ModelToStorageUser converts pkg/model.User to pkg/storage/model.User
func ModelToStorageUser(u *model.User) *storagemodel.User {
	if u == nil {
		return nil
	}
	return &storagemodel.User{
		UserID:           u.UserID,
		Nickname:         u.Nickname,
		FaceURL:          u.FaceURL,
		Ex:               u.Ex,
		CreateTime:       u.CreateTime,
		AppMangerLevel:   u.AppMangerLevel,
		GlobalRecvMsgOpt: u.GlobalRecvMsgOpt,
	}
}

// ModelToStorageGroup converts pkg/model.Group to pkg/storage/model.Group
func ModelToStorageGroup(g *model.Group) *storagemodel.Group {
	if g == nil {
		return nil
	}
	return &storagemodel.Group{
		GroupID:                g.GroupID,
		GroupName:              g.GroupName,
		Notification:           g.Notification,
		Introduction:           g.Introduction,
		FaceURL:                g.FaceURL,
		CreateTime:             g.CreateTime,
		Ex:                     g.Ex,
		Status:                 g.Status,
		CreatorUserID:          g.CreatorUserID,
		GroupType:              g.GroupType,
		NeedVerification:       g.NeedVerification,
		LookMemberInfo:         g.LookMemberInfo,
		ApplyMemberFriend:      g.ApplyMemberFriend,
		NotificationUpdateTime: g.NotificationUpdateTime,
		NotificationUserID:     g.NotificationUserID,
	}
}

// ModelToStorageGroupMember converts pkg/model.GroupMember to pkg/storage/model.GroupMember
func ModelToStorageGroupMember(m *model.GroupMember) *storagemodel.GroupMember {
	if m == nil {
		return nil
	}
	return &storagemodel.GroupMember{
		GroupID:        m.GroupID,
		UserID:         m.UserID,
		Nickname:       m.Nickname,
		FaceURL:        m.FaceURL,
		RoleLevel:      m.RoleLevel,
		JoinTime:       m.JoinTime,
		JoinSource:     m.JoinSource,
		InviterUserID:  m.InviterUserID,
		OperatorUserID: m.OperatorUserID,
		MuteEndTime:    m.MuteEndTime,
		Ex:             m.Ex,
	}
}

// ModelGroupMembersDB2Pb converts []*model.GroupMember to []*sdkws.GroupMemberFullInfo
func ModelGroupMembersDB2Pb(members []*model.GroupMember) []*sdkws.GroupMemberFullInfo {
	return datautil.Slice(members, ModelGroupMemberDB2Pb)
}

// ModelFriendRequestsDB2Pb converts []*model.FriendRequest to []*sdkws.FriendRequest
func ModelFriendRequestsDB2Pb(requests []*model.FriendRequest) []*sdkws.FriendRequest {
	return datautil.Slice(requests, ModelFriendRequestDB2Pb)
}

// ModelFriendsDB2Pb converts []*model.Friend to []*sdkws.FriendInfo (without user info)
func ModelFriendsDB2Pb(friends []*model.Friend) []*sdkws.FriendInfo {
	result := make([]*sdkws.FriendInfo, 0, len(friends))
	for _, f := range friends {
		result = append(result, ModelFriendDB2Pb(f, nil))
	}
	return result
}

// ModelBlacksDB2Pb converts []*model.Black to []*sdkws.BlackInfo (without user info)
func ModelBlacksDB2Pb(blacks []*model.Black) []*sdkws.BlackInfo {
	result := make([]*sdkws.BlackInfo, 0, len(blacks))
	for _, b := range blacks {
		result = append(result, ModelBlackDB2Pb(b, nil))
	}
	return result
}

// ModelGroupsDB2Pb converts []*model.Group to []*sdkws.GroupInfo
func ModelGroupsDB2Pb(groups []*model.Group) []*sdkws.GroupInfo {
	return datautil.Slice(groups, ModelGroupDB2Pb)
}

// Pb2ModelGroupInfo converts sdkws.GroupInfo to pkg/model.Group
func Pb2ModelGroupInfo(m *sdkws.GroupInfo) *model.Group {
	if m == nil {
		return nil
	}
	return &model.Group{
		GroupID:                m.GroupID,
		GroupName:              m.GroupName,
		Notification:           m.Notification,
		Introduction:           m.Introduction,
		FaceURL:                m.FaceURL,
		Ex:                     m.Ex,
		Status:                 m.Status,
		CreatorUserID:          m.CreatorUserID,
		GroupType:              m.GroupType,
		NeedVerification:       m.NeedVerification,
		LookMemberInfo:         m.LookMemberInfo,
		ApplyMemberFriend:      m.ApplyMemberFriend,
		NotificationUserID:     m.NotificationUserID,
	}
}
