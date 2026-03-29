package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSpecifiedUserGroupRequestInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSpecifiedUserGroupRequestInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSpecifiedUserGroupRequestInfoLogic {
	return &GetSpecifiedUserGroupRequestInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSpecifiedUserGroupRequestInfoLogic) GetSpecifiedUserGroupRequestInfo(req *user.GetSpecifiedUserGroupRequestInfoReq) (*user.GetSpecifiedUserGroupRequestInfoResp, error) {
	// TODO: ʵֻȡָûȺϢ߼
	return nil, fmt.Errorf("GetSpecifiedUserGroupRequestInfo not fully implemented yet")
}
