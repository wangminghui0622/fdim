package logic

import (
	"context"
	"fmt"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecvMsgNotNotifyUserIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRecvMsgNotNotifyUserIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecvMsgNotNotifyUserIDsLogic {
	return &GetRecvMsgNotNotifyUserIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRecvMsgNotNotifyUserIDsLogic) GetRecvMsgNotNotifyUserIDs(req *conversation.GetRecvMsgNotNotifyUserIDsReq) (*conversation.GetRecvMsgNotNotifyUserIDsResp, error) {
	resp := &conversation.GetRecvMsgNotNotifyUserIDsResp{}

	// 根据 groupID 构建 conversationID
	// conversationID 格式: "group_" + groupID
	conversationID := "group_" + req.GroupID

	// 获取不接收消息通知的用户ID列表（recv_msg_opt == 1，ReceiveNotNotifyMessage?
	userIDs, err := l.svcCtx.ConversationDB.FindRecvMsgUserIDs(l.ctx, conversationID, []int32{1}) // 1 = ReceiveNotNotifyMessage
	if err != nil {
		return nil, fmt.Errorf("failed to get recv msg not notify user IDs: %w", err)
	}

	resp.UserIDs = userIDs
	return resp, nil
}
