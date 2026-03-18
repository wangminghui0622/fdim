package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
)

type GetRecvGroupApplicationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRecvGroupApplicationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecvGroupApplicationListLogic {
	return &GetRecvGroupApplicationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRecvGroupApplicationListLogic) GetRecvGroupApplicationList(req *types.GetRecvGroupApplicationListReq) (resp *types.GetRecvGroupApplicationListResp, err error) {
	rpcReq := &user.GetGroupApplicationListReq{
		FromUserID:    req.FromUserID,
		GroupIDs:      req.GroupIDs,
		HandleResults: req.HandleResults,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	rpcResp, err := l.svcCtx.GroupClient.GetGroupApplicationList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var groupRequests []interface{}
	for _, gr := range rpcResp.GroupRequests {
		groupRequests = append(groupRequests, gr)
	}

	resp = &types.GetRecvGroupApplicationListResp{
		Total:         rpcResp.Total,
		GroupRequests: groupRequests,
	}

	return resp, nil
}
