package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
)

type GetGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupsLogic {
	return &GetGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupsLogic) GetGroups(req *types.GetGroupsReq) (resp *types.GetGroupsResp, err error) {
	rpcReq := &user.GetGroupsReq{
		GroupName: req.GroupName,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	rpcResp, err := l.svcCtx.GroupClient.GetGroups(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var groups []interface{}
	for _, g := range rpcResp.Groups {
		groups = append(groups, g)
	}

	resp = &types.GetGroupsResp{
		Total:  rpcResp.Total,
		Groups: groups,
	}

	return resp, nil
}
