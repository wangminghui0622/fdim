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

type DeleteMsgsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMsgsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMsgsLogic {
	return &DeleteMsgsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteMsgs 选择性删除指定会话中的部分消息。
// 删除后发送 DeleteMsgsNotification (2102) 通知客户端同步。
func (l *DeleteMsgsLogic) DeleteMsgs(req *msg.DeleteMsgsReq) (*msg.DeleteMsgsResp, error) {
	if req.ConversationID == "" || len(req.Seqs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("conversationID and seqs are required")
	}
	if l.svcCtx.MsgDB == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message database not initialized")
	}

	// 从 Mongo 中删除
	if err := l.svcCtx.MsgDB.DeleteMessagesBySeq(l.ctx, req.ConversationID, req.Seqs); err != nil {
		l.Errorw("DeleteMessagesBySeq failed", logx.Field("conversationID", req.ConversationID), logx.Field("seqs", req.Seqs), logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to delete messages")
	}

	// 从 Redis 中清理这些消息缓存（忽略错误）
	if l.svcCtx.MsgCache != nil {
		for _, seq := range req.Seqs {
			key := l.svcCtx.MsgCacheKey(req.ConversationID, seq)
			if err := l.svcCtx.Redis.Del(l.ctx, key).Err(); err != nil {
				l.Errorw("failed to delete message cache", logx.Field("key", key), logx.Field("error", err))
			}
		}
	}

	// 发送 DeleteMsgsNotification 通知（与官方一致）
	if l.svcCtx.SendMsgFunc != nil && req.UserID != "" {
		go l.sendDeleteMsgsNotification(req.UserID, req.ConversationID, req.Seqs)
	}

	return &msg.DeleteMsgsResp{}, nil
}

func (l *DeleteMsgsLogic) sendDeleteMsgsNotification(userID, conversationID string, seqs []int64) {
	tips := &sdkws.DeleteMsgsTips{
		UserID:         userID,
		ConversationID: conversationID,
		Seqs:           seqs,
	}

	n := sdkws.NotificationElem{Detail: marshalDeleteToString(tips)}
	content, err := json.Marshal(&n)
	if err != nil {
		l.Errorw("marshal DeleteMsgsNotification failed", logx.Field("error", err))
		return
	}

	msgData := &sdkws.MsgData{
		SendID:      userID,
		RecvID:      userID,
		Content:     content,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.DeleteMsgsNotification,
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
		l.Errorw("send DeleteMsgsNotification failed", logx.Field("error", err))
	}
}

func marshalDeleteToString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
