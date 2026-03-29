package logic

import (
	"context"
	"fmt"

	"fdim/protocol/auth"
	constantpb "fdim/protocol/constant"
	"fdim/rpc/auth/internal/svc"
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

func (l *InvalidateTokenLogic) InvalidateToken(req *auth.InvalidateTokenReq) (*auth.InvalidateTokenResp, error) {
	if l.svcCtx.AuthDB == nil {
		return nil, fmt.Errorf("auth database not initialized")
	}

	statusMap, err := l.svcCtx.AuthDB.GetTokens(l.ctx, req.UserID, req.PlatformID)
	if err != nil {
		return nil, err
	}

	var toKick []string
	for tk := range statusMap {
		if tk == req.PreservedToken {
			continue
		}
		toKick = append(toKick, tk)
	}

	if err := l.svcCtx.AuthDB.BatchKickTokens(l.ctx, req.UserID, req.PlatformID, toKick); err != nil {
		return nil, err
	}

	// 保留?token 明确标记?Normal
	if req.PreservedToken != "" {
		_ = l.svcCtx.AuthDB.SetTokenStatus(l.ctx, req.UserID, req.PlatformID, req.PreservedToken, constantpb.NormalToken)
	}

	return &auth.InvalidateTokenResp{}, nil
}
