package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetIncrementalBlacksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIncrementalBlacksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalBlacksLogic {
	return &GetIncrementalBlacksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIncrementalBlacksLogic) GetIncrementalBlacks(req *types.GetIncrementalBlacksReq) (resp *types.GetIncrementalBlacksResp, err error) {
	rpcReq := &user.GetIncrementalBlacksReq{
		UserID:    req.UserID,
		VersionID: req.VersionID,
		Version:   req.Version,
	}

	rpcResp, err := l.svcCtx.FriendClient.GetIncrementalBlacks(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var insert []types.BlackInfo
	for _, bi := range rpcResp.Insert {
		blackUserID := ""
		if bi.BlackUserInfo != nil {
			blackUserID = bi.BlackUserInfo.UserID
		}
		insert = append(insert, types.BlackInfo{
			OwnerUserID: bi.OwnerUserID,
			BlackUserID: blackUserID,
			CreateTime:  bi.CreateTime,
			AddSource:   bi.AddSource,
			Ex:          bi.Ex,
		})
	}

	var update []types.BlackInfo
	for _, bi := range rpcResp.Update {
		blackUserID := ""
		if bi.BlackUserInfo != nil {
			blackUserID = bi.BlackUserInfo.UserID
		}
		update = append(update, types.BlackInfo{
			OwnerUserID: bi.OwnerUserID,
			BlackUserID: blackUserID,
			CreateTime:  bi.CreateTime,
			AddSource:   bi.AddSource,
			Ex:          bi.Ex,
		})
	}

	resp = &types.GetIncrementalBlacksResp{
		Insert:    insert,
		Update:    update,
		Delete:    rpcResp.Delete,
		Version:   rpcResp.Version,
		VersionID: rpcResp.VersionID,
		Full:      rpcResp.Full,
	}

	return resp, nil
}
