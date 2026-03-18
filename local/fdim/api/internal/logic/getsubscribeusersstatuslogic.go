package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetSubscribeUsersStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSubscribeUsersStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscribeUsersStatusLogic {
	return &GetSubscribeUsersStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSubscribeUsersStatusLogic) GetSubscribeUsersStatus(req *types.GetSubscribeUsersStatusReq) (resp *types.GetSubscribeUsersStatusResp, err error) {
	// 转换请求参数
	rpcReq := &user.GetSubscribeUsersStatusReq{
		UserID: req.UserID,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.GetSubscribeUsersStatus(l.ctx, rpcReq)
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

	resp = &types.GetSubscribeUsersStatusResp{
		StatusList: statusList,
	}

	return resp, nil
}
