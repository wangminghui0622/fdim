package logic

import (
	"context"
	"strings"
	"time"

	"fdim/pkg/errs"
	"fdim/pkg/mcontext"
	"fdim/pkg/model"
	"fdim/pkg/msgprocessor"
	"fdim/pkg/util/idutil"
	"fdim/pkg/util/timeutil"
	"fdim/protocol/constant"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type SendMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMsgLogic) SendMsg(req *msg.SendMsgReq) (*msg.SendMsgResp, error) {
	// 1. 基本参数校验
	if req.GetMsgData() == nil {
		return nil, errs.ErrArgs.WrapMsg("msgData is required")
	}

	msgData := req.GetMsgData()
	if msgData.ClientMsgID == "" {
		return nil, errs.ErrArgs.WrapMsg("clientMsgID is required")
	}
	if msgData.SendID == "" {
		return nil, errs.ErrArgs.WrapMsg("sendID is required")
	}

	// 单聊必须有 recvID，群聊必须有 groupID
	if msgData.SessionType == constant.SingleChatType && msgData.RecvID == "" {
		return nil, errs.ErrArgs.WrapMsg("recvID is required for single chat")
	}
	if (msgData.SessionType == constant.WriteGroupChatType || msgData.SessionType == constant.ReadGroupChatType) && msgData.GroupID == "" {
		return nil, errs.ErrArgs.WrapMsg("groupID is required for group chat")
	}

	// 2. 计算会话 ID（与官方一致：通知消息使用 n_ 前缀，普通消息使用 si_/sg_ 前缀）
	conversationID := msgprocessor.GetConversationIDByMsg(msgData)
	if conversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("failed to generate conversationID")
	}
	
	// 验证 conversationID 格式
	if msgData.SessionType == constant.SingleChatType {
		// 单聊 conversationID 必须包含两个用户ID
		if !strings.Contains(conversationID, "_") || strings.HasSuffix(conversationID, "_") || strings.Contains(conversationID, "__") {
			l.Errorw("Invalid conversationID format",
				logx.Field("conversationID", conversationID),
				logx.Field("sendID", msgData.SendID),
				logx.Field("recvID", msgData.RecvID))
			return nil, errs.ErrArgs.WrapMsg("invalid conversationID format")
		}
	}

	// 3. 分配 seq（通过 MsgCache.IncrMaxSeq，内部委托给官方 SeqConversation.Malloc）
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	seq, err := l.svcCtx.MsgCache.IncrMaxSeq(l.ctx, conversationID)
	if err != nil {
		l.Errorw("IncrMaxSeq failed", logx.Field("error", err))
		return nil, errs.WrapMsg(err, "failed to allocate seq")
	}

	// 填充 protobuf 消息字段
	msgData.Seq = seq
	// 方案A：sendTime = 服务器 UTC 毫秒（不信任客户端时间，不加时区偏移）
	// 客户端显示时用 .toLocal() 转为手机本地时区
	// senderTimeZone 仅作为元数据保留在消息中
	serverNowMs := timeutil.GetCurrentTimestampByMill()
	msgData.CreateTime = serverNowMs
	msgData.SendTime = serverNowMs
	msgData.ServerMsgID = idutil.GetMsgIDByMD5(msgData.SendID)

	now := time.Now()

	// 4. 投递到 NATS [toRedis] topic → msgtransfer 负责：
	//    - 写 Redis 消息缓存
	//    - 转发到 [toMongo] → MongoDB 持久化
	//    - 转发到 [toPush]  → push 服务推送
	
	// 确保 context 包含必需的 operationID
	ctx := l.ctx
	if mcontext.GetOperationID(ctx) == "" {
		// 如果 context 中没有 operationID，使用 clientMsgID 作为 operationID
		ctx = mcontext.SetOperationID(ctx, msgData.ClientMsgID)
		l.Infow("Added operationID to context", logx.Field("operationID", msgData.ClientMsgID))
	}
	
	// 尝试发送到 NATS
	natsSuccess := false
	if l.svcCtx.ToRedisProducer != nil {
		pushMsgToMQ := &msg.PushMsgDataToMQ{
			MsgData:        msgData,
			ConversationID: conversationID,
		}
		data, err := proto.Marshal(pushMsgToMQ)
		if err != nil {
			l.Errorw("marshal PushMsgDataToMQ failed", logx.Field("error", err))
		} else {
			l.Infow("Sending message to NATS toRedis",
				logx.Field("conversationID", conversationID),
				logx.Field("seq", seq),
				logx.Field("clientMsgID", msgData.ClientMsgID))
			if err := l.svcCtx.ToRedisProducer.SendMessage(ctx, conversationID, data); err != nil {
				l.Errorw("produce to NATS toRedis failed", logx.Field("error", err))
			} else {
				natsSuccess = true
				l.Infow("Successfully sent message to NATS toRedis",
					logx.Field("conversationID", conversationID),
					logx.Field("seq", seq))
			}
		}
	} else {
		l.Errorw("ToRedisProducer is nil, cannot send to NATS", logx.Field("conversationID", conversationID))
	}
	
	// NATS 不可用或发送失败时，回退到直接写 MongoDB + Redis
	if !natsSuccess {
		l.Infow("Falling back to direct write", logx.Field("conversationID", conversationID))
		doc := msgDataToDoc(conversationID, msgData, now)
		
		// 1. 写入 MongoDB
		if l.svcCtx.MsgDB != nil {
			if err := l.svcCtx.MsgDB.BatchInsertChat2DB(l.ctx, conversationID, []*model.MsgDoc{doc}, seq); err != nil {
				l.Errorw("BatchInsertChat2DB fallback failed", logx.Field("error", err))
			} else {
				l.Infow("Message saved to MongoDB", logx.Field("conversationID", conversationID), logx.Field("seq", seq))
			}
		}
		
		// 2. 写入 Redis 缓存
		if l.svcCtx.MsgCache != nil {
			if err := l.svcCtx.MsgCache.SetMessagesToCache(l.ctx, conversationID, []*model.MsgDoc{doc}); err != nil {
				l.Errorw("SetMessagesToCache fallback failed", logx.Field("error", err))
			} else {
				l.Infow("Message cached to Redis", logx.Field("conversationID", conversationID), logx.Field("seq", seq))
			}
		}
		
		// 3. 推送给在线用户
		if l.svcCtx.ToPushProducer != nil {
			pushMsgToMQ := &msg.PushMsgDataToMQ{
				MsgData:        msgData,
				ConversationID: conversationID,
			}
			pushData, err := proto.Marshal(pushMsgToMQ)
			if err != nil {
				l.Errorw("marshal PushMsgDataToMQ for push failed", logx.Field("error", err))
			} else {
				if err := l.svcCtx.ToPushProducer.SendMessage(ctx, conversationID, pushData); err != nil {
					l.Errorw("send to toPush failed", logx.Field("error", err))
				} else {
					l.Infow("Message pushed to toPush queue", logx.Field("conversationID", conversationID))
				}
			}
		}
	}

	// 会话自动创建（与官方一致：非通知消息才创建会话）
	if l.svcCtx.ConversationClient != nil && !msgprocessor.IsNotificationByMsg(msgData) {
		go l.ensureConversation(msgData, conversationID)
	}

	// 详细日志：验证会话ID传递
	l.Infow("SendMsg success: conversationID generated and validated",
		logx.Field("conversationID", conversationID),
		logx.Field("seq", seq),
		logx.Field("clientMsgID", msgData.ClientMsgID),
		logx.Field("sendID", msgData.SendID),
		logx.Field("recvID", msgData.RecvID),
		logx.Field("groupID", msgData.GroupID),
		logx.Field("sessionType", msgData.SessionType))

	return &msg.SendMsgResp{
		ServerMsgID: msgData.ServerMsgID,
		ClientMsgID: msgData.ClientMsgID,
		SendTime:    msgData.SendTime,
		Modify:      msgData,
	}, nil
}

