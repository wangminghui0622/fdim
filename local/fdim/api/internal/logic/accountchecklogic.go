package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type AccountCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAccountCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountCheckLogic {
	return &AccountCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountCheckLogic) AccountCheck(req *types.AccountCheckReq) (resp *types.AccountCheckResp, err error) {
	// 转换请求参数
	rpcReq := &user.AccountCheckReq{
		CheckUserIDs: req.CheckUserIDs,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.AccountCheck(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var results []types.AccountCheckResult
	for _, r := range rpcResp.Results {
		results = append(results, types.AccountCheckResult{
			UserID: r.UserID,
			Result: r.AccountStatus,
		})
	}

	resp = &types.AccountCheckResp{
		Results: results,
	}

	return resp, nil
}
