package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
)

type GetPinnedConversationIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPinnedConversationIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPinnedConversationIDsLogic {
	return &GetPinnedConversationIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPinnedConversationIDsLogic) GetPinnedConversationIDs(req *types.GetPinnedConversationIDsReq) (resp *types.GetPinnedConversationIDsResp, err error) {
	// 转换请求参数
	rpcReq := &conversation.GetPinnedConversationIDsReq{
		UserID: req.UserID,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ConversationClient.GetPinnedConversationIDs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.GetPinnedConversationIDsResp{
		ConversationIDs: rpcResp.ConversationIDs,
	}

	return resp, nil
}