// ensureConversation 确保会话存在（与官方一致：仅单聊在 SendMsg 后自动创建，群聊由 group RPC 管理）
func (l *SendMsgLogic) ensureConversation(msgData *sdkws.MsgData, conversationID string) {
	if msgData.SessionType != constant.SingleChatType {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := l.svcCtx.ConversationClient.CreateSingleChatConversations(ctx,
		&conversation.CreateSingleChatConversationsReq{
			RecvID:           msgData.RecvID,
			SendID:           msgData.SendID,
			ConversationID:   conversationID,
			ConversationType: msgData.SessionType,
		})
	if err != nil {
		l.Infow("CreateSingleChatConversations may already exist",
			logx.Field("conversationID", conversationID), logx.Field("error", err))
	}
}

// msgDataToDoc 将 protobuf MsgData 转换为 model.MsgDoc（仅 NATS 不可用时的回退路径使用）
func msgDataToDoc(conversationID string, pb *sdkws.MsgData, now time.Time) *model.MsgDoc {
	return &model.MsgDoc{
		ConversationID:   conversationID,
		Seq:              pb.Seq,
		SendID:           pb.SendID,
		RecvID:           pb.RecvID,
		GroupID:          pb.GroupID,
		ClientMsgID:      pb.ClientMsgID,
		ServerMsgID:      pb.ServerMsgID,
		SenderPlatformID: pb.SenderPlatformID,
		SenderNickname:   pb.SenderNickname,
		SenderFaceURL:    pb.SenderFaceURL,
		SessionType:      pb.SessionType,
		MsgFrom:          pb.MsgFrom,
		ContentType:      pb.ContentType,
		Content:          pb.Content,
		CreateTime:       now.UnixMilli(),  // 转换为毫秒时间戳
		SendTime:         pb.SendTime,      // 直接使用 int64 时间戳
		Status:           pb.Status,
		Options:          pb.Options,
		AtUserIDs:        pb.AtUserIDList,
		AttachedInfo:     pb.AttachedInfo,
		Ex:               pb.Ex,
		IsRead:           false, // 默认未读
		ReadTime:         0,      // 0 表示未读
		BurnTime:         0,      // 0 表示未设置
	}
}
