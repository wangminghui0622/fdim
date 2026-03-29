package logic

import (
	"context"
	"fmt"

	"fdim/protocol/auth"
	"fdim/rpc/auth/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ParseTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewParseTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseTokenLogic {
	return &ParseTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ParseTokenLogic) ParseToken(req *auth.ParseTokenReq) (*auth.ParseTokenResp, error) {
	if l.svcCtx.Token == nil {
		return nil, fmt.Errorf("token verifier not initialized")
	}

	// 1. 解析 Token，获?userID ?userType
	userID, _, err := l.svcCtx.Token.GetToken(req.Token)
	if err != nil {
		return nil, err
	}

	// 当前 tokenverify.Token 不包?platformID 和到期时间的解析能力（简化版），这里?proto 返回 0?
	return &auth.ParseTokenResp{
		UserID:            userID,
		PlatformID:        0,
		ExpireTimeSeconds: 0,
	}, nil
}
