package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type KickGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewKickGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickGroupMemberLogic {
	return &KickGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *KickGroupMemberLogic) KickGroupMember(req *types.KickGroupMemberReq) (resp *types.KickGroupMemberResp, err error) {
	rpcReq := &user.KickGroupMemberReq{
		GroupID:       req.GroupID,
		KickedUserIDs: req.KickedUserIDs,
		Reason:        req.Reason,
	}

	if req.SendMessage != nil {
		rpcReq.SendMessage = req.SendMessage
	}

	_, err = l.svcCtx.GroupClient.KickGroupMember(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.KickGroupMemberResp{
	}

	return resp, nil
}
