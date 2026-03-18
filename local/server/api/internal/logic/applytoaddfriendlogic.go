package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type ApplyToAddFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyToAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyToAddFriendLogic {
	return &ApplyToAddFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyToAddFriendLogic) ApplyToAddFriend(req *types.ApplyToAddFriendReq) (resp *types.ApplyToAddFriendResp, err error) {
	reqMsg := req.ReqMsg
	if reqMsg == "" {
		reqMsg = " "
	}
	ex := req.Ex
	if ex == "" {
		ex = "{}"
	}
	rpcReq := &user.ApplyToAddFriendReq{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		ReqMsg:     reqMsg,
		Ex:         ex,
	}

	_, err = l.svcCtx.FriendClient.ApplyToAddFriend(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.ApplyToAddFriendResp{
	}

	return resp, nil
}
