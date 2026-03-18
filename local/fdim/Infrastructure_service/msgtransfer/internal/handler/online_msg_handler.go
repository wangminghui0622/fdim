package handler

import (
	"context"

	"fdim/pkg/cache"
	"fdim/pkg/model"
	"fdim/pkg/mq"
	"fdim/protocol/constant"
	"fdim/protocol/conversation"
	pbmsg "fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

// OnlineMsgHandler 在线消息处理器：消费 [toRedis] topic，分发到 Redis 缓存、MongoDB 持久化、Push 推送
type OnlineMsgHandler struct {
	msgCache           cache.MsgCache
	conversationClient conversation.ConversationClient
	toMongoProducer    mq.Producer
	toPushProducer     mq.Producer
}

// NewOnlineMsgHandler 创建在线消息处理器
func NewOnlineMsgHandler(msgCache cache.MsgCache, conversationClient conversation.ConversationClient, toMongoProducer, toPushProducer mq.Producer) *OnlineMsgHandler {
	return &OnlineMsgHandler{
		msgCache:           msgCache,
		conversationClient: conversationClient,
		toMongoProducer:    toMongoProducer,
		toPushProducer:     toPushProducer,
	}
}

// HandleMessage 处理消息
func (h *OnlineMsgHandler) HandleMessage(msg mq.Message) error {
	ctx := msg.Context()
	key := msg.Key()
	value := msg.Value()

	// 解析消息（toRedis topic 发送的是 PushMsgDataToMQ）
	var pushMsg pbmsg.PushMsgDataToMQ
	if err := proto.Unmarshal(value, &pushMsg); err != nil {
		logx.Errorf("Failed to unmarshal message: %v, marking as processed to avoid redelivery", err)
		// 【修复】解析失败时也要标记消息已处理，避免重复消费
		msg.Mark()
		msg.Commit()
		return nil // 返回 nil 避免重试
	}

	// PushMsgDataToMQ 中的 msgData 是单个消息，不是数组
	if pushMsg.MsgData == nil {
		logx.Infof("No message to process")
		msg.Mark()
		msg.Commit()
		return nil
	}

	// 转换为 MsgDoc
	msgDoc := onlinePbToMsgDoc(pushMsg.ConversationID, pushMsg.MsgData)
	if msgDoc == nil {
		logx.Errorf("Failed to convert message to MsgDoc")
		msg.Mark()
		msg.Commit()
		return nil
	}

	// 创建会话（第一次发消息时，与官方一致）
	if h.conversationClient != nil && pushMsg.MsgData != nil {
		h.createConversationIfNeeded(ctx, pushMsg.ConversationID, pushMsg.MsgData)
	}

	// 存储到 Redis（使用 SetMessagesToCache，保留消息自带的 seq，不重新分配）
	if err := h.msgCache.SetMessagesToCache(ctx, pushMsg.ConversationID, []*model.MsgDoc{msgDoc}); err != nil {
		logx.Errorf("Failed to set message to cache: %v", err)
		// 继续处理，不返回错误
	}

	// 转发到 toMongo
	if h.toMongoProducer != nil {
		// 构造 MsgDataToMongoByMQ 消息
		mongoMsg := &pbmsg.MsgDataToMongoByMQ{
			ConversationID: pushMsg.ConversationID,
			MsgData:        []*sdkws.MsgData{pushMsg.MsgData},
		}

		// 获取当前最大序列号（从缓存中）
		maxSeq, err := h.msgCache.GetMaxSeq(ctx, pushMsg.ConversationID)
		if err == nil {
			mongoMsg.LastSeq = maxSeq
		}

		mongoMsgBytes, err := proto.Marshal(mongoMsg)
		if err != nil {
			logx.Errorf("Failed to marshal mongo message: %v", err)
		} else {
			if err := h.toMongoProducer.SendMessage(ctx, key, mongoMsgBytes); err != nil {
				logx.Errorf("Failed to forward message to toMongo: %v", err)
			} else {
				logx.Debugf("Forwarded message to toMongo: conversationID=%s", pushMsg.ConversationID)
			}
		}
	}

	// 转发到 toPush
	if h.toPushProducer != nil {
		logx.Infof("Forwarding message to toPush: conversationID=%s", pushMsg.ConversationID)
		// 重新序列化 PushMsgDataToMQ 并发送
		pushMsgBytes, err := proto.Marshal(&pushMsg)
		if err != nil {
			logx.Errorf("Failed to marshal push message: %v", err)
		} else {
			if err := h.toPushProducer.SendMessage(ctx, key, pushMsgBytes); err != nil {
				logx.Errorf("Failed to forward message to toPush: %v", err)
			} else {
				logx.Infof("Successfully forwarded message to toPush: conversationID=%s", pushMsg.ConversationID)
			}
		}
	} else {
		logx.Errorf("toPushProducer is nil, cannot forward message to push")
	}

	logx.Infof("Processed message: key=%s, conversationID=%s", key, pushMsg.ConversationID)

	// 标记消息已处理
	msg.Mark()
	msg.Commit()

	return nil
}

// onlinePbToMsgDoc converts protobuf MsgData to model.MsgDoc
func onlinePbToMsgDoc(conversationID string, pb *sdkws.MsgData) *model.MsgDoc {
	if pb == nil {
		return nil
	}
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
		CreateTime:       pb.CreateTime,  // 已经是 int64 时间戳
		SendTime:         pb.SendTime,
		Status:           pb.Status,
		Options:          pb.Options,
		AtUserIDs:        pb.AtUserIDList,
		AttachedInfo:     pb.AttachedInfo,
		Ex:               pb.Ex,
	}
}

// createConversationIfNeeded 创建会话（如果需要）
func (h *OnlineMsgHandler) createConversationIfNeeded(ctx context.Context, conversationID string, msgData *sdkws.MsgData) {
	// 根据 sessionType 判断是否需要创建会话
	switch msgData.SessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		// 单聊和通知类型：创建单聊会话
		req := &conversation.CreateSingleChatConversationsReq{
			RecvID:           msgData.RecvID,
			SendID:           msgData.SendID,
			ConversationID:   conversationID,
			ConversationType: msgData.SessionType,
		}
		if _, err := h.conversationClient.CreateSingleChatConversations(ctx, req); err != nil {
			// 会话可能已存在，只记录警告
			logx.Errorf("Create single chat conversation error: conversationID=%s, sessionType=%d, error=%v",
				conversationID, msgData.SessionType, err)
		} else {
			logx.Infof("Successfully created single chat conversation: conversationID=%s, sendID=%s, recvID=%s",
				conversationID, msgData.SendID, msgData.RecvID)
		}
	case constant.WriteGroupChatType, constant.ReadGroupChatType:
		// 群聊类型：创建群聊会话（需要调用 CreateGroupChatConversations）
		// 注意：群聊会话创建需要群成员列表，这里暂时跳过
		// 官方实现中会调用 GroupClient.GetGroupMemberUserIDs 获取成员列表
		logx.Infof("Group chat conversation creation skipped (not implemented yet): conversationID=%s, groupID=%s",
			conversationID, msgData.GroupID)
	default:
		logx.Errorf("Unknown session type, skip conversation creation: conversationID=%s, sessionType=%d",
			conversationID, msgData.SessionType)
	}
}
