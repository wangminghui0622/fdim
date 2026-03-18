package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
)

type GetSelfApplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSelfApplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelfApplyListLogic {
	return &GetSelfApplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSelfApplyListLogic) GetSelfApplyList(req *types.GetSelfApplyListReq) (resp *types.GetSelfApplyListResp, err error) {
	rpcReq := &user.GetPaginationFriendsApplyFromReq{
		UserID: req.UserID,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	rpcResp, err := l.svcCtx.FriendClient.GetPaginationFriendsApplyFrom(l.ctx, rpcReq)
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

	resp = &types.GetSelfApplyListResp{
		FriendRequests: friendRequests,
		Total:          rpcResp.Total,
	}

	return resp, nil
}
