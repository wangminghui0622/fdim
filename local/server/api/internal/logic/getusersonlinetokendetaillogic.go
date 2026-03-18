package logic

import (
	"context"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/constant"
	"fdim/protocol/msggateway"
)

type GetUsersOnlineTokenDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersOnlineTokenDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersOnlineTokenDetailLogic {
	return &GetUsersOnlineTokenDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersOnlineTokenDetailLogic) GetUsersOnlineTokenDetail(req *types.GetUsersOnlineTokenDetailReq) (resp *types.GetUsersOnlineTokenDetailResp, err error) {
	if l.svcCtx.MsgGatewayClient == nil {
		return nil, errs.ErrInternalServer.WrapMsg("msgGateway service is not configured")
	}

	// 转换请求参数
	rpcReq := &msggateway.GetUsersOnlineStatusReq{
		UserIDs: req.UserIDs,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.MsgGatewayClient.GetUsersOnlineStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var details []interface{}
	wsResult := rpcResp.SuccessResult

	// 遍历请求中的每个 userID
	for _, userID := range req.UserIDs {
		// 创建平台 token 映射
		platformTokenMap := make(map[int32][]string)

		for _, wsRes := range wsResult {
			if wsRes.UserID == userID {
				for _, status := range wsRes.DetailPlatformStatus {
					if tokens, ok := platformTokenMap[status.PlatformID]; ok {
						platformTokenMap[status.PlatformID] = append(tokens, status.Token)
					} else {
						platformTokenMap[status.PlatformID] = []string{status.Token}
					}
				}

				// 构建 SingleDetail
				singleDetail := &msggateway.SingleDetail{
					UserID: userID,
					Status: constant.Online,
				}

				for platformID, tokens := range platformTokenMap {
					singleDetail.SinglePlatformToken = append(singleDetail.SinglePlatformToken, &msggateway.SinglePlatformToken{
						PlatformID: platformID,
						Total:      int32(len(tokens)),
						Token:      tokens,
					})
				}

				details = append(details, singleDetail)
				break
			}
		}
	}

	resp = &types.GetUsersOnlineTokenDetailResp{
		Details: details,
	}

	return resp, nil
}
