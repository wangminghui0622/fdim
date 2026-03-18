package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/third"
)

type UploadLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogsLogic {
	return &UploadLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogsLogic) UploadLogs(req *types.UploadLogsReq) (resp *types.UploadLogsResp, err error) {
	// 转换请求参数
	var fileURLs []*third.FileURL
	for _, logItem := range req.Logs {
		if logMap, ok := logItem.(map[string]interface{}); ok {
			if urlStr, ok := logMap["url"].(string); ok {
				fileURL := &third.FileURL{
					Filename: urlStr,
					URL:      urlStr,
				}
				fileURLs = append(fileURLs, fileURL)
			}
		}
	}

	rpcReq := &third.UploadLogsReq{
		FileURLs: fileURLs,
		Platform:     0,
		AppFramework: "",
		Version:      "",
		Ex:           "",
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.ThirdClient.UploadLogs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.UploadLogsResp{
	}

	return resp, nil
}
