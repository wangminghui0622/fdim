package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type RespondFriendApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRespondFriendApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RespondFriendApplyLogic {
	return &RespondFriendApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RespondFriendApplyLogic) RespondFriendApply(req *types.RespondFriendApplyReq) (resp *types.RespondFriendApplyResp, err error) {
	handleMsg := req.HandleMsg
	if handleMsg == "" {
		handleMsg = " "
	}
	rpcReq := &user.RespondFriendApplyReq{
		FromUserID:   req.FromUserID,
		ToUserID:     req.ToUserID,
		HandleResult: req.HandleResult,
		HandleMsg:    handleMsg,
	}

	_, err = l.svcCtx.FriendClient.RespondFriendApply(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.RespondFriendApplyResp{
	}

	return resp, nil
}
