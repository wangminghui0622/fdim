package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
	"fdim/protocol/sdkws"
)

type GetOwnerConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOwnerConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOwnerConversationLogic {
	return &GetOwnerConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOwnerConversationLogic) GetOwnerConversation(req *types.GetOwnerConversationReq) (resp *types.GetOwnerConversationResp, err error) {
	// 转换请求参数
	rpcReq := &conversation.GetOwnerConversationReq{
		UserID: req.UserID,
	}

	if req.Pagination.PageNumber > 0 || req.Pagination.ShowNumber > 0 {
		rpcReq.Pagination = &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		}
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ConversationClient.GetOwnerConversation(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var conversations []interface{}
	for _, c := range rpcResp.Conversations {
		conversations = append(conversations, c)
	}

	resp = &types.GetOwnerConversationResp{
		Total:         rpcResp.Total,
		Conversations: conversations,
	}

	return resp, nil
}
