package server

import (
	"context"

	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/logic"
	"fdim/rpc/msg/internal/svc"
)

type MsgServer struct {
	msg.UnimplementedMsgServer
	svcCtx *svc.ServiceContext
}

func NewMsgServer(svcCtx *svc.ServiceContext) *MsgServer {
	return &MsgServer{
		svcCtx: svcCtx,
	}
}

// GetMaxSeq 获取最大序列号
func (s *MsgServer) GetMaxSeq(ctx context.Context, req *sdkws.GetMaxSeqReq) (*sdkws.GetMaxSeqResp, error) {
	l := logic.NewGetMaxSeqLogic(ctx, s.svcCtx)
	return l.GetMaxSeq(req)
}

// GetMaxSeqs 获取多个最大序列号
func (s *MsgServer) GetMaxSeqs(ctx context.Context, req *msg.GetMaxSeqsReq) (*msg.SeqsInfoResp, error) {
	l := logic.NewGetMaxSeqsLogic(ctx, s.svcCtx)
	return l.GetMaxSeqs(req)
}

// GetHasReadSeqs 获取已读序列号
func (s *MsgServer) GetHasReadSeqs(ctx context.Context, req *msg.GetHasReadSeqsReq) (*msg.SeqsInfoResp, error) {
	l := logic.NewGetHasReadSeqsLogic(ctx, s.svcCtx)
	return l.GetHasReadSeqs(req)
}

// GetMsgByConversationIDs 根据会话ID获取消息
func (s *MsgServer) GetMsgByConversationIDs(ctx context.Context, req *msg.GetMsgByConversationIDsReq) (*msg.GetMsgByConversationIDsResp, error) {
	l := logic.NewGetMsgByConversationIDsLogic(ctx, s.svcCtx)
	return l.GetMsgByConversationIDs(req)
}

// GetConversationMaxSeq 获取会话最大序列号
func (s *MsgServer) GetConversationMaxSeq(ctx context.Context, req *msg.GetConversationMaxSeqReq) (*msg.GetConversationMaxSeqResp, error) {
	l := logic.NewGetConversationMaxSeqLogic(ctx, s.svcCtx)
	return l.GetConversationMaxSeq(req)
}

// PullMessageBySeqs 根据序列号拉取消息
func (s *MsgServer) PullMessageBySeqs(ctx context.Context, req *sdkws.PullMessageBySeqsReq) (*sdkws.PullMessageBySeqsResp, error) {
	l := logic.NewPullMessageBySeqsLogic(ctx, s.svcCtx)
	return l.PullMessageBySeqs(req)
}

// GetSeqMessage 获取序列号消息
func (s *MsgServer) GetSeqMessage(ctx context.Context, req *msg.GetSeqMessageReq) (*msg.GetSeqMessageResp, error) {
	l := logic.NewGetSeqMessageLogic(ctx, s.svcCtx)
	return l.GetSeqMessage(req)
}

// SearchMessage 搜索消息
func (s *MsgServer) SearchMessage(ctx context.Context, req *msg.SearchMessageReq) (*msg.SearchMessageResp, error) {
	l := logic.NewSearchMessageLogic(ctx, s.svcCtx)
	return l.SearchMessage(req)
}

// SendMsg 发送消息
func (s *MsgServer) SendMsg(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error) {
	l := logic.NewSendMsgLogic(ctx, s.svcCtx)
	return l.SendMsg(req)
}

// SendSimpleMsg 发送简单消息
func (s *MsgServer) SendSimpleMsg(ctx context.Context, req *msg.SendSimpleMsgReq) (*msg.SendSimpleMsgResp, error) {
	l := logic.NewSendSimpleMsgLogic(ctx, s.svcCtx)
	return l.SendSimpleMsg(req)
}

// SetUserConversationsMinSeq 设置用户会话最小序列号
func (s *MsgServer) SetUserConversationsMinSeq(ctx context.Context, req *msg.SetUserConversationsMinSeqReq) (*msg.SetUserConversationsMinSeqResp, error) {
	l := logic.NewSetUserConversationsMinSeqLogic(ctx, s.svcCtx)
	return l.SetUserConversationsMinSeq(req)
}

