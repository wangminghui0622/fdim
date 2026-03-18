package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type DeleteLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogsLogic {
	return &DeleteLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogsLogic) DeleteLogs(req *types.DeleteLogsReq) (resp *types.DeleteLogsResp, err error) {
	// 转换请求参数
	rpcReq := &third.DeleteLogsReq{
		LogIDs: req.LogIDs,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.ThirdClient.DeleteLogs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.DeleteLogsResp{
	}

	return resp, nil
}
