package server

import (
	"context"

	"fdim/protocol/conversation"
	"fdim/rpc/conversation/internal/logic"
	"fdim/rpc/conversation/internal/svc"
)

type ConversationServer struct {
	conversation.UnimplementedConversationServer
	svcCtx *svc.ServiceContext
}

func NewConversationServer(svcCtx *svc.ServiceContext) *ConversationServer {
	return &ConversationServer{
		svcCtx: svcCtx,
	}
}

// GetAllConversations 获取所有会话
func (s *ConversationServer) GetAllConversations(ctx context.Context, req *conversation.GetAllConversationsReq) (*conversation.GetAllConversationsResp, error) {
	l := logic.NewGetAllConversationsLogic(ctx, s.svcCtx)
	return l.GetAllConversations(req)
}

// GetSortedConversationList 获取排序的会话列表
func (s *ConversationServer) GetSortedConversationList(ctx context.Context, req *conversation.GetSortedConversationListReq) (*conversation.GetSortedConversationListResp, error) {
	l := logic.NewGetSortedConversationListLogic(ctx, s.svcCtx)
	return l.GetSortedConversationList(req)
}

// GetConversation 获取会话
func (s *ConversationServer) GetConversation(ctx context.Context, req *conversation.GetConversationReq) (*conversation.GetConversationResp, error) {
	l := logic.NewGetConversationLogic(ctx, s.svcCtx)
	return l.GetConversation(req)
}

// GetConversations 获取多个会话
func (s *ConversationServer) GetConversations(ctx context.Context, req *conversation.GetConversationsReq) (*conversation.GetConversationsResp, error) {
	l := logic.NewGetConversationsLogic(ctx, s.svcCtx)
	return l.GetConversations(req)
}

// SetConversation 设置会话
func (s *ConversationServer) SetConversation(ctx context.Context, req *conversation.SetConversationReq) (*conversation.SetConversationResp, error) {
	l := logic.NewSetConversationLogic(ctx, s.svcCtx)
	return l.SetConversation(req)
}

// GetRecvMsgNotNotifyUserIDs 获取不接收消息通知的用户ID列表
func (s *ConversationServer) GetRecvMsgNotNotifyUserIDs(ctx context.Context, req *conversation.GetRecvMsgNotNotifyUserIDsReq) (*conversation.GetRecvMsgNotNotifyUserIDsResp, error) {
	l := logic.NewGetRecvMsgNotNotifyUserIDsLogic(ctx, s.svcCtx)
	return l.GetRecvMsgNotNotifyUserIDs(req)
}

// CreateSingleChatConversations 创建单聊会话
func (s *ConversationServer) CreateSingleChatConversations(ctx context.Context, req *conversation.CreateSingleChatConversationsReq) (*conversation.CreateSingleChatConversationsResp, error) {
	l := logic.NewCreateSingleChatConversationsLogic(ctx, s.svcCtx)
	return l.CreateSingleChatConversations(req)
}

// CreateGroupChatConversations 创建群聊会话
func (s *ConversationServer) CreateGroupChatConversations(ctx context.Context, req *conversation.CreateGroupChatConversationsReq) (*conversation.CreateGroupChatConversationsResp, error) {
	l := logic.NewCreateGroupChatConversationsLogic(ctx, s.svcCtx)
	return l.CreateGroupChatConversations(req)
}

// SetConversationMaxSeq 设置会话最大序列号
func (s *ConversationServer) SetConversationMaxSeq(ctx context.Context, req *conversation.SetConversationMaxSeqReq) (*conversation.SetConversationMaxSeqResp, error) {
	l := logic.NewSetConversationMaxSeqLogic(ctx, s.svcCtx)
	return l.SetConversationMaxSeq(req)
}

// SetConversationMinSeq 设置会话最小序列号
func (s *ConversationServer) SetConversationMinSeq(ctx context.Context, req *conversation.SetConversationMinSeqReq) (*conversation.SetConversationMinSeqResp, error) {
	l := logic.NewSetConversationMinSeqLogic(ctx, s.svcCtx)
	return l.SetConversationMinSeq(req)
}

// GetConversationIDs 获取会话ID列表
func (s *ConversationServer) GetConversationIDs(ctx context.Context, req *conversation.GetConversationIDsReq) (*conversation.GetConversationIDsResp, error) {
	l := logic.NewGetConversationIDsLogic(ctx, s.svcCtx)
	return l.GetConversationIDs(req)
}

// SetConversations 设置多个会话
func (s *ConversationServer) SetConversations(ctx context.Context, req *conversation.SetConversationsReq) (*conversation.SetConversationsResp, error) {
	l := logic.NewSetConversationsLogic(ctx, s.svcCtx)
	return l.SetConversations(req)
}

// GetUserConversationIDsHash 获取用户会话ID哈希
func (s *ConversationServer) GetUserConversationIDsHash(ctx context.Context, req *conversation.GetUserConversationIDsHashReq) (*conversation.GetUserConversationIDsHashResp, error) {
	l := logic.NewGetUserConversationIDsHashLogic(ctx, s.svcCtx)
	return l.GetUserConversationIDsHash(req)
}

