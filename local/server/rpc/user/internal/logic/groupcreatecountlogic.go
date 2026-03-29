package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateCountLogic {
	return &GroupCreateCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupCreateCountLogic) GroupCreateCount(req *user.GroupCreateCountReq) (*user.GroupCreateCountResp, error) {
	resp := &user.GroupCreateCountResp{}

	// ???????????????????
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	// ??????????
	total, err := l.svcCtx.GroupDB.CountGroups(l.ctx)
	if err != nil {
		return nil, err
	}

	// ??????????????????
	var before int64
	if req.Start > 0 {
		before, err = l.svcCtx.GroupDB.CountGroupsBeforeTime(l.ctx, req.Start)
		if err != nil {
			return nil, err
		}
	}

	resp.Total = int64(total)
	resp.Before = int64(before)

	return resp, nil
}
