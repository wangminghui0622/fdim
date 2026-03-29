package logic

import (
	"context"
	"fmt"

	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserTokenLogic {
	return &GetUserTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserTokenLogic) GetUserToken(req *admin.GetUserTokenReq) (*admin.GetUserTokenResp, error) {
	// ȡû Token
	tokensMap, err := l.svcCtx.AdminDB.GetTokens(l.ctx, req.UserID)
	if err != nil {
		l.Errorf("GetTokens failed: %v", err)
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}

	return &admin.GetUserTokenResp{
		TokensMap: tokensMap,
	}, nil
}
