package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
)

type ClearConversationsMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearConversationsMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearConversationsMsgLogic {
	return &ClearConversationsMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearConversationsMsgLogic) ClearConversationsMsg(req *types.ClearConversationsMsgReq) (resp *types.ClearConversationsMsgResp, err error) {
	// 转换请求参数
	rpcReq := &msg.ClearConversationsMsgReq{
		ConversationIDs: req.ConversationIDs,
		UserID:          req.UserID,
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.MsgClient.ClearConversationsMsg(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.ClearConversationsMsgResp{
	}

	return resp, nil
}
