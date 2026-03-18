package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetSpecifiedUserGroupRequestInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSpecifiedUserGroupRequestInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSpecifiedUserGroupRequestInfoLogic {
	return &GetSpecifiedUserGroupRequestInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSpecifiedUserGroupRequestInfoLogic) GetSpecifiedUserGroupRequestInfo(req *types.GetSpecifiedUserGroupRequestInfoReq) (resp *types.GetSpecifiedUserGroupRequestInfoResp, err error) {
	rpcReq := &user.GetSpecifiedUserGroupRequestInfoReq{
		GroupID: req.GroupID,
		UserID:  req.UserID,
	}

	rpcResp, err := l.svcCtx.GroupClient.GetSpecifiedUserGroupRequestInfo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var groupRequests []interface{}
	for _, gr := range rpcResp.GroupRequests {
		groupRequests = append(groupRequests, gr)
	}

	resp = &types.GetSpecifiedUserGroupRequestInfoResp{
		Total:         rpcResp.Total,
		GroupRequests: groupRequests,
	}

	return resp, nil
}
