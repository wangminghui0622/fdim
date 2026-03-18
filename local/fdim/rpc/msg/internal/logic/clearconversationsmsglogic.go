package logic

import (
	"context"
	"encoding/json"
	"time"

	"fdim/pkg/errs"
	"fdim/pkg/util/idutil"
	"fdim/protocol/constant"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ClearConversationsMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearConversationsMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearConversationsMsgLogic {
	return &ClearConversationsMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ClearConversationsMsg 清除指定会话消息：将 minSeq 设为 maxSeq+1，
// 然后发送 ClearConversationNotification (1703) 通知客户端。
func (l *ClearConversationsMsgLogic) ClearConversationsMsg(req *msg.ClearConversationsMsgReq) (*msg.ClearConversationsMsgResp, error) {
	if len(req.ConversationIDs) == 0 || req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationIDs and userID are required")
	}
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	// 对每个会话，将 minSeq 设为 maxSeq+1（逻辑清除）
	for _, convID := range req.ConversationIDs {
		maxSeq, err := l.svcCtx.MsgCache.GetMaxSeq(l.ctx, convID)
		if err != nil {
			l.Errorw("GetMaxSeq failed", logx.Field("conversationID", convID), logx.Field("error", err))
			continue
		}
		if err := l.svcCtx.MsgCache.SetMinSeq(l.ctx, convID, maxSeq+1); err != nil {
			l.Errorw("SetMinSeq failed", logx.Field("conversationID", convID), logx.Field("error", err))
		}
	}

	// 发送 ClearConversationNotification (1703) 通知（与官方一致）
	if l.svcCtx.SendMsgFunc != nil {
		go l.sendClearConversationNotification(req.UserID, req.ConversationIDs)
	}

	return &msg.ClearConversationsMsgResp{}, nil
}

func (l *ClearConversationsMsgLogic) sendClearConversationNotification(userID string, conversationIDs []string) {
	tips := &sdkws.ClearConversationTips{
		UserID:          userID,
		ConversationIDs: conversationIDs,
	}

	n := sdkws.NotificationElem{Detail: marshalClearToString(tips)}
	content, err := json.Marshal(&n)
	if err != nil {
		l.Errorw("marshal ClearConversationNotification failed", logx.Field("error", err))
		return
	}

	msgData := &sdkws.MsgData{
		SendID:      userID,
		RecvID:      userID,
		Content:     content,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.ClearConversationNotification,
		SessionType: constant.SingleChatType,
		CreateTime:  time.Now().UnixMilli(),
		ClientMsgID: idutil.GetMsgIDByMD5(userID),
		Options: map[string]bool{
			constant.IsNotNotification:          false,
			constant.IsConversationUpdate:       false,
			constant.IsSenderConversationUpdate: false,
			constant.IsUnreadCount:              false,
			constant.IsOfflinePush:              false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := l.svcCtx.SendMsgFunc(ctx, &msg.SendMsgReq{MsgData: msgData}); err != nil {
		l.Errorw("send ClearConversationNotification failed", logx.Field("error", err))
	}
}

func marshalClearToString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
