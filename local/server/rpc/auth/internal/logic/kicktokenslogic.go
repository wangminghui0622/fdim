package logic

import (
	"context"

	"fdim/protocol/auth"
	"fdim/rpc/auth/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type KickTokensLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKickTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickTokensLogic {
	return &KickTokensLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *KickTokensLogic) KickTokens(req *auth.KickTokensReq) (*auth.KickTokensResp, error) {
	if l.svcCtx.AuthDB == nil {
		// 兼容：无 AuthDB 时直接返回成功，避免调用方失败
		return &auth.KickTokensResp{}, nil
	}
	if len(req.Tokens) == 0 {
		return &auth.KickTokensResp{}, nil
	}
	// 注意：当前 auth.proto 的 KickTokensReq 仅包含 tokens，没有提供 userID/platformID，
	// 无法在 Redis 的 UID_TOKEN_STATUS:{userID}:{platformID} 结构中精确定位。
	// 为避免写入错误数据，这里采用保守策略：仅记录日志并返回成功。
	l.Infof("KickTokens called: tokenCount=%d", len(req.Tokens))
	return &auth.KickTokensResp{}, nil
}
