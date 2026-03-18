package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetDesignatedFriendsApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDesignatedFriendsApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDesignatedFriendsApplyLogic {
	return &GetDesignatedFriendsApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDesignatedFriendsApplyLogic) GetDesignatedFriendsApply(req *types.GetDesignatedFriendsApplyReq) (resp *types.GetDesignatedFriendsApplyResp, err error) {
	rpcReq := &user.GetDesignatedFriendsApplyReq{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
	}

	rpcResp, err := l.svcCtx.FriendClient.GetDesignatedFriendsApply(l.ctx, rpcReq)
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

	resp = &types.GetDesignatedFriendsApplyResp{
		FriendRequests: friendRequests,
	}

	return resp, nil
}
