package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetGlobalRecvMessageOptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetGlobalRecvMessageOptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGlobalRecvMessageOptLogic {
	return &SetGlobalRecvMessageOptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetGlobalRecvMessageOptLogic) SetGlobalRecvMessageOpt(req *user.SetGlobalRecvMessageOptReq) (*user.SetGlobalRecvMessageOptResp, error) {
	resp := &user.SetGlobalRecvMessageOptResp{}

	// Ȩ֤
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// ֤
	if req.UserID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	// ûǷ
	_, err := l.svcCtx.UserDB.Take(l.ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// ȫֽϢѡ
	data := map[string]interface{}{
		"global_recv_msg_opt": req.GlobalRecvMsgOpt,
	}
	if err := l.svcCtx.UserDB.UpdateByMap(l.ctx, req.UserID, data); err != nil {
		return nil, err
	}

	// ûϢ֪֪ͨͨûԼ
	if l.svcCtx.UserNotification != nil {
		l.svcCtx.UserNotification.UserInfoUpdatedNotification(l.ctx, req.UserID, req.UserID)
	}

	return resp, nil
}
