package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
)

type GetJoinedGroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetJoinedGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJoinedGroupListLogic {
	return &GetJoinedGroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJoinedGroupListLogic) GetJoinedGroupList(req *types.GetJoinedGroupListReq) (resp *types.GetJoinedGroupListResp, err error) {
	rpcReq := &user.GetJoinedGroupListReq{
		FromUserID: req.FromUserID,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	rpcResp, err := l.svcCtx.GroupClient.GetJoinedGroupList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var groups []interface{}
	for _, g := range rpcResp.Groups {
		groups = append(groups, g)
	}

	resp = &types.GetJoinedGroupListResp{
		Total:  rpcResp.Total,
		Groups: groups,
	}

	return resp, nil
}
