package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelDefaultFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelDefaultFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelDefaultFriendLogic {
	return &DelDefaultFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelDefaultFriendLogic) DelDefaultFriend(req *admin.DelDefaultFriendReq) (*admin.DelDefaultFriendResp, error) {
	// 1. 验证参数
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("userIDs cannot be empty")
	}

	// 2. 检查是否存在
	exists, err := l.svcCtx.AdminDB.FindDefaultFriend(l.ctx, req.UserIDs)
	if err != nil {
		l.Errorf("FindDefaultFriend failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to check default friends")
	}
	if len(exists) == 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg("no default friends found for the given userIDs")
	}

	// 3. 删除默认好友
	if err := l.svcCtx.AdminDB.DelDefaultFriend(l.ctx, req.UserIDs); err != nil {
		l.Errorf("DelDefaultFriend failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to delete default friends")
	}

	l.Infof("Deleted default friends: userIDs=%v", req.UserIDs)
	return &admin.DelDefaultFriendResp{}, nil
}
