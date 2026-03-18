package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type ApplicationGroupResponseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplicationGroupResponseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplicationGroupResponseLogic {
	return &ApplicationGroupResponseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplicationGroupResponseLogic) ApplicationGroupResponse(req *types.ApplicationGroupResponseReq) (resp *types.ApplicationGroupResponseResp, err error) {
	rpcReq := &user.GroupApplicationResponseReq{
		GroupID:      req.GroupID,
		FromUserID:   req.FromUserID,
		HandledMsg:   req.HandledMsg,
		HandleResult: req.HandleResult,
	}

	_, err = l.svcCtx.GroupClient.GroupApplicationResponse(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.ApplicationGroupResponseResp{
	}

	return resp, nil
}
