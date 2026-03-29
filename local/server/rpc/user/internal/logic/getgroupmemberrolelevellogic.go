package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberRoleLevelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberRoleLevelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberRoleLevelLogic {
	return &GetGroupMemberRoleLevelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMemberRoleLevelLogic) GetGroupMemberRoleLevel(req *user.GetGroupMemberRoleLevelReq) (*user.GetGroupMemberRoleLevelResp, error) {
	// TODO: ʵֻȡȺԱɫ߼
	return nil, fmt.Errorf("GetGroupMemberRoleLevel not fully implemented yet")
}
