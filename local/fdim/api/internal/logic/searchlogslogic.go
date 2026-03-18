package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/sdkws"
	"fdim/protocol/third"
)

type SearchLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogsLogic {
	return &SearchLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogsLogic) SearchLogs(req *types.SearchLogsReq) (resp *types.SearchLogsResp, err error) {
	// 转换请求参数
	rpcReq := &third.SearchLogsReq{
		Keyword:   req.Keyword,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if req.Pagination.PageNumber > 0 || req.Pagination.ShowNumber > 0 {
		rpcReq.Pagination = &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		}
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ThirdClient.SearchLogs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var logs []interface{}
	for _, log := range rpcResp.LogsInfos {
		logs = append(logs, log)
	}

	resp = &types.SearchLogsResp{
		Logs:  logs,
		Total: int64(rpcResp.Total),
	}

	return resp, nil
}
