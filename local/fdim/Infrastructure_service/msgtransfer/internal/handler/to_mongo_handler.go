package handler

import (
	"fdim/pkg/database"
	"fdim/pkg/model"
	"fdim/pkg/mq"
	pbmsg "fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

// ToMongoHandler 处理 toMongo topic 的消息
type ToMongoHandler struct {
	msgDB database.MsgDatabase
}

// NewToMongoHandler 创建 toMongo 处理器
func NewToMongoHandler(msgDB database.MsgDatabase) *ToMongoHandler {
	return &ToMongoHandler{
		msgDB: msgDB,
	}
}

// HandleMessage 处理消息
func (h *ToMongoHandler) HandleMessage(msg mq.Message) error {
	ctx := msg.Context()
	key := msg.Key()
	value := msg.Value()

	// 解析消息（toMongo topic 发送的是 MsgDataToMongoByMQ）
	var mongoMsg pbmsg.MsgDataToMongoByMQ
	if err := proto.Unmarshal(value, &mongoMsg); err != nil {
		logx.Errorf("Failed to unmarshal message: %v, marking as processed to avoid redelivery", err)
		// 【修复】解析失败时也要标记消息已处理，避免重复消费
		msg.Mark()
		msg.Commit()
		return nil // 返回 nil 避免重试
	}

	// 获取当前最大序列号
	currentMaxSeq, err := h.msgDB.GetMaxSeq(ctx, mongoMsg.ConversationID)
	if err != nil {
		logx.Errorf("Failed to get max seq: %v, marking as processed", err)
		// 【修复】获取失败时也标记消息已处理，避免重复消费
		msg.Mark()
		msg.Commit()
		return nil
	}

	// 转换消息为 MsgDoc
	msgDocs := make([]*model.MsgDoc, 0, len(mongoMsg.MsgData))
	for _, msgData := range mongoMsg.MsgData {
		msgDoc := pbToMsgDoc(mongoMsg.ConversationID, msgData)
		if msgDoc != nil {
			msgDocs = append(msgDocs, msgDoc)
		}
	}

	if len(msgDocs) == 0 {
		logx.Infof("No messages to persist")
		msg.Mark()
		msg.Commit()
		return nil
	}

	// 批量插入到 MongoDB
	if err := h.msgDB.BatchInsertChat2DB(ctx, mongoMsg.ConversationID, msgDocs, currentMaxSeq); err != nil {
		logx.Errorf("Failed to insert messages to MongoDB: %v, marking as processed", err)
		// 【修复】插入失败时也标记消息已处理，避免重复消费
		msg.Mark()
		msg.Commit()
		return nil
	}

	logx.Infof("Persisted %d messages to MongoDB: key=%s, conversationID=%s", len(msgDocs), key, mongoMsg.ConversationID)

	// 标记消息已处理
	msg.Mark()
	msg.Commit()

	return nil
}

// pbToMsgDoc converts protobuf MsgData to model.MsgDoc
func pbToMsgDoc(conversationID string, pb *sdkws.MsgData) *model.MsgDoc {
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
		CreateTime:       pb.CreateTime,  // 直接使用 int64 时间戳
		SendTime:         pb.SendTime,    // 直接使用 int64 时间戳
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
