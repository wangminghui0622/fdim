package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
)

type GetFullOwnerConversationIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFullOwnerConversationIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFullOwnerConversationIDsLogic {
	return &GetFullOwnerConversationIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFullOwnerConversationIDsLogic) GetFullOwnerConversationIDs(req *types.GetFullOwnerConversationIDsReq) (resp *types.GetFullOwnerConversationIDsResp, err error) {
	// 转换请求参数
	rpcReq := &conversation.GetFullOwnerConversationIDsReq{
		IdHash: req.IdHash,
		UserID: req.UserID,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ConversationClient.GetFullOwnerConversationIDs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.GetFullOwnerConversationIDsResp{
		ConversationIDs: rpcResp.ConversationIDs,
	}

	return resp, nil
}
