package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindDefaultFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindDefaultFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindDefaultFriendLogic {
	return &FindDefaultFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindDefaultFriendLogic) FindDefaultFriend(req *admin.FindDefaultFriendReq) (*admin.FindDefaultFriendResp, error) {
	// FindDefaultFriendReq û UserIDs ֶΣĬϺ
	userIDs, err := l.svcCtx.AdminDB.FindDefaultFriend(l.ctx, nil)
	if err != nil {
		l.Errorf("FindDefaultFriend failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to find default friends")
	}

	return &admin.FindDefaultFriendResp{
		UserIDs: userIDs,
	}, nil
}
