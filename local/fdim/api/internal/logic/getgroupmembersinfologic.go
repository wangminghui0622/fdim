package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetGroupMembersInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupMembersInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMembersInfoLogic {
	return &GetGroupMembersInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupMembersInfoLogic) GetGroupMembersInfo(req *types.GetGroupMembersInfoReq) (resp *types.GetGroupMembersInfoResp, err error) {
	rpcReq := &user.GetGroupMembersInfoReq{
		GroupID: req.GroupID,
		UserIDs: req.UserIDs,
	}

	rpcResp, err := l.svcCtx.GroupClient.GetGroupMembersInfo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var members []interface{}
	for _, m := range rpcResp.Members {
		members = append(members, m)
	}

	resp = &types.GetGroupMembersInfoResp{
		Members: members,
	}

	return resp, nil
}
