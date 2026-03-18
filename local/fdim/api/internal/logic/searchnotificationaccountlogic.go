package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/sdkws"
	"fdim/protocol/user"
)

type SearchNotificationAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchNotificationAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchNotificationAccountLogic {
	return &SearchNotificationAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchNotificationAccountLogic) SearchNotificationAccount(req *types.SearchNotificationAccountReq) (resp *types.SearchNotificationAccountResp, err error) {
	// 转换请求参数
	rpcReq := &user.SearchNotificationAccountReq{
		Keyword: req.Keyword,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	}

	if req.AppManagerLevel != nil {
		rpcReq.AppManagerLevel = req.AppManagerLevel
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.SearchNotificationAccount(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var accounts []types.NotificationAccountInfo
	for _, acc := range rpcResp.NotificationAccounts {
		accounts = append(accounts, types.NotificationAccountInfo{
			UserID:         acc.UserID,
			FaceURL:        acc.FaceURL,
			NickName:       acc.NickName,
			AppMangerLevel: acc.AppMangerLevel,
		})
	}

	resp = &types.SearchNotificationAccountResp{
		Total:                rpcResp.Total,
		NotificationAccounts: accounts,
	}

	return resp, nil
}
