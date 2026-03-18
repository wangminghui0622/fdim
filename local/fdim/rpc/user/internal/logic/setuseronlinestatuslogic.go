package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetUserOnlineStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUserOnlineStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserOnlineStatusLogic {
	return &SetUserOnlineStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetUserOnlineStatusLogic) SetUserOnlineStatus(req *user.SetUserOnlineStatusReq) (*user.SetUserOnlineStatusResp, error) {
	resp := &user.SetUserOnlineStatusResp{}

	// ������֤
	if len(req.Status) == 0 {
		return nil, fmt.Errorf("status is empty")
	}

	// ���������û�����״̬
	for _, status := range req.Status {
		if status.UserID == "" {
			continue
		}
		if err := l.svcCtx.UserCache.SetUserOnlineStatus(l.ctx, status.UserID, status.Online, status.Offline); err != nil {
			return nil, err
		}
	}

	return resp, nil
}
