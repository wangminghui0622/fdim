package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
)

type DeleteFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFavoriteLogic {
	return &DeleteFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFavoriteLogic) DeleteFavorite(req *types.DeleteFavoriteReq, userID string) (resp *types.DeleteFavoriteResp, err error) {
	err = l.svcCtx.FavoriteDB.Delete(userID, req.FavoriteID)
	if err != nil {
		return nil, err
	}

	resp = &types.DeleteFavoriteResp{
	}

	return resp, nil
}
