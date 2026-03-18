package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
)

type GetUserReqGroupApplicationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserReqGroupApplicationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserReqGroupApplicationListLogic {
	return &GetUserReqGroupApplicationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserReqGroupApplicationListLogic) GetUserReqGroupApplicationList(req *types.GetUserReqGroupApplicationListReq) (resp *types.GetUserReqGroupApplicationListResp, err error) {
	rpcReq := &user.GetUserReqApplicationListReq{
		UserID:        req.UserID,
		GroupIDs:      req.GroupIDs,
		HandleResults: req.HandleResults,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	rpcResp, err := l.svcCtx.GroupClient.GetUserReqApplicationList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var groupRequests []interface{}
	for _, gr := range rpcResp.GroupRequests {
		groupRequests = append(groupRequests, gr)
	}

	resp = &types.GetUserReqGroupApplicationListResp{
		Total:         rpcResp.Total,
		GroupRequests: groupRequests,
	}

	return resp, nil
}
