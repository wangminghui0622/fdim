package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetGroupAbstractInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupAbstractInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupAbstractInfoLogic {
	return &GetGroupAbstractInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupAbstractInfoLogic) GetGroupAbstractInfo(req *types.GetGroupAbstractInfoReq) (resp *types.GetGroupAbstractInfoResp, err error) {
	rpcReq := &user.GetGroupAbstractInfoReq{
		GroupIDs: req.GroupIDs,
	}

	rpcResp, err := l.svcCtx.GroupClient.GetGroupAbstractInfo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var groupAbstractInfos []interface{}
	for _, gai := range rpcResp.GroupAbstractInfos {
		groupAbstractInfos = append(groupAbstractInfos, gai)
	}

	resp = &types.GetGroupAbstractInfoResp{
		GroupAbstractInfos: groupAbstractInfos,
	}

	return resp, nil
}
