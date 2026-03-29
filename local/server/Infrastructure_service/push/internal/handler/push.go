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

// PushHandler ʹ
type PushHandler struct {
	config          *config.Config
	svcCtx          *svc.ServiceContext
	pushConsumer    mq.Consumer
	offlineConsumer mq.Consumer
}

// NewPushHandler ʹ
func NewPushHandler(ctx context.Context, cfg *config.Config, svcCtx *svc.ServiceContext) (*PushHandler, error) {
	return &PushHandler{
		config:          cfg,
		svcCtx:          svcCtx,
		pushConsumer:    svcCtx.PushConsumer,
		offlineConsumer: svcCtx.OfflinePushConsumer,
	}, nil
}

// Start ͷ
func (h *PushHandler) Start(ctx context.Context) error {
	// 
	if h.pushConsumer != nil {
		go h.handlePush(ctx)
	}

	// 
	if h.offlineConsumer != nil {
		go h.handleOfflinePush(ctx)
	}

	logx.Info("Push handler started")
	return nil
}

// handlePush Ϣ
func (h *PushHandler) handlePush(ctx context.Context) {
	logx.Info("Starting to subscribe to toPush topic...")
	
	// ѭϢ
	for {
		err := h.pushConsumer.Subscribe(ctx, func(msg mq.Message) error {
			logx.Infof("Received message from toPush topic: key=%s", msg.Key())
			// Ϣ
			var pushMsg pbmsg.PushMsgDataToMQ
			if err := proto.Unmarshal(msg.Value(), &pushMsg); err != nil {
				logx.Errorf("Failed to unmarshal push message: %v", err)
				msg.Mark()
				msg.Commit()
				return err
			}
			
			logx.Infof("Processing push message: conversationID=%s, sendID=%s, recvID=%s, sessionType=%d",
				pushMsg.ConversationID, pushMsg.MsgData.SendID, pushMsg.MsgData.RecvID, pushMsg.MsgData.SessionType)

			//  MessageGateway 
			if h.svcCtx.MessageGatewayClient != nil && pushMsg.MsgData != nil {
				if pushMsg.MsgData.SessionType == 1 { // SingleChatType - 
					// ģ͸ + ߣ豸ͬٷһ£
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
							// ֻԽ߽߲ͣҪ
							offlineUserIDs = append(offlineUserIDs, uid)
						}
					}
					// ʧ  ͻˣٷһ£
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
					// Ⱥģ + ͻˣٷһ£
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

							// 
							batchResp, err := h.svcCtx.MessageGatewayClient.OnlineBatchPushOneMsg(ctx, &msggateway.OnlineBatchPushOneMsgReq{
								MsgData:       pushMsg.MsgData,
								PushToUserIDs: pushToUserIDs,
							})
							if err != nil {
								logx.Errorf("Failed to push group message online: %v", err)
							}

							// ͻˣҳδɹ͵ûųߣ
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
										continue // ߲Ҫ
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

// handleOfflinePush Ϣ
func (h *PushHandler) handleOfflinePush(ctx context.Context) {
	// ѭϢ
	for {
		err := h.offlineConsumer.Subscribe(ctx, func(msg mq.Message) error {
			// Ϣ
			var pushMsg pbmsg.PushMsgDataToMQ
			if err := proto.Unmarshal(msg.Value(), &pushMsg); err != nil {
				logx.Errorf("Failed to unmarshal offline push message: %v", err)
				msg.Mark()
				msg.Commit()
				return err
			}

			// ʹϢ
			if h.svcCtx.OfflinePusher != nil && pushMsg.MsgData != nil {
				// ϢȡûID
				userIDs := []string{}
				if pushMsg.MsgData.SessionType == 1 { // SingleChatType - 
					// ʱ RecvIDųߣ
					if pushMsg.MsgData.RecvID != "" && pushMsg.MsgData.RecvID != pushMsg.MsgData.SendID {
						userIDs = append(userIDs, pushMsg.MsgData.RecvID)
					}
				} else {
					// ȺҪ Group ȡԱб
					groupID := pushMsg.MsgData.GroupID
					if groupID == "" {
						groupID = getGroupIDFromConversationID(pushMsg.ConversationID)
					}

					if groupID != "" && h.svcCtx.GroupClient != nil {
						// ȡȺԱûIDб
						memberResp, err := h.svcCtx.GroupClient.GetGroupMemberUserIDs(ctx, &user.GetGroupMemberUserIDsReq{
							GroupID: groupID,
						})
						if err != nil {
							logx.Errorf("Failed to get group member user IDs for offline push: %v", err)
						} else if len(memberResp.UserIDs) > 0 {
							// ų
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

				// ȡ
				title := "New Message"
				content := "You have a new message"

				// Ϣʹ
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

// getGroupIDFromConversationID ӻỰIDȡȺID
func getGroupIDFromConversationID(conversationID string) string {
	if strings.HasPrefix(conversationID, "g_") {
		return strings.TrimPrefix(conversationID, "g_")
	} else if strings.HasPrefix(conversationID, "sg_") {
		return strings.TrimPrefix(conversationID, "sg_")
	}
	return ""
}

// Stop ֹͣͷ
func (h *PushHandler) Stop(ctx context.Context) error {
	logx.Info("Stopping push handler...")

	// ر NATS consumersֹͣϢ
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
