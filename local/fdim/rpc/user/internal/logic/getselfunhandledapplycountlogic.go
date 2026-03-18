package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSelfUnhandledApplyCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSelfUnhandledApplyCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelfUnhandledApplyCountLogic {
	return &GetSelfUnhandledApplyCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSelfUnhandledApplyCountLogic) GetSelfUnhandledApplyCount(req *user.GetSelfUnhandledApplyCountReq) (*user.GetSelfUnhandledApplyCountResp, error) {
	resp := &user.GetSelfUnhandledApplyCountResp{}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// ��ȡ�Լ����͵�δ������������
	count, err := l.svcCtx.FriendDB.GetSelfUnhandledFriendRequestCount(l.ctx, req.UserID, req.Time)
	if err != nil {
		return nil, err
	}

	resp.Count = int64(count)
	return resp, nil
}
