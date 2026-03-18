package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
)

type GetPaginationBlacksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaginationBlacksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaginationBlacksLogic {
	return &GetPaginationBlacksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaginationBlacksLogic) GetPaginationBlacks(req *types.GetPaginationBlacksReq) (resp *types.GetPaginationBlacksResp, err error) {
	rpcReq := &user.GetPaginationBlacksReq{
		UserID: req.UserID,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	rpcResp, err := l.svcCtx.FriendClient.GetPaginationBlacks(l.ctx, rpcReq)
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

	resp = &types.GetPaginationBlacksResp{
		Blacks: blacks,
		Total:  rpcResp.Total,
	}

	return resp, nil
}
