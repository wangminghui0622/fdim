package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/conversation"
	"fdim/protocol/sdkws"
)

type GetSortedConversationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSortedConversationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSortedConversationListLogic {
	return &GetSortedConversationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSortedConversationListLogic) GetSortedConversationList(req *types.GetSortedConversationListReq) (resp *types.GetSortedConversationListResp, err error) {
	// 转换请求参数
	rpcReq := &conversation.GetSortedConversationListReq{
		UserID:          req.UserID,
		ConversationIDs: req.ConversationIDs,
	}

	if req.Pagination.PageNumber > 0 || req.Pagination.ShowNumber > 0 {
		rpcReq.Pagination = &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		}
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.ConversationClient.GetSortedConversationList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应，保留 MsgInfo 以便前端获取 showName/faceURL
	var conversationElems []types.ConversationElem
	for _, elem := range rpcResp.ConversationElems {
		conversationElems = append(conversationElems, types.ConversationElem{
			ConversationID: elem.ConversationID,
			RecvMsgOpt:     elem.RecvMsgOpt,
			UnreadCount:    elem.UnreadCount,
			IsPinned:       elem.IsPinned,
			MsgInfo:        elem.MsgInfo,
		})
	}

	resp = &types.GetSortedConversationListResp{
		ConversationTotal: rpcResp.ConversationTotal,
		UnreadTotal:       rpcResp.UnreadTotal,
		ConversationElems: conversationElems,
	}

	return resp, nil
}
