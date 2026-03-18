package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
)

type GetGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupMemberListLogic) GetGroupMemberList(req *types.GetGroupMemberListReq) (resp *types.GetGroupMemberListResp, err error) {
	rpcReq := &user.GetGroupMemberListReq{
		GroupID: req.GroupID,
		Filter:  req.Filter,
		Keyword: req.Keyword,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	rpcResp, err := l.svcCtx.GroupClient.GetGroupMemberList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var members []interface{}
	for _, m := range rpcResp.Members {
		members = append(members, m)
	}

	resp = &types.GetGroupMemberListResp{
		Total:   rpcResp.Total,
		Members: members,
	}

	return resp, nil
}