// ClearConversationsMsg 清除会话消息
func (s *MsgServer) ClearConversationsMsg(ctx context.Context, req *msg.ClearConversationsMsgReq) (*msg.ClearConversationsMsgResp, error) {
	l := logic.NewClearConversationsMsgLogic(ctx, s.svcCtx)
	return l.ClearConversationsMsg(req)
}

// UserClearAllMsg 用户清除所有消息
func (s *MsgServer) UserClearAllMsg(ctx context.Context, req *msg.UserClearAllMsgReq) (*msg.UserClearAllMsgResp, error) {
	l := logic.NewUserClearAllMsgLogic(ctx, s.svcCtx)
	return l.UserClearAllMsg(req)
}

// DeleteMsgs 删除消息
func (s *MsgServer) DeleteMsgs(ctx context.Context, req *msg.DeleteMsgsReq) (*msg.DeleteMsgsResp, error) {
	l := logic.NewDeleteMsgsLogic(ctx, s.svcCtx)
	return l.DeleteMsgs(req)
}

// DeleteMsgPhysicalBySeq 根据序列号物理删除消息
func (s *MsgServer) DeleteMsgPhysicalBySeq(ctx context.Context, req *msg.DeleteMsgPhysicalBySeqReq) (*msg.DeleteMsgPhysicalBySeqResp, error) {
	l := logic.NewDeleteMsgPhysicalBySeqLogic(ctx, s.svcCtx)
	return l.DeleteMsgPhysicalBySeq(req)
}

// DeleteMsgPhysical 物理删除消息
func (s *MsgServer) DeleteMsgPhysical(ctx context.Context, req *msg.DeleteMsgPhysicalReq) (*msg.DeleteMsgPhysicalResp, error) {
	l := logic.NewDeleteMsgPhysicalLogic(ctx, s.svcCtx)
	return l.DeleteMsgPhysical(req)
}

// SetSendMsgStatus 设置发送消息状态
func (s *MsgServer) SetSendMsgStatus(ctx context.Context, req *msg.SetSendMsgStatusReq) (*msg.SetSendMsgStatusResp, error) {
	l := logic.NewSetSendMsgStatusLogic(ctx, s.svcCtx)
	return l.SetSendMsgStatus(req)
}

// GetSendMsgStatus 获取发送消息状态
func (s *MsgServer) GetSendMsgStatus(ctx context.Context, req *msg.GetSendMsgStatusReq) (*msg.GetSendMsgStatusResp, error) {
	l := logic.NewGetSendMsgStatusLogic(ctx, s.svcCtx)
	return l.GetSendMsgStatus(req)
}

// RevokeMsg 撤回消息
func (s *MsgServer) RevokeMsg(ctx context.Context, req *msg.RevokeMsgReq) (*msg.RevokeMsgResp, error) {
	l := logic.NewRevokeMsgLogic(ctx, s.svcCtx)
	return l.RevokeMsg(req)
}

// MarkMsgsAsRead 标记消息为已读
func (s *MsgServer) MarkMsgsAsRead(ctx context.Context, req *msg.MarkMsgsAsReadReq) (*msg.MarkMsgsAsReadResp, error) {
	l := logic.NewMarkMsgsAsReadLogic(ctx, s.svcCtx)
	return l.MarkMsgsAsRead(req)
}

// MarkConversationAsRead 标记会话为已读
func (s *MsgServer) MarkConversationAsRead(ctx context.Context, req *msg.MarkConversationAsReadReq) (*msg.MarkConversationAsReadResp, error) {
	l := logic.NewMarkConversationAsReadLogic(ctx, s.svcCtx)
	return l.MarkConversationAsRead(req)
}

// SetConversationHasReadSeq 设置会话已读序列号
func (s *MsgServer) SetConversationHasReadSeq(ctx context.Context, req *msg.SetConversationHasReadSeqReq) (*msg.SetConversationHasReadSeqResp, error) {
	l := logic.NewSetConversationHasReadSeqLogic(ctx, s.svcCtx)
	return l.SetConversationHasReadSeq(req)
}

