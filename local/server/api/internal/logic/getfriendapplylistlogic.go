package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
)

type GetFriendApplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendApplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendApplyListLogic {
	return &GetFriendApplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendApplyListLogic) GetFriendApplyList(req *types.GetFriendApplyListReq) (resp *types.GetFriendApplyListResp, err error) {
	rpcReq := &user.GetPaginationFriendsApplyToReq{
		UserID: req.UserID,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	rpcResp, err := l.svcCtx.FriendClient.GetPaginationFriendsApplyTo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var friendRequests []types.FriendRequest
	for _, fr := range rpcResp.FriendRequests {
		friendRequests = append(friendRequests, types.FriendRequest{
			FromUserID:   fr.FromUserID,
			ToUserID:     fr.ToUserID,
			HandleResult: fr.HandleResult,
			ReqMsg:       fr.ReqMsg,
			HandleMsg:    fr.HandleMsg,
			CreateTime:   fr.CreateTime,
			HandleTime:   fr.HandleTime,
			Ex:           fr.Ex,
		})
	}

	resp = &types.GetFriendApplyListResp{
		FriendRequests: friendRequests,
		Total:          rpcResp.Total,
	}

	return resp, nil
}
