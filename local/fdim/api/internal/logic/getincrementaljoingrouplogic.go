package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetIncrementalJoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIncrementalJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalJoinGroupLogic {
	return &GetIncrementalJoinGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIncrementalJoinGroupLogic) GetIncrementalJoinGroup(req *types.GetIncrementalJoinGroupReq) (resp *types.GetIncrementalJoinGroupResp, err error) {
	rpcReq := &user.GetIncrementalJoinGroupReq{
		UserID:    req.UserID,
		VersionID: req.VersionID,
		Version:   req.Version,
	}

	rpcResp, err := l.svcCtx.GroupClient.GetIncrementalJoinGroup(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var insert []interface{}
	for _, g := range rpcResp.Insert {
		insert = append(insert, g)
	}

	var update []interface{}
	for _, g := range rpcResp.Update {
		update = append(update, g)
	}

	groups := append(insert, update...)

	resp = &types.GetIncrementalJoinGroupResp{
		Version:   rpcResp.Version,
		VersionID: rpcResp.VersionID,
		Full:      rpcResp.Full,
		Delete:    rpcResp.Delete,
		Insert:    insert,
		Update:    update,
		Groups:    groups,
	}

	return resp, nil
}
