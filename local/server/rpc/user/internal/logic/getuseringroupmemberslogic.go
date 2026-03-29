package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInGroupMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInGroupMembersLogic {
	return &GetUserInGroupMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInGroupMembersLogic) GetUserInGroupMembers(req *user.GetUserInGroupMembersReq) (*user.GetUserInGroupMembersResp, error) {
	// TODO: ʵֻȡûȺԱеϢ߼
	return nil, fmt.Errorf("GetUserInGroupMembers not fully implemented yet")
}
