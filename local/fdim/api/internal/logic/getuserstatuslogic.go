package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserStatusLogic {
	return &GetUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserStatusLogic) GetUserStatus(req *types.GetUserStatusReq) (resp *types.GetUserStatusResp, err error) {
	// 转换请求参数
	rpcReq := &user.GetUserStatusReq{
		UserID:  req.UserID,
		UserIDs: req.UserIDs,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.GetUserStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var statusList []types.OnlineStatus
	for _, s := range rpcResp.StatusList {
		statusList = append(statusList, types.OnlineStatus{
			UserID:      s.UserID,
			Status:      s.Status,
			PlatformIDs: s.PlatformIDs,
		})
	}

	resp = &types.GetUserStatusResp{
		StatusList: statusList,
	}

	return resp, nil
}
