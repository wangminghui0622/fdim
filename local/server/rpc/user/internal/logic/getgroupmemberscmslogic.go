package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMembersCMSLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMembersCMSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMembersCMSLogic {
	return &GetGroupMembersCMSLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMembersCMSLogic) GetGroupMembersCMS(req *user.GetGroupMembersCMSReq) (*user.GetGroupMembersCMSResp, error) {
	// TODO: ʵֻȡȺԱCMS߼
	return nil, fmt.Errorf("GetGroupMembersCMS not fully implemented yet")
}
