package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type TransferGroupOwnerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransferGroupOwnerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferGroupOwnerLogic {
	return &TransferGroupOwnerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransferGroupOwnerLogic) TransferGroupOwner(req *types.TransferGroupOwnerReq) (resp *types.TransferGroupOwnerResp, err error) {
	rpcReq := &user.TransferGroupOwnerReq{
		GroupID:        req.GroupID,
		OldOwnerUserID: req.OldOwnerUserID,
		NewOwnerUserID: req.NewOwnerUserID,
	}

	_, err = l.svcCtx.GroupClient.TransferGroupOwner(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.TransferGroupOwnerResp{
	}

	return resp, nil
}