// GetConversationsByConversationID 根据会话ID获取会话
func (s *ConversationServer) GetConversationsByConversationID(ctx context.Context, req *conversation.GetConversationsByConversationIDReq) (*conversation.GetConversationsByConversationIDResp, error) {
	l := logic.NewGetConversationsByConversationIDLogic(ctx, s.svcCtx)
	return l.GetConversationsByConversationID(req)
}

// GetConversationOfflinePushUserIDs 获取会话离线推送用户ID列表
func (s *ConversationServer) GetConversationOfflinePushUserIDs(ctx context.Context, req *conversation.GetConversationOfflinePushUserIDsReq) (*conversation.GetConversationOfflinePushUserIDsResp, error) {
	l := logic.NewGetConversationOfflinePushUserIDsLogic(ctx, s.svcCtx)
	return l.GetConversationOfflinePushUserIDs(req)
}

// GetConversationNotReceiveMessageUserIDs 获取不接收消息的用户ID列表
func (s *ConversationServer) GetConversationNotReceiveMessageUserIDs(ctx context.Context, req *conversation.GetConversationNotReceiveMessageUserIDsReq) (*conversation.GetConversationNotReceiveMessageUserIDsResp, error) {
	l := logic.NewGetConversationNotReceiveMessageUserIDsLogic(ctx, s.svcCtx)
	return l.GetConversationNotReceiveMessageUserIDs(req)
}

// UpdateConversation 更新会话
func (s *ConversationServer) UpdateConversation(ctx context.Context, req *conversation.UpdateConversationReq) (*conversation.UpdateConversationResp, error) {
	l := logic.NewUpdateConversationLogic(ctx, s.svcCtx)
	return l.UpdateConversation(req)
}

// GetFullOwnerConversationIDs 获取完整拥有者会话ID列表
func (s *ConversationServer) GetFullOwnerConversationIDs(ctx context.Context, req *conversation.GetFullOwnerConversationIDsReq) (*conversation.GetFullOwnerConversationIDsResp, error) {
	l := logic.NewGetFullOwnerConversationIDsLogic(ctx, s.svcCtx)
	return l.GetFullOwnerConversationIDs(req)
}

// GetIncrementalConversation 获取增量会话
func (s *ConversationServer) GetIncrementalConversation(ctx context.Context, req *conversation.GetIncrementalConversationReq) (*conversation.GetIncrementalConversationResp, error) {
	l := logic.NewGetIncrementalConversationLogic(ctx, s.svcCtx)
	return l.GetIncrementalConversation(req)
}

// GetOwnerConversation 获取拥有者会话
func (s *ConversationServer) GetOwnerConversation(ctx context.Context, req *conversation.GetOwnerConversationReq) (*conversation.GetOwnerConversationResp, error) {
	l := logic.NewGetOwnerConversationLogic(ctx, s.svcCtx)
	return l.GetOwnerConversation(req)
}

// GetConversationsNeedClearMsg 获取需要清除消息的会话
func (s *ConversationServer) GetConversationsNeedClearMsg(ctx context.Context, req *conversation.GetConversationsNeedClearMsgReq) (*conversation.GetConversationsNeedClearMsgResp, error) {
	l := logic.NewGetConversationsNeedClearMsgLogic(ctx, s.svcCtx)
	return l.GetConversationsNeedClearMsg(req)
}

// GetNotNotifyConversationIDs 获取不通知的会话ID列表
func (s *ConversationServer) GetNotNotifyConversationIDs(ctx context.Context, req *conversation.GetNotNotifyConversationIDsReq) (*conversation.GetNotNotifyConversationIDsResp, error) {
	l := logic.NewGetNotNotifyConversationIDsLogic(ctx, s.svcCtx)
	return l.GetNotNotifyConversationIDs(req)
}

// GetPinnedConversationIDs 获取置顶的会话ID列表
func (s *ConversationServer) GetPinnedConversationIDs(ctx context.Context, req *conversation.GetPinnedConversationIDsReq) (*conversation.GetPinnedConversationIDsResp, error) {
	l := logic.NewGetPinnedConversationIDsLogic(ctx, s.svcCtx)
	return l.GetPinnedConversationIDs(req)
}

// ClearUserConversationMsg 清除用户会话消息
func (s *ConversationServer) ClearUserConversationMsg(ctx context.Context, req *conversation.ClearUserConversationMsgReq) (*conversation.ClearUserConversationMsgResp, error) {
	l := logic.NewClearUserConversationMsgLogic(ctx, s.svcCtx)
	return l.ClearUserConversationMsg(req)
}

// UpdateConversationsByUser 按用户更新会话
func (s *ConversationServer) UpdateConversationsByUser(ctx context.Context, req *conversation.UpdateConversationsByUserReq) (*conversation.UpdateConversationsByUserResp, error) {
	l := logic.NewUpdateConversationsByUserLogic(ctx, s.svcCtx)
	return l.UpdateConversationsByUser(req)
}

// DeleteConversations 删除会话
func (s *ConversationServer) DeleteConversations(ctx context.Context, req *conversation.DeleteConversationsReq) (*conversation.DeleteConversationsResp, error) {
	l := logic.NewDeleteConversationsLogic(ctx, s.svcCtx)
	return l.DeleteConversations(req)
}
