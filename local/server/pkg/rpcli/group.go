package rpcli

import (
	"context"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"google.golang.org/grpc"
)

func NewGroupClient(cc grpc.ClientConnInterface) *GroupClient {
	return &GroupClient{user.NewGroupClient(cc)}
}

type GroupClient struct {
	user.GroupClient
}

func (x *GroupClient) GetGroupsInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	req := &user.GetGroupsInfoReq{GroupIDs: groupIDs}
	return extractField(ctx, x.GroupClient.GetGroupsInfo, req, (*user.GetGroupsInfoResp).GetGroupInfos)
}

func (x *GroupClient) GetGroupInfo(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	return firstValue(x.GetGroupsInfo(ctx, []string{groupID}))
}

func (x *GroupClient) GetGroupInfoCache(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	req := &user.GetGroupInfoCacheReq{GroupID: groupID}
	return extractField(ctx, x.GroupClient.GetGroupInfoCache, req, (*user.GetGroupInfoCacheResp).GetGroupInfo)
}

func (x *GroupClient) GetGroupMemberCache(ctx context.Context, groupID string, userID string) (*sdkws.GroupMemberFullInfo, error) {
	req := &user.GetGroupMemberCacheReq{GroupID: groupID, GroupMemberID: userID}
	return extractField(ctx, x.GroupClient.GetGroupMemberCache, req, (*user.GetGroupMemberCacheResp).GetMember)
}

func (x *GroupClient) DismissGroup(ctx context.Context, groupID string, deleteMember bool) error {
	req := &user.DismissGroupReq{GroupID: groupID, DeleteMember: deleteMember}
	return ignoreResp(x.GroupClient.DismissGroup(ctx, req))
}

func (x *GroupClient) GetGroupMemberUserIDs(ctx context.Context, groupID string) ([]string, error) {
	req := &user.GetGroupMemberUserIDsReq{GroupID: groupID}
	return extractField(ctx, x.GroupClient.GetGroupMemberUserIDs, req, (*user.GetGroupMemberUserIDsResp).GetUserIDs)
}

func (x *GroupClient) GetGroupMembersInfo(ctx context.Context, groupID string, userIDs []string) ([]*sdkws.GroupMemberFullInfo, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	req := &user.GetGroupMembersInfoReq{GroupID: groupID, UserIDs: userIDs}
	return extractField(ctx, x.GroupClient.GetGroupMembersInfo, req, (*user.GetGroupMembersInfoResp).GetMembers)
}

func (x *GroupClient) GetGroupMemberInfo(ctx context.Context, groupID string, userID string) (*sdkws.GroupMemberFullInfo, error) {
	return firstValue(x.GetGroupMembersInfo(ctx, groupID, []string{userID}))
}

func (x *GroupClient) GetGroupMemberMapInfo(ctx context.Context, groupID string, userIDs []string) (map[string]*sdkws.GroupMemberFullInfo, error) {
	members, err := x.GetGroupMembersInfo(ctx, groupID, userIDs)
	if err != nil {
		return nil, err
	}
	memberMap := make(map[string]*sdkws.GroupMemberFullInfo)
	for _, member := range members {
		memberMap[member.UserID] = member
	}
	return memberMap, nil
}
