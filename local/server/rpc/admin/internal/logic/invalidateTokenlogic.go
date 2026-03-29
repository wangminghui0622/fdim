package logic

import (
	"context"
	"fmt"

	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type InvalidateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInvalidateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InvalidateTokenLogic {
	return &InvalidateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InvalidateTokenLogic) InvalidateToken(req *admin.InvalidateTokenReq) (*admin.InvalidateTokenResp, error) {
	// ʹָû Token ʧЧ
	if err := l.svcCtx.AdminDB.InvalidateToken(l.ctx, req.UserID); err != nil {
		l.Errorf("InvalidateToken failed: %v", err)
		return nil, fmt.Errorf("failed to invalidate token: %w", err)
	}

	return &admin.InvalidateTokenResp{}, nil
}
