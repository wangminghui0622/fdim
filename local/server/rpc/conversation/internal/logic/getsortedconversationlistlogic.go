package logic

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"fdim/pkg/util/datautil"
	"fdim/protocol/constant"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/protocol/user"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSortedConversationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSortedConversationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSortedConversationListLogic {
	return &GetSortedConversationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSortedConversationListLogic) GetSortedConversationList(req *conversation.GetSortedConversationListReq) (*conversation.GetSortedConversationListResp, error) {
	resp := &conversation.GetSortedConversationListResp{}

	// 1. 获取会话ID列表
	var conversationIDs []string
	var err error
	if len(req.ConversationIDs) == 0 {
		allConvIDs, err := l.svcCtx.ConversationDB.FindUserIDAllConversationID(l.ctx, req.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to find conversation IDs: %w", err)
		}
		// 过滤掉通知会话（n_ 前缀），这些不应该显示在会话列表?
		conversationIDs = make([]string, 0, len(allConvIDs))
		for _, convID := range allConvIDs {
			if !strings.HasPrefix(convID, "n_") {
				conversationIDs = append(conversationIDs, convID)
			}
		}
	} else {
		conversationIDs = req.ConversationIDs
	}

	// 2. 获取会话详情
	conversations, err := l.svcCtx.ConversationDB.Find(l.ctx, req.UserID, conversationIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find conversations: %w", err)
	}

	if len(conversations) == 0 {
		return &conversation.GetSortedConversationListResp{
			ConversationTotal:  0,
			ConversationElems:  []*conversation.ConversationElem{},
			UnreadTotal:        0,
		}, nil
	}

	// 3. 获取每个会话的最大序列号
	convIDs := make([]string, 0, len(conversations))
	for _, conv := range conversations {
		if conv.ConversationID != "" {
			convIDs = append(convIDs, conv.ConversationID)
		}
	}

	l.Infof("Found %d conversations, convIDs: %v", len(convIDs), convIDs)

	maxSeqs := make(map[string]int64)
	if l.svcCtx.MsgClient != nil && len(convIDs) > 0 {
		l.Infof("Calling GetMaxSeqs with %d conversationIDs", len(convIDs))
		maxSeqResp, err := l.svcCtx.MsgClient.GetMaxSeqs(l.ctx, &msg.GetMaxSeqsReq{
			ConversationIDs: convIDs,
		})
		if err != nil {
			l.Errorf("failed to get max seqs: %v", err)
		} else if maxSeqResp != nil {
			l.Infof("GetMaxSeqs returned %d seqs: %v", len(maxSeqResp.MaxSeqs), maxSeqResp.MaxSeqs)
			maxSeqs = maxSeqResp.MaxSeqs
		}
	} else {
		if l.svcCtx.MsgClient == nil {
			l.Errorf("MsgClient is nil")
		}
		if len(convIDs) == 0 {
			l.Infof("No conversations found")
		}
	}

	// 4. 获取每个会话的最新消?
	chatLogs := make(map[string]*sdkws.MsgData)
	if l.svcCtx.MsgClient != nil && len(maxSeqs) > 0 {
		l.Infof("Calling GetMsgByConversationIDs with %d conversations, maxSeqs: %v", len(convIDs), maxSeqs)
		chatLogResp, err := l.svcCtx.MsgClient.GetMsgByConversationIDs(l.ctx, &msg.GetMsgByConversationIDsReq{
			ConversationIDs: convIDs,
			MaxSeqs:         maxSeqs,
		})
		if err != nil {
			l.Errorf("failed to get messages by conversation IDs: %v", err)
		} else if chatLogResp != nil {
			l.Infof("GetMsgByConversationIDs returned %d messages", len(chatLogResp.MsgDatas))
			chatLogs = chatLogResp.MsgDatas
		} else {
			l.Errorf("GetMsgByConversationIDs returned nil response")
		}
	} else {
		if l.svcCtx.MsgClient == nil {
			l.Errorf("MsgClient is nil, cannot get messages")
		}
		if len(maxSeqs) == 0 {
			l.Infof("maxSeqs is empty, no messages to fetch")
		}
	}

	// 5. 获取会话消息信息（用户头像、昵称等?
	conversationMsg, err := l.getConversationInfo(l.ctx, chatLogs, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation info: %w", err)
	}

	// 6. 获取已读序列?
	hasReadSeqs := make(map[string]int64)
	if l.svcCtx.MsgClient != nil && len(convIDs) > 0 {
		hasReadResp, err := l.svcCtx.MsgClient.GetHasReadSeqs(l.ctx, &msg.GetHasReadSeqsReq{
			UserID:          req.UserID,
			ConversationIDs: convIDs,
		})
		if err != nil {
			l.Errorf("failed to get has read seqs: %v", err)
		} else if hasReadResp != nil {
			hasReadSeqs = hasReadResp.MaxSeqs // 注意：返回的字段名是 MaxSeqs
		}
	}

	// 7. 计算未读?
	var unreadTotal int64
	conversationUnreadCount := make(map[string]int64)
	for conversationID, maxSeq := range maxSeqs {
		unreadCount := maxSeq - hasReadSeqs[conversationID]
		if unreadCount < 0 {
			unreadCount = 0
		}
		conversationUnreadCount[conversationID] = unreadCount
		unreadTotal += unreadCount
	}

	// 8. 按置顶和时间排序
	conversationIsPinTime := make(map[int64]string)
	conversationNotPinTime := make(map[int64]string)

	// 收集没有消息的会话的对方 userID，用于查询用户信?
	noMsgConversations := make(map[string]string) // conversationID -> otherUserID
	for _, v := range conversations {
		if _, ok := conversationMsg[v.ConversationID]; !ok {
			// 单聊会话：v.UserID 是对方的 userID
			if v.ConversationType == constant.SingleChatType && v.UserID != "" {
				noMsgConversations[v.ConversationID] = v.UserID
			}
		}
	}

	// 查询没有消息的会话的用户信息
	if len(noMsgConversations) > 0 && l.svcCtx.UserClient != nil {
		userIDs := make([]string, 0, len(noMsgConversations))
		for _, userID := range noMsgConversations {
			userIDs = append(userIDs, userID)
		}
		
		l.Infof("Fetching user info for conversations without messages: %d users", len(userIDs))
		userInfoResp, err := l.svcCtx.UserClient.GetDesignateUsers(l.ctx, &user.GetDesignateUsersReq{
			UserIDs: userIDs,
		})
		if err != nil {
			l.Errorf("failed to get user info for no-msg conversations: %v", err)
		} else if userInfoResp != nil {
			userMap := make(map[string]*sdkws.UserInfo)
			for _, userInfo := range userInfoResp.UsersInfo {
				userMap[userInfo.UserID] = userInfo
			}
			
			// 为没有消息的会话创建 ConversationElem，并填充用户信息
			for convID, otherUserID := range noMsgConversations {
				if userInfo, ok := userMap[otherUserID]; ok {
					conversationMsg[convID] = &conversation.ConversationElem{
						ConversationID: convID,
						MsgInfo: &conversation.MsgInfo{
							FaceURL:    userInfo.FaceURL,
							SenderName: userInfo.Nickname,
						},
					}
					l.Infof("Set user info for no-msg conversation: convID=%s, userID=%s, nickname=%s, faceURL=%s",
						convID, otherUserID, userInfo.Nickname, userInfo.FaceURL)
				}
			}
		}
	}

	for _, v := range conversations {
		conversationID := v.ConversationID
		var time int64
		if msgInfo, ok := conversationMsg[conversationID]; ok {
			if msgInfo.MsgInfo != nil {
				time = msgInfo.MsgInfo.LatestMsgRecvTime
			} else {
				// 没有消息，使用会话创建时?
				time = v.CreateTime.UnixMilli()
			}
		} else {
			// 如果还是没有找到（比如群聊或其他类型），创建空的
			conversationMsg[conversationID] = &conversation.ConversationElem{
				ConversationID: conversationID,
				IsPinned:       v.IsPinned,
				MsgInfo:        nil,
			}
			time = v.CreateTime.UnixMilli()
		}

		conversationMsg[conversationID].RecvMsgOpt = v.RecvMsgOpt
		if v.IsPinned {
			conversationMsg[conversationID].IsPinned = v.IsPinned
			conversationIsPinTime[time] = conversationID
			continue
		}
		conversationNotPinTime[time] = conversationID
	}

	resp = &conversation.GetSortedConversationListResp{
		ConversationTotal:  int64(len(chatLogs)),
		ConversationElems:  []*conversation.ConversationElem{},
		UnreadTotal:        unreadTotal,
	}

	// 9. 排序并添加到响应
	l.conversationSort(conversationIsPinTime, resp, conversationUnreadCount, conversationMsg)
	l.conversationSort(conversationNotPinTime, resp, conversationUnreadCount, conversationMsg)

	// 10. 分页
	if req.Pagination != nil {
		pageNumber := int(req.Pagination.PageNumber)
		showNumber := int(req.Pagination.ShowNumber)
		resp.ConversationElems = datautil.Paginate(resp.ConversationElems, pageNumber, showNumber)
	}

	return resp, nil
}

