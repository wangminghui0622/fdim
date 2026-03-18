package handler

import (
	"context"
	"strings"

	"fdim/Infrastructure_service/push/internal/config"
	"fdim/Infrastructure_service/push/internal/svc"

	"fdim/pkg/mq"
	pbmsg "fdim/protocol/msg"
	"fdim/protocol/msggateway"
	"fdim/protocol/user"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

// PushHandler 推送处理器
type PushHandler struct {
	config          *config.Config
	svcCtx          *svc.ServiceContext
	pushConsumer    mq.Consumer
	offlineConsumer mq.Consumer
}

// NewPushHandler 创建推送处理器
func NewPushHandler(ctx context.Context, cfg *config.Config, svcCtx *svc.ServiceContext) (*PushHandler, error) {
	return &PushHandler{
		config:          cfg,
		svcCtx:          svcCtx,
		pushConsumer:    svcCtx.PushConsumer,
		offlineConsumer: svcCtx.OfflinePushConsumer,
	}, nil
}

// Start 启动推送服务
func (h *PushHandler) Start(ctx context.Context) error {
	// 启动在线推送消费者
	if h.pushConsumer != nil {
		go h.handlePush(ctx)
	}

	// 启动离线推送消费者
	if h.offlineConsumer != nil {
		go h.handleOfflinePush(ctx)
	}

	logx.Info("Push handler started")
	return nil
}

// handlePush 处理在线推送消息
func (h *PushHandler) handlePush(ctx context.Context) {
	logx.Info("Starting to subscribe to toPush topic...")
	
	// 持续循环消费消息
	for {
		err := h.pushConsumer.Subscribe(ctx, func(msg mq.Message) error {
			logx.Infof("Received message from toPush topic: key=%s", msg.Key())
			// 解析消息
			var pushMsg pbmsg.PushMsgDataToMQ
			if err := proto.Unmarshal(msg.Value(), &pushMsg); err != nil {
				logx.Errorf("Failed to unmarshal push message: %v", err)
				msg.Mark()
				msg.Commit()
				return err
			}
			
			logx.Infof("Processing push message: conversationID=%s, sendID=%s, recvID=%s, sessionType=%d",
				pushMsg.ConversationID, pushMsg.MsgData.SendID, pushMsg.MsgData.RecvID, pushMsg.MsgData.SessionType)

			// 调用 MessageGateway 进行在线推送
			if h.svcCtx.MessageGatewayClient != nil && pushMsg.MsgData != nil {
				if pushMsg.MsgData.SessionType == 1 { // SingleChatType - 单聊
					// 单聊：推送给接收者 + 发送者（多设备同步，与官方一致）
					pushToUserIDs := []string{}
					if pushMsg.MsgData.RecvID != "" {
						pushToUserIDs = append(pushToUserIDs, pushMsg.MsgData.RecvID)
					}
					if pushMsg.MsgData.SendID != "" && pushMsg.MsgData.SendID != pushMsg.MsgData.RecvID {
						pushToUserIDs = append(pushToUserIDs, pushMsg.MsgData.SendID)
					}

					offlineUserIDs := []string{}
					for _, uid := range pushToUserIDs {
						onlinePush := false
						pushResp, err := h.svcCtx.MessageGatewayClient.OnlinePushMsg(ctx, &msggateway.OnlinePushMsgReq{
							MsgData:      pushMsg.MsgData,
							PushToUserID: uid,
						})
						if err != nil {
							logx.Errorf("Failed to push message online to %s: %v", uid, err)
						} else {
							for _, platform := range pushResp.Resp {
								if platform.ResultCode == 0 {
									onlinePush = true
									break
								}
							}
						}
						if !onlinePush && uid != pushMsg.MsgData.SendID {
							// 只对接收者进行离线推送，发送者不需要
							offlineUserIDs = append(offlineUserIDs, uid)
						}
					}
					// 在线推送失败 → 离线推送回退（与官方一致）
					if len(offlineUserIDs) > 0 && h.svcCtx.OfflinePusher != nil {
						title := "New Message"
						content := "You have a new message"
						if pushMsg.MsgData.OfflinePushInfo != nil {
							if pushMsg.MsgData.OfflinePushInfo.Title != "" {
								title = pushMsg.MsgData.OfflinePushInfo.Title
							}
							if pushMsg.MsgData.OfflinePushInfo.Desc != "" {
								content = pushMsg.MsgData.OfflinePushInfo.Desc
							}
						}
						if err := h.svcCtx.OfflinePusher.Push(ctx, offlineUserIDs, title, content, nil); err != nil {
							logx.Errorf("Offline push fallback failed: %v", err)
						} else {
							logx.Infof("Offline push fallback sent to users: %v", offlineUserIDs)
						}
					}
				} else {
					// 群聊：在线推送 + 离线推送回退（与官方一致）
					groupID := pushMsg.MsgData.GroupID
					if groupID == "" {
						groupID = getGroupIDFromConversationID(pushMsg.ConversationID)
					}

					if groupID != "" && h.svcCtx.GroupClient != nil {
						memberResp, err := h.svcCtx.GroupClient.GetGroupMemberUserIDs(ctx, &user.GetGroupMemberUserIDsReq{
							GroupID: groupID,
						})
						if err != nil {
							logx.Errorf("Failed to get group member user IDs: %v", err)
						} else if len(memberResp.UserIDs) > 0 {
							pushToUserIDs := memberResp.UserIDs

							// 在线推送
							batchResp, err := h.svcCtx.MessageGatewayClient.OnlineBatchPushOneMsg(ctx, &msggateway.OnlineBatchPushOneMsgReq{
								MsgData:       pushMsg.MsgData,
								PushToUserIDs: pushToUserIDs,
							})
							if err != nil {
								logx.Errorf("Failed to push group message online: %v", err)
							}

							// 离线推送回退：找出未成功在线推送的用户（排除发送者）
							if h.svcCtx.OfflinePusher != nil {
								onlineUserIDs := make(map[string]bool)
								if batchResp != nil {
									for _, result := range batchResp.SinglePushResult {
										for _, platform := range result.Resp {
											if platform.ResultCode == 0 {
												onlineUserIDs[result.UserID] = true
												break
											}
										}
									}
								}
								offlineUserIDs := make([]string, 0)
								for _, uid := range pushToUserIDs {
									if uid == pushMsg.MsgData.SendID {
										continue // 发送者不需要离线推送
									}
									if !onlineUserIDs[uid] {
										offlineUserIDs = append(offlineUserIDs, uid)
									}
								}
								if len(offlineUserIDs) > 0 {
									title := "New Message"
									content := "You have a new message"
									if pushMsg.MsgData.OfflinePushInfo != nil {
										if pushMsg.MsgData.OfflinePushInfo.Title != "" {
											title = pushMsg.MsgData.OfflinePushInfo.Title
										}
										if pushMsg.MsgData.OfflinePushInfo.Desc != "" {
											content = pushMsg.MsgData.OfflinePushInfo.Desc
										}
									}
									if err := h.svcCtx.OfflinePusher.Push(ctx, offlineUserIDs, title, content, nil); err != nil {
										logx.Errorf("Group offline push fallback failed: %v", err)
									} else {
										logx.Infof("Group offline push sent to %d users", len(offlineUserIDs))
									}
								}
							}
						}
					}
				}
			}

			msg.Mark()
			msg.Commit()
			return nil
		})
		if err != nil {
			logx.Errorf("Failed to subscribe to push topic: %v", err)
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
	}
}

// handleOfflinePush 处理离线推送消息
func (h *PushHandler) handleOfflinePush(ctx context.Context) {
	// 持续循环消费消息
	for {
		err := h.offlineConsumer.Subscribe(ctx, func(msg mq.Message) error {
			// 解析消息
			var pushMsg pbmsg.PushMsgDataToMQ
			if err := proto.Unmarshal(msg.Value(), &pushMsg); err != nil {
				logx.Errorf("Failed to unmarshal offline push message: %v", err)
				msg.Mark()
				msg.Commit()
				return err
			}

			// 使用离线推送器推送消息
			if h.svcCtx.OfflinePusher != nil && pushMsg.MsgData != nil {
				// 从消息中提取用户ID
				userIDs := []string{}
				if pushMsg.MsgData.SessionType == 1 { // SingleChatType - 单聊
					// 单聊时，接收者是 RecvID（排除发送者）
					if pushMsg.MsgData.RecvID != "" && pushMsg.MsgData.RecvID != pushMsg.MsgData.SendID {
						userIDs = append(userIDs, pushMsg.MsgData.RecvID)
					}
				} else {
					// 群聊需要从 Group 服务获取成员列表
					groupID := pushMsg.MsgData.GroupID
					if groupID == "" {
						groupID = getGroupIDFromConversationID(pushMsg.ConversationID)
					}

					if groupID != "" && h.svcCtx.GroupClient != nil {
						// 获取群成员用户ID列表
						memberResp, err := h.svcCtx.GroupClient.GetGroupMemberUserIDs(ctx, &user.GetGroupMemberUserIDsReq{
							GroupID: groupID,
						})
						if err != nil {
							logx.Errorf("Failed to get group member user IDs for offline push: %v", err)
						} else if len(memberResp.UserIDs) > 0 {
							// 排除发送者
							for _, userID := range memberResp.UserIDs {
								if userID != pushMsg.MsgData.SendID {
									userIDs = append(userIDs, userID)
								}
							}
						}
					}
				}

				if len(userIDs) == 0 {
					logx.Debugf("No user IDs to push offline")
					msg.Mark()
					msg.Commit()
					return nil
				}

				// 提取推送内容
				title := "New Message"
				content := "You have a new message"

				// 如果有离线推送信息，使用它
				if pushMsg.MsgData.OfflinePushInfo != nil {
					if pushMsg.MsgData.OfflinePushInfo.Title != "" {
						title = pushMsg.MsgData.OfflinePushInfo.Title
					}
					if pushMsg.MsgData.OfflinePushInfo.Desc != "" {
						content = pushMsg.MsgData.OfflinePushInfo.Desc
					}
				}

				if err := h.svcCtx.OfflinePusher.Push(ctx, userIDs, title, content, nil); err != nil {
					logx.Errorf("Failed to push message offline: %v", err)
				}
			}

			msg.Mark()
			msg.Commit()
			return nil
		})
		if err != nil {
			logx.Errorf("Failed to subscribe to offline push topic: %v", err)
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
	}
}

// getGroupIDFromConversationID 从会话ID中提取群组ID
func getGroupIDFromConversationID(conversationID string) string {
	if strings.HasPrefix(conversationID, "g_") {
		return strings.TrimPrefix(conversationID, "g_")
	} else if strings.HasPrefix(conversationID, "sg_") {
		return strings.TrimPrefix(conversationID, "sg_")
	}
	return ""
}

// Stop 停止推送服务
func (h *PushHandler) Stop(ctx context.Context) error {
	logx.Info("Stopping push handler...")

	// 关闭 NATS consumers（停止接收新消息）
	if h.pushConsumer != nil {
		if err := h.pushConsumer.Close(); err != nil {
			logx.Errorf("Failed to close push consumer: %v", err)
		} else {
			logx.Info("Push consumer closed")
		}
	}
	if h.offlineConsumer != nil {
		if err := h.offlineConsumer.Close(); err != nil {
			logx.Errorf("Failed to close offline push consumer: %v", err)
		} else {
			logx.Info("Offline push consumer closed")
		}
	}

	logx.Info("Push handler stopped")
	return nil
}
