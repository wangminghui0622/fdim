package rpcli

import (
	"context"
	"fdim/protocol/user"
	"google.golang.org/grpc"
)

func NewRelationClient(cc grpc.ClientConnInterface) *RelationClient {
	return &RelationClient{user.NewFriendClient(cc)}
}

type RelationClient struct {
	user.FriendClient
}

func (x *RelationClient) GetFriendsInfo(ctx context.Context, ownerUserID string, friendUserIDs []string) ([]*user.FriendInfoOnly, error) {
	if len(friendUserIDs) == 0 {
		return nil, nil
	}
	req := &user.GetFriendInfoReq{OwnerUserID: ownerUserID, FriendUserIDs: friendUserIDs}
	return extractField(ctx, x.FriendClient.GetFriendInfo, req, (*user.GetFriendInfoResp).GetFriendInfos)
}