// conversationSort 对会话进行排序（按时间倒序?
func (l *GetSortedConversationListLogic) conversationSort(
	conversations map[int64]string,
	resp *conversation.GetSortedConversationListResp,
	conversationUnreadCount map[string]int64,
	conversationMsg map[string]*conversation.ConversationElem,
) {
	keys := make([]int64, 0, len(conversations))
	for key := range conversations {
		keys = append(keys, key)
	}

	// 按时间倒序排序
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	cons := make([]*conversation.ConversationElem, 0, len(conversations))
	for _, v := range keys {
		conversationID := conversations[v]
		conversationElem := conversationMsg[conversationID]
		conversationElem.UnreadCount = conversationUnreadCount[conversationID]
		cons = append(cons, conversationElem)
	}
	resp.ConversationElems = append(resp.ConversationElems, cons...)
}

// getConversationInfo 获取会话信息（用户头像、昵称、群组信息等?
func (l *GetSortedConversationListLogic) getConversationInfo(
	ctx context.Context,
	chatLogs map[string]*sdkws.MsgData,
	userID string,
) (map[string]*conversation.ConversationElem, error) {
	var (
		sendIDs         []string
		groupIDs        []string
		sendMap         = make(map[string]*sdkws.UserInfo)
		groupMap        = make(map[string]*sdkws.GroupInfo)
		conversationMsg = make(map[string]*conversation.ConversationElem)
	)

	// 收集需要查询的用户ID和群组ID
	for _, chatLog := range chatLogs {
		switch chatLog.SessionType {
		case constant.SingleChatType:
			if chatLog.SendID == userID {
				sendIDs = append(sendIDs, chatLog.RecvID)
			}
			sendIDs = append(sendIDs, chatLog.SendID)
		case constant.WriteGroupChatType, constant.ReadGroupChatType:
			groupIDs = append(groupIDs, chatLog.GroupID)
			sendIDs = append(sendIDs, chatLog.SendID)
		}
	}

	// 查询用户信息
	if len(sendIDs) != 0 && l.svcCtx.UserClient != nil {
		l.Infof("Calling GetDesignateUsers with %d userIDs: %v", len(sendIDs), sendIDs)
		sendInfoResp, err := l.svcCtx.UserClient.GetDesignateUsers(ctx, &user.GetDesignateUsersReq{
			UserIDs: sendIDs,
		})
		if err != nil {
			l.Errorf("failed to get users info: %v", err)
		} else if sendInfoResp != nil {
			l.Infof("GetDesignateUsers returned %d users", len(sendInfoResp.UsersInfo))
			for _, sendInfo := range sendInfoResp.UsersInfo {
				sendMap[sendInfo.UserID] = sendInfo
				l.Infof("User info: userID=%s, nickname=%s, faceURL=%s", sendInfo.UserID, sendInfo.Nickname, sendInfo.FaceURL)
			}
		}
	} else {
		if len(sendIDs) == 0 {
			l.Infof("No sendIDs to query")
		}
		if l.svcCtx.UserClient == nil {
			l.Errorf("UserClient is nil, cannot get user info")
		}
	}

	// 查询群组信息
	if len(groupIDs) != 0 && l.svcCtx.GroupClient != nil {
		groupInfoResp, err := l.svcCtx.GroupClient.GetGroupsInfo(ctx, &user.GetGroupsInfoReq{
			GroupIDs: groupIDs,
		})
		if err != nil {
			l.Errorf("failed to get groups info: %v", err)
		} else if groupInfoResp != nil {
			for _, groupInfo := range groupInfoResp.GroupInfos {
				groupMap[groupInfo.GroupID] = groupInfo
			}
		}
	}

	// 构建会话消息信息
	for conversationID, chatLog := range chatLogs {
		pbchatLog := &conversation.ConversationElem{}
		msgInfo := &conversation.MsgInfo{}

		// 复制基本字段
		msgInfo.ServerMsgID = chatLog.ServerMsgID
		msgInfo.ClientMsgID = chatLog.ClientMsgID
		msgInfo.SessionType = chatLog.SessionType
		msgInfo.SendID = chatLog.SendID
		msgInfo.RecvID = chatLog.RecvID
		msgInfo.GroupID = chatLog.GroupID
		msgInfo.ContentType = chatLog.ContentType
		msgInfo.Content = string(chatLog.Content) // Content ?[]byte，需要转换为 string
		msgInfo.Ex = chatLog.Ex
		msgInfo.LatestMsgRecvTime = chatLog.SendTime
		msgInfo.MsgFrom = chatLog.MsgFrom

		switch chatLog.SessionType {
		case constant.SingleChatType:
			// 单聊：显示对方的信息
			if chatLog.SendID == userID {
				// 我发送的消息，显示接收方信息
				l.Infof("Single chat: I am sender, showing receiver info. conversationID=%s, recvID=%s", conversationID, chatLog.RecvID)
				if recv, ok := sendMap[chatLog.RecvID]; ok {
					msgInfo.FaceURL = recv.FaceURL
					msgInfo.SenderName = recv.Nickname
					l.Infof("Set receiver info: faceURL=%s, senderName=%s", msgInfo.FaceURL, msgInfo.SenderName)
				} else {
					l.Errorf("Receiver info not found in sendMap: recvID=%s", chatLog.RecvID)
				}
				break
			}
			// 对方发送的消息，显示发送方信息
			l.Infof("Single chat: Other is sender, showing sender info. conversationID=%s, sendID=%s", conversationID, chatLog.SendID)
			if send, ok := sendMap[chatLog.SendID]; ok {
				msgInfo.FaceURL = send.FaceURL
				msgInfo.SenderName = send.Nickname
				l.Infof("Set sender info: faceURL=%s, senderName=%s", msgInfo.FaceURL, msgInfo.SenderName)
			} else {
				l.Errorf("Sender info not found in sendMap: sendID=%s", chatLog.SendID)
			}
		case constant.WriteGroupChatType, constant.ReadGroupChatType:
			// 群聊：显示群组信息和发送者昵?
			msgInfo.GroupID = chatLog.GroupID
			if group, ok := groupMap[chatLog.GroupID]; ok {
				msgInfo.GroupName = group.GroupName
				msgInfo.GroupFaceURL = group.FaceURL
				msgInfo.GroupMemberCount = group.MemberCount
				msgInfo.GroupType = group.GroupType
			}
			if send, ok := sendMap[chatLog.SendID]; ok {
				msgInfo.SenderName = send.Nickname
			}
		}

		pbchatLog.ConversationID = conversationID
		pbchatLog.MsgInfo = msgInfo
		conversationMsg[conversationID] = pbchatLog
	}

	return conversationMsg, nil
}
