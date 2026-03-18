package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetGroupsInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupsInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupsInfoLogic {
	return &GetGroupsInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupsInfoLogic) GetGroupsInfo(req *types.GetGroupsInfoReq) (resp *types.GetGroupsInfoResp, err error) {
	rpcReq := &user.GetGroupsInfoReq{
		GroupIDs: req.GroupIDs,
	}

	rpcResp, err := l.svcCtx.GroupClient.GetGroupsInfo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var groupInfos []interface{}
	for _, gi := range rpcResp.GroupInfos {
		groupInfos = append(groupInfos, gi)
	}

	resp = &types.GetGroupsInfoResp{
		GroupInfos: groupInfos,
	}

	return resp, nil
}
