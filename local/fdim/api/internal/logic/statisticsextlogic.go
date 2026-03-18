package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ========== GetActiveUser ==========

type GetActiveUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetActiveUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetActiveUserLogic {
	return &GetActiveUserLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetActiveUserLogic) GetActiveUser(req *types.GetActiveUserReq) (*types.GetActiveUserResp, error) {
	return &types.GetActiveUserResp{MsgCount: 0, Users: []types.ActiveUser{}}, nil
}

// ========== GroupCreateCount ==========

type GroupCreateCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateCountLogic {
	return &GroupCreateCountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GroupCreateCountLogic) GroupCreateCount(req *types.GroupCreateCountReq) (*types.GroupCreateCountResp, error) {
	return &types.GroupCreateCountResp{Total: 0, Before: 0, Counts: []types.DateCount{}}, nil
}

// ========== GetActiveGroup ==========

type GetActiveGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetActiveGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetActiveGroupLogic {
	return &GetActiveGroupLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetActiveGroupLogic) GetActiveGroup(req *types.GetActiveGroupReq) (*types.GetActiveGroupResp, error) {
	return &types.GetActiveGroupResp{MsgCount: 0, Groups: []types.ActiveGroup{}}, nil
}
