package logic

import (
	"context"
	"time"

	"fdim/pkg/authverify"
	"fdim/pkg/errs"
	pkgmodel "fdim/pkg/model"
	"fdim/pkg/util/jsonutil"
	"fdim/protocol/constant"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
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

// RevokeMsg 处理消息撤回：
// 1. 校验权限
// 2. 从简化存储层（MsgCache/MsgDB）查找原消息
// 3. 更新 MongoDB 中的消息（标记为已撤回）并清除 Redis 缓存
// 4. 发送 MsgRevokeNotification 给会话另一端/群
func (l *RevokeMsgLogic) RevokeMsg(req *msg.RevokeMsgReq) (*msg.RevokeMsgResp, error) {
	l.Infof("[Revoke][Server] request: userID=%s conversationID=%s seq=%d", req.UserID, req.ConversationID, req.Seq)
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID is required")
	}
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID is required")
	}
	if req.Seq < 0 {
		return nil, errs.ErrArgs.WrapMsg("seq is invalid")
	}
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// 从简化存储层查找原消息（MsgCache → 简化 Redis 缓存 + MongoDB stream_msg）
	msgDocs, err := l.svcCtx.MsgCache.GetMessagesBySeq(l.ctx, req.ConversationID, []int64{req.Seq})
	if err != nil {
		l.Errorf("[Revoke][Server] MsgCache.GetMessagesBySeq failed: %v", err)
		return nil, err
	}
	if len(msgDocs) == 0 || msgDocs[0] == nil {
		// 缓存未命中，尝试从 MongoDB 查找
		if l.svcCtx.MsgDB != nil {
			msgDocs, err = l.svcCtx.MsgDB.GetMessagesBySeq(l.ctx, req.ConversationID, []int64{req.Seq})
			if err != nil {
				l.Errorf("[Revoke][Server] MsgDB.GetMessagesBySeq failed: %v", err)
				return nil, err
			}
		}
	}
	if len(msgDocs) == 0 || msgDocs[0] == nil {
		return nil, errs.ErrRecordNotFound.WrapMsg("msg not found")
	}

	origMsg := msgDocs[0]
	l.Infof("[Revoke][Server] target msg: clientMsgID=%s sendID=%s recvID=%s groupID=%s sessionType=%d contentType=%d seq=%d",
		origMsg.ClientMsgID, origMsg.SendID, origMsg.RecvID, origMsg.GroupID, origMsg.SessionType, origMsg.ContentType, origMsg.Seq)

	if origMsg.ContentType == constant.MsgRevokeNotification {
		return nil, errs.ErrMsgAlreadyRevoke.WrapMsg("msg already revoke")
	}

	var role int32
	if authverify.IsAdmin(l.ctx) {
		role = constant.AppAdmin
	}

	now := time.Now().UnixMilli()

	// 构建撤回后的消息内容（与官方格式一致）
	revokeContent := sdkws.MessageRevokedContent{
		RevokerID:                   req.UserID,
		RevokerRole:                 role,
		ClientMsgID:                 origMsg.ClientMsgID,
		RevokeTime:                  now,
		SourceMessageSendTime:       origMsg.SendTime,
		SourceMessageSendID:         origMsg.SendID,
		SourceMessageSenderNickname: origMsg.SenderNickname,
		SessionType:                 origMsg.SessionType,
		Seq:                         origMsg.Seq,
		Ex:                          origMsg.Ex,
	}
	revokeData, _ := jsonutil.JsonMarshal(&revokeContent)
	elem := sdkws.NotificationElem{Detail: string(revokeData)}
	contentBytes, _ := jsonutil.JsonMarshal(&elem)

	// 更新 MongoDB 中的消息（标记为已撤回）
	if l.svcCtx.MsgDB != nil {
		if err := l.svcCtx.MsgDB.RevokeMsg(l.ctx, req.ConversationID, req.Seq, constant.MsgRevokeNotification, contentBytes); err != nil {
			l.Errorf("[Revoke][Server] MsgDB.RevokeMsg failed: %v", err)
			return nil, err
		}
	}

	// 更新 Redis 缓存中的消息为已撤回状态（覆盖原缓存，保证 PullMsgBySeqs 返回正确状态）
	if l.svcCtx.MsgCache != nil {
		revokedMsg := *origMsg
		revokedMsg.ContentType = constant.MsgRevokeNotification
		revokedMsg.Content = contentBytes
		if err := l.svcCtx.MsgCache.SetMessagesToCache(l.ctx, req.ConversationID, []*pkgmodel.MsgDoc{&revokedMsg}); err != nil {
			l.Errorf("[Revoke][Server] MsgCache.SetMessagesToCache (revoked) failed: %v", err)
		}
	}
	l.Infof("[Revoke][Server] revoke stored: conversationID=%s seq=%d", req.ConversationID, req.Seq)

	// 发送撤回通知
	// 注意：官方架构中 Go SDK 会将 RevokeMsgTips 转换为 MessageRevokedContent 再回调给 Flutter 层
	// 本地架构无 Go SDK 中间层，所以直接发送 MessageRevokedContent，字段名与客户端 RevokedInfo 一致
	if l.svcCtx.NotificationSender != nil {
		var recvID string
		if origMsg.SessionType == constant.ReadGroupChatType {
			recvID = origMsg.GroupID
		} else {
			recvID = origMsg.RecvID
		}
		l.Infof("[Revoke][Server] notify: sendID=%s recvID=%s sessionType=%d clientMsgID=%s", req.UserID, recvID, origMsg.SessionType, origMsg.ClientMsgID)
		l.svcCtx.NotificationSender.NotificationWithSessionType(
			l.ctx,
			req.UserID,
			recvID,
			constant.MsgRevokeNotification,
			origMsg.SessionType,
			&revokeContent,
		)
	} else {
		l.Errorf("[Revoke][Server] NotificationSender is nil")
	}

	l.Infof("[Revoke][Server] request done: conversationID=%s seq=%d", req.ConversationID, req.Seq)
	return &msg.RevokeMsgResp{}, nil
}

