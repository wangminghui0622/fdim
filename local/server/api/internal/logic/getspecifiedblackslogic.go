package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetSpecifiedBlacksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSpecifiedBlacksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSpecifiedBlacksLogic {
	return &GetSpecifiedBlacksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSpecifiedBlacksLogic) GetSpecifiedBlacks(req *types.GetSpecifiedBlacksReq) (resp *types.GetSpecifiedBlacksResp, err error) {
	rpcReq := &user.GetSpecifiedBlacksReq{
		OwnerUserID: req.OwnerUserID,
		UserIDList:  req.UserIDList,
	}

	rpcResp, err := l.svcCtx.FriendClient.GetSpecifiedBlacks(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var blacks []types.BlackInfo
	for _, bi := range rpcResp.Blacks {
		blackUserID := ""
		if bi.BlackUserInfo != nil {
			blackUserID = bi.BlackUserInfo.UserID
		}
		blacks = append(blacks, types.BlackInfo{
			OwnerUserID: bi.OwnerUserID,
			BlackUserID: blackUserID,
			CreateTime:  bi.CreateTime,
			AddSource:   bi.AddSource,
			Ex:          bi.Ex,
		})
	}

	resp = &types.GetSpecifiedBlacksResp{
		Blacks: blacks,
	}

	return resp, nil
}
