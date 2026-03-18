package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
)

type DeleteConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConversationsLogic {
	return &DeleteConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteConversationsLogic) DeleteConversations(req *types.DeleteConversationsReq) (resp *types.DeleteConversationsResp, err error) {
	// 转换请求参数
	rpcReq := &conversation.DeleteConversationsReq{
		OwnerUserID:     req.OwnerUserID,
		NeedDeleteTime:  req.NeedDeleteTime,
		ConversationIDs: req.ConversationIDs,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.ConversationClient.DeleteConversations(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.DeleteConversationsResp{
	}

	return resp, nil
}
