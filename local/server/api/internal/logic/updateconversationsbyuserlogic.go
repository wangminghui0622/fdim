package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
	"fdim/protocol/wrapperspb"
)

type UpdateConversationsByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateConversationsByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConversationsByUserLogic {
	return &UpdateConversationsByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateConversationsByUserLogic) UpdateConversationsByUser(req *types.UpdateConversationsByUserReq) (resp *types.UpdateConversationsByUserResp, err error) {
	// 转换请求参数
	rpcReq := &conversation.UpdateConversationsByUserReq{
		UserID: req.UserID,
	}

	if req.Ex != nil {
		rpcReq.Ex = wrapperspb.String(*req.Ex)
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.ConversationClient.UpdateConversationsByUser(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.UpdateConversationsByUserResp{
	}

	return resp, nil
}
