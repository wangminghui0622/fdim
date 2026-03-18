package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetIncrementalGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIncrementalGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalGroupMemberLogic {
	return &GetIncrementalGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIncrementalGroupMemberLogic) GetIncrementalGroupMember(req *types.GetIncrementalGroupMemberReq) (resp *types.GetIncrementalGroupMemberResp, err error) {
	rpcReq := &user.GetIncrementalGroupMemberReq{
		GroupID:   req.GroupID,
		VersionID: req.VersionID,
		Version:   req.Version,
	}

	rpcResp, err := l.svcCtx.GroupClient.GetIncrementalGroupMember(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var insert []interface{}
	for _, m := range rpcResp.Insert {
		insert = append(insert, m)
	}

	var update []interface{}
	for _, m := range rpcResp.Update {
		update = append(update, m)
	}

	resp = &types.GetIncrementalGroupMemberResp{
		Version:   rpcResp.Version,
		VersionID: rpcResp.VersionID,
		Full:      rpcResp.Full,
		Delete:    rpcResp.Delete,
		Insert:    insert,
		Update:    update,
	}

	return resp, nil
}
