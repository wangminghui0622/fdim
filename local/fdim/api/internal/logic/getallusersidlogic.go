package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/sdkws"
	"fdim/protocol/user"
)

type GetAllUsersIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllUsersIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllUsersIDLogic {
	return &GetAllUsersIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllUsersIDLogic) GetAllUsersID(req *types.GetAllUsersIDReq) (resp *types.GetAllUsersIDResp, err error) {
	// 转换请求参数
	rpcReq := &user.GetAllUserIDReq{
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.GetAllUserID(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.GetAllUsersIDResp{
		Total:   rpcResp.Total,
		UserIDs: rpcResp.UserIDs,
	}

	return resp, nil
}
