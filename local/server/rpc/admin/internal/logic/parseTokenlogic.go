package logic

import (
	"context"
	"fmt"

	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
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

func (l *ParseTokenLogic) ParseToken(req *admin.ParseTokenReq) (*admin.ParseTokenResp, error) {
	// 1.  Token
	userID, userType, err := l.svcCtx.Token.GetToken(req.Token)
	if err != nil {
		l.Errorf("GetToken failed: %v", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// 2.  Redis ֤ Token Ƿ
	tokensMap, err := l.svcCtx.AdminDB.GetTokens(l.ctx, userID)
	if err != nil {
		l.Errorf("GetTokens failed: %v", err)
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}

	// 3.  Token Ƿڻ
	if _, ok := tokensMap[req.Token]; !ok {
		return nil, fmt.Errorf("token not found in cache")
	}

	return &admin.ParseTokenResp{
		UserID:   userID,
		UserType: userType,
	}, nil
}
