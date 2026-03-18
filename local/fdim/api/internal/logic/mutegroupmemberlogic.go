package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type MuteGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMuteGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteGroupMemberLogic {
	return &MuteGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MuteGroupMemberLogic) MuteGroupMember(req *types.MuteGroupMemberReq) (resp *types.MuteGroupMemberResp, err error) {
	rpcReq := &user.MuteGroupMemberReq{
		GroupID:      req.GroupID,
		UserID:       req.UserID,
		MutedSeconds: req.MutedSeconds,
	}

	_, err = l.svcCtx.GroupClient.MuteGroupMember(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.MuteGroupMemberResp{
	}

	return resp, nil
}
