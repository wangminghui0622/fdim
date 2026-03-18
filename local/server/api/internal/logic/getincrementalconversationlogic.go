package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
)

type GetIncrementalConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIncrementalConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalConversationLogic {
	return &GetIncrementalConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIncrementalConversationLogic) GetIncrementalConversation(req *types.GetIncrementalConversationReq) (resp *types.GetIncrementalConversationResp, err error) {
	// 转换请求参数
	rpcReq := &conversation.GetIncrementalConversationReq{
		UserID:    req.UserID,
		VersionID: req.VersionID,
		Version:   req.Version,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ConversationClient.GetIncrementalConversation(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var insert []interface{}
	for _, c := range rpcResp.Insert {
		insert = append(insert, c)
	}

	var update []interface{}
	for _, c := range rpcResp.Update {
		update = append(update, c)
	}

	resp = &types.GetIncrementalConversationResp{
		Version:   rpcResp.Version,
		VersionID: rpcResp.VersionID,
		Full:      rpcResp.Full,
		Delete:    rpcResp.Delete,
		Insert:    insert,
		Update:    update,
	}

	return resp, nil
}
