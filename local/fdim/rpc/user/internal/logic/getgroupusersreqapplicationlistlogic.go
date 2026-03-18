package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupUsersReqApplicationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupUsersReqApplicationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupUsersReqApplicationListLogic {
	return &GetGroupUsersReqApplicationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupUsersReqApplicationListLogic) GetGroupUsersReqApplicationList(req *user.GetGroupUsersReqApplicationListReq) (*user.GetGroupUsersReqApplicationListResp, error) {
	// TODO: 实现获取群组用户申请列表的逻辑
	return nil, fmt.Errorf("GetGroupUsersReqApplicationList not fully implemented yet")
}
