package server

import (
	"context"
	logic2 "fdim/Infrastructure_service/push/internal/logic"
	"fdim/Infrastructure_service/push/internal/svc"

	pbpush "fdim/protocol/push"
)

type PushServer struct {
	pbpush.UnimplementedPushMsgServiceServer
	svcCtx *svc.ServiceContext
}

func NewPushServer(svcCtx *svc.ServiceContext) *PushServer {
	return &PushServer{
		svcCtx: svcCtx,
	}
}

// PushMsg 推送消?
func (s *PushServer) PushMsg(ctx context.Context, req *pbpush.PushMsgReq) (*pbpush.PushMsgResp, error) {
	l := logic2.NewPushMsgLogic(ctx, s.svcCtx)
	return l.PushMsg(req)
}

// DelUserPushToken 删除用户推送Token
func (s *PushServer) DelUserPushToken(ctx context.Context, req *pbpush.DelUserPushTokenReq) (*pbpush.DelUserPushTokenResp, error) {
	l := logic2.NewDelUserPushTokenLogic(ctx, s.svcCtx)
	return l.DelUserPushToken(req)
}
