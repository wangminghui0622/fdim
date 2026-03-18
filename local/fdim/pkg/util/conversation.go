package util

import (
	"strings"
)

// ParseConversationID 解析会话ID，提取用户ID
// 单聊: si_{userID1}_{userID2} (按字典序排序)
// 群聊: g_{groupID} 或 sg_{groupID}
func ParseConversationID(conversationID string) (sessionType string, userIDs []string, groupID string) {
	if strings.HasPrefix(conversationID, "si_") {
		// 单聊
		parts := strings.Split(conversationID, "_")
		if len(parts) == 3 {
			return "single", []string{parts[1], parts[2]}, ""
		}
	} else if strings.HasPrefix(conversationID, "g_") {
		// 可写群聊
		groupID = strings.TrimPrefix(conversationID, "g_")
		return "group", nil, groupID
	} else if strings.HasPrefix(conversationID, "sg_") {
		// 只读群聊（超级群）
		groupID = strings.TrimPrefix(conversationID, "sg_")
		return "super_group", nil, groupID
	}
	return "", nil, ""
}

// GetUserIDsFromConversationID 从会话ID中提取用户ID列表（仅适用于单聊）
func GetUserIDsFromConversationID(conversationID string) []string {
	_, userIDs, _ := ParseConversationID(conversationID)
	return userIDs
}

// GetGroupIDFromConversationID 从会话ID中提取群组ID（仅适用于群聊）
func GetGroupIDFromConversationID(conversationID string) string {
	_, _, groupID := ParseConversationID(conversationID)
	return groupID
}

// GetRecvUserIDFromMsg 从消息中获取接收者用户ID
// 对于单聊，返回 RecvID
// 对于群聊，需要调用 Group 服务获取群成员列表
func GetRecvUserIDFromMsg(sessionType int32, recvID, groupID string) string {
	if sessionType == 1 { // SingleChatType
		return recvID
	}
	// 群聊需要从 Group 服务获取成员列表，这里只返回空字符串
	return ""
}

// ExtractUserIDsFromMsg 从消息中提取需要推送的用户ID列表
func ExtractUserIDsFromMsg(sessionType int32, sendID, recvID, groupID, conversationID string) []string {
	var userIDs []string
	
	switch sessionType {
	case 1: // SingleChatType - 单聊
		// 单聊时，接收者是 RecvID（排除发送者）
		if recvID != "" && recvID != sendID {
			userIDs = append(userIDs, recvID)
		}
		// 也可以从 conversationID 解析
		if len(userIDs) == 0 {
			parsedUserIDs := GetUserIDsFromConversationID(conversationID)
			for _, uid := range parsedUserIDs {
				if uid != sendID {
					userIDs = append(userIDs, uid)
				}
			}
		}
	case 2, 3: // WriteGroupChatType, ReadGroupChatType - 群聊
		// 群聊需要从 Group 服务获取成员列表
		// 这里返回空，需要调用 Group 服务
		userIDs = []string{}
	default:
		userIDs = []string{}
	}
	
	return userIDs
}
