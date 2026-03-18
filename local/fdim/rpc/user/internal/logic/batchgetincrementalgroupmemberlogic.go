package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetIncrementalGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetIncrementalGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetIncrementalGroupMemberLogic {
	return &BatchGetIncrementalGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetIncrementalGroupMemberLogic) BatchGetIncrementalGroupMember(req *user.BatchGetIncrementalGroupMemberReq) (*user.BatchGetIncrementalGroupMemberResp, error) {
	// TODO: 实现批量获取增量群成员的逻辑
	return nil, fmt.Errorf("BatchGetIncrementalGroupMember not fully implemented yet")
}
