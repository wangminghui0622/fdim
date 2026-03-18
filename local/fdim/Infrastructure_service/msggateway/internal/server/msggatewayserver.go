package server

import (
	"context"
	logic2 "fdim/Infrastructure_service/msggateway/internal/logic"
	"fdim/Infrastructure_service/msggateway/internal/svc"

	"fdim/protocol/msggateway"
)

type MsgGatewayServer struct {
	msggateway.UnimplementedMsgGatewayServer
	svcCtx *svc.ServiceContext
}

func NewMsgGatewayServer(svcCtx *svc.ServiceContext) *MsgGatewayServer {
	return &MsgGatewayServer{
		svcCtx: svcCtx,
	}
}

// OnlinePushMsg 在线推送消息
func (s *MsgGatewayServer) OnlinePushMsg(ctx context.Context, req *msggateway.OnlinePushMsgReq) (*msggateway.OnlinePushMsgResp, error) {
	l := logic2.NewOnlinePushMsgLogic(ctx, s.svcCtx)
	return l.OnlinePushMsg(req)
}

// GetUsersOnlineStatus 获取用户在线状态
func (s *MsgGatewayServer) GetUsersOnlineStatus(ctx context.Context, req *msggateway.GetUsersOnlineStatusReq) (*msggateway.GetUsersOnlineStatusResp, error) {
	l := logic2.NewGetUsersOnlineStatusLogic(ctx, s.svcCtx)
	return l.GetUsersOnlineStatus(req)
}

// OnlineBatchPushOneMsg 在线批量推送单条消息
func (s *MsgGatewayServer) OnlineBatchPushOneMsg(ctx context.Context, req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	l := logic2.NewOnlineBatchPushOneMsgLogic(ctx, s.svcCtx)
	return l.OnlineBatchPushOneMsg(req)
}

// SuperGroupOnlineBatchPushOneMsg 超级群在线批量推送单条消息
func (s *MsgGatewayServer) SuperGroupOnlineBatchPushOneMsg(ctx context.Context, req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	l := logic2.NewSuperGroupOnlineBatchPushOneMsgLogic(ctx, s.svcCtx)
	return l.SuperGroupOnlineBatchPushOneMsg(req)
}

// KickUserOffline 踢用户下线
func (s *MsgGatewayServer) KickUserOffline(ctx context.Context, req *msggateway.KickUserOfflineReq) (*msggateway.KickUserOfflineResp, error) {
	l := logic2.NewKickUserOfflineLogic(ctx, s.svcCtx)
	return l.KickUserOffline(req)
}

// MultiTerminalLoginCheck 多端登录检查
func (s *MsgGatewayServer) MultiTerminalLoginCheck(ctx context.Context, req *msggateway.MultiTerminalLoginCheckReq) (*msggateway.MultiTerminalLoginCheckResp, error) {
	l := logic2.NewMultiTerminalLoginCheckLogic(ctx, s.svcCtx)
	return l.MultiTerminalLoginCheck(req)
}
