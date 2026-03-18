package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ========== JSSdkGetConversations ==========

type JSSdkGetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJSSdkGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JSSdkGetConversationsLogic {
	return &JSSdkGetConversationsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *JSSdkGetConversationsLogic) JSSdkGetConversations(req *types.JSSdkGetConversationsReq) (*types.JSSdkGetConversationsResp, error) {
	return &types.JSSdkGetConversationsResp{Total: 0, Conversations: []interface{}{}}, nil
}

// ========== JSSdkGetActiveConversations ==========

type JSSdkGetActiveConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJSSdkGetActiveConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JSSdkGetActiveConversationsLogic {
	return &JSSdkGetActiveConversationsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *JSSdkGetActiveConversationsLogic) JSSdkGetActiveConversations(req *types.JSSdkGetActiveConversationsReq) (*types.JSSdkGetActiveConversationsResp, error) {
	return &types.JSSdkGetActiveConversationsResp{Conversations: []interface{}{}}, nil
}
