package logic

import (
	"context"
	"encoding/json"
	"time"

	"fdim/pkg/errs"
	"fdim/pkg/model"
	"fdim/pkg/util/idutil"
	"fdim/protocol/constant"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
)

type RevokeMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeMsgLogic {
	return &RevokeMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RevokeMsg 撤回消息：
// - 在 Mongo 中找到对应会话+seq 的消息，将其 contentType 标记为"撤回通知"，并清空内容；
// - 发送 MsgRevokeNotification (2101) 通知所有客户端。
func (l *RevokeMsgLogic) RevokeMsg(req *msg.RevokeMsgReq) (*msg.RevokeMsgResp, error) {
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID is required")
	}
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID is required")
	}
	if req.Seq < 0 {
		return nil, errs.ErrArgs.WrapMsg("seq is invalid")
	}

	coll := l.svcCtx.MongoDB.GetCollection("stream_msg")
	if coll == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message collection not initialized")
	}

	// 查找消息
	filter := bson.M{
		"conversation_id": req.ConversationID,
		"seq":             req.Seq,
	}
	var doc model.MsgDoc
	if err := coll.FindOne(l.ctx, filter).Decode(&doc); err != nil {
		l.Errorw("revoke FindOne failed", logx.Field("conversationID", req.ConversationID), logx.Field("seq", req.Seq), logx.Field("error", err))
		return nil, errs.ErrRecordNotFound.WrapMsg("msg not found")
	}

	// 将消息标记为"已撤回"：保留 Seq 等元数据，清空内容和附加字段。
	update := bson.M{
		"$set": bson.M{
			"content_type": 0,
			"content":      []byte{},
			"ex":           "",
		},
	}
	if _, err := coll.UpdateOne(l.ctx, filter, update); err != nil {
		l.Errorw("revoke UpdateOne failed", logx.Field("conversationID", req.ConversationID), logx.Field("seq", req.Seq), logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to revoke message")
	}

	// 发送 MsgRevokeNotification 通知（与官方一致）
	if l.svcCtx.SendMsgFunc != nil {
		go l.sendRevokeNotification(req.UserID, req.ConversationID, &doc)
	}

	return &msg.RevokeMsgResp{}, nil
}

func (l *RevokeMsgLogic) sendRevokeNotification(revokerID, conversationID string, doc *model.MsgDoc) {
	tips := &sdkws.MessageRevokedContent{
		RevokerID:                   revokerID,
		ClientMsgID:                 doc.ClientMsgID,
		RevokeTime:                  time.Now().UnixMilli(),
		Seq:                         doc.Seq,
		SessionType:                 doc.SessionType,
		SourceMessageSendID:         doc.SendID,
		SourceMessageSendTime:       doc.SendTime,
	}

	n := sdkws.NotificationElem{Detail: marshalRevokeToString(tips)}
	content, err := json.Marshal(&n)
	if err != nil {
		l.Errorw("marshal MsgRevokeNotification failed", logx.Field("error", err))
		return
	}

	// 撤回通知发给自己（多设备同步），sessionType=SingleChat
	msgData := &sdkws.MsgData{
		SendID:      revokerID,
		RecvID:      revokerID,
		Content:     content,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.MsgRevokeNotification,
		SessionType: constant.SingleChatType,
		CreateTime:  time.Now().UnixMilli(),
		ClientMsgID: idutil.GetMsgIDByMD5(revokerID),
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
		l.Errorw("send MsgRevokeNotification failed", logx.Field("error", err))
	}
}

func marshalRevokeToString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
