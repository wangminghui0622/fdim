package logic

import (
	"context"
	"fmt"

	"fdim/protocol/auth"
	"fdim/rpc/auth/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetExistingTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetExistingTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExistingTokenLogic {
	return &GetExistingTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetExistingTokenLogic) GetExistingToken(req *auth.GetExistingTokenReq) (*auth.GetExistingTokenResp, error) {
	if l.svcCtx.AuthDB == nil {
		return nil, fmt.Errorf("auth database not initialized")
	}

	statusMap, err := l.svcCtx.AuthDB.GetTokens(l.ctx, req.UserID, req.PlatformID)
	if err != nil {
		return nil, err
	}
	// auth.proto 定义?map<string,int32> tokenStates
	resp := &auth.GetExistingTokenResp{TokenStates: make(map[string]int32, len(statusMap))}
	for tk, st := range statusMap {
		resp.TokenStates[tk] = st
	}
	return resp, nil
}
