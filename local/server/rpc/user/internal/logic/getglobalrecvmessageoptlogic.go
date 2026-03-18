package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGlobalRecvMessageOptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGlobalRecvMessageOptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGlobalRecvMessageOptLogic {
	return &GetGlobalRecvMessageOptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGlobalRecvMessageOptLogic) GetGlobalRecvMessageOpt(req *user.GetGlobalRecvMessageOptReq) (*user.GetGlobalRecvMessageOptResp, error) {
	resp := &user.GetGlobalRecvMessageOptResp{}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// ������֤
	if req.UserID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	// �����û�
	user, err := l.svcCtx.UserDB.Take(l.ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	resp.GlobalRecvMsgOpt = user.GlobalRecvMsgOpt
	return resp, nil
}