// GetConversationsHasReadAndMaxSeq 获取会话已读和最大序列号
func (s *MsgServer) GetConversationsHasReadAndMaxSeq(ctx context.Context, req *msg.GetConversationsHasReadAndMaxSeqReq) (*msg.GetConversationsHasReadAndMaxSeqResp, error) {
	l := logic.NewGetConversationsHasReadAndMaxSeqLogic(ctx, s.svcCtx)
	return l.GetConversationsHasReadAndMaxSeq(req)
}

// GetActiveUser 获取活跃用户
func (s *MsgServer) GetActiveUser(ctx context.Context, req *msg.GetActiveUserReq) (*msg.GetActiveUserResp, error) {
	l := logic.NewGetActiveUserLogic(ctx, s.svcCtx)
	return l.GetActiveUser(req)
}

// GetActiveGroup 获取活跃群组
func (s *MsgServer) GetActiveGroup(ctx context.Context, req *msg.GetActiveGroupReq) (*msg.GetActiveGroupResp, error) {
	l := logic.NewGetActiveGroupLogic(ctx, s.svcCtx)
	return l.GetActiveGroup(req)
}

// GetServerTime 获取服务器时间
func (s *MsgServer) GetServerTime(ctx context.Context, req *msg.GetServerTimeReq) (*msg.GetServerTimeResp, error) {
	l := logic.NewGetServerTimeLogic(ctx, s.svcCtx)
	return l.GetServerTime(req)
}

// ClearMsg 清除消息
func (s *MsgServer) ClearMsg(ctx context.Context, req *msg.ClearMsgReq) (*msg.ClearMsgResp, error) {
	l := logic.NewClearMsgLogic(ctx, s.svcCtx)
	return l.ClearMsg(req)
}

// DestructMsgs 销毁消息
func (s *MsgServer) DestructMsgs(ctx context.Context, req *msg.DestructMsgsReq) (*msg.DestructMsgsResp, error) {
	l := logic.NewDestructMsgsLogic(ctx, s.svcCtx)
	return l.DestructMsgs(req)
}

// GetActiveConversation 获取活跃会话
func (s *MsgServer) GetActiveConversation(ctx context.Context, req *msg.GetActiveConversationReq) (*msg.GetActiveConversationResp, error) {
	l := logic.NewGetActiveConversationLogic(ctx, s.svcCtx)
	return l.GetActiveConversation(req)
}

// SetUserConversationMaxSeq 设置用户会话最大序列号
func (s *MsgServer) SetUserConversationMaxSeq(ctx context.Context, req *msg.SetUserConversationMaxSeqReq) (*msg.SetUserConversationMaxSeqResp, error) {
	l := logic.NewSetUserConversationMaxSeqLogic(ctx, s.svcCtx)
	return l.SetUserConversationMaxSeq(req)
}

// SetUserConversationMinSeq 设置用户会话最小序列号
func (s *MsgServer) SetUserConversationMinSeq(ctx context.Context, req *msg.SetUserConversationMinSeqReq) (*msg.SetUserConversationMinSeqResp, error) {
	l := logic.NewSetUserConversationMinSeqLogic(ctx, s.svcCtx)
	return l.SetUserConversationMinSeq(req)
}

// GetLastMessageSeqByTime 根据时间获取最后消息序列号
func (s *MsgServer) GetLastMessageSeqByTime(ctx context.Context, req *msg.GetLastMessageSeqByTimeReq) (*msg.GetLastMessageSeqByTimeResp, error) {
	l := logic.NewGetLastMessageSeqByTimeLogic(ctx, s.svcCtx)
	return l.GetLastMessageSeqByTime(req)
}

// GetLastMessage 获取最后消息
func (s *MsgServer) GetLastMessage(ctx context.Context, req *msg.GetLastMessageReq) (*msg.GetLastMessageResp, error) {
	l := logic.NewGetLastMessageLogic(ctx, s.svcCtx)
	return l.GetLastMessage(req)
}
