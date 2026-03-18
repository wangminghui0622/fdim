package logic

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddDefaultFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddDefaultFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDefaultFriendLogic {
	return &AddDefaultFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddDefaultFriendLogic) AddDefaultFriend(req *admin.AddDefaultFriendReq) (*admin.AddDefaultFriendResp, error) {
	// 1. 验证参数
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("userIDs cannot be empty")
	}

	// 2. 检查是否已存在
	exists, err := l.svcCtx.AdminDB.FindDefaultFriend(l.ctx, req.UserIDs)
	if err != nil {
		l.Errorf("FindDefaultFriend failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to check default friends")
	}
	if len(exists) > 0 {
		return nil, fmt.Errorf("some userIDs already exist as default friends: %v", exists)
	}

	// 3. 添加默认好友
	if err := l.svcCtx.AdminDB.AddDefaultFriend(l.ctx, req.UserIDs); err != nil {
		l.Errorf("AddDefaultFriend failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to add default friends")
	}

	l.Infof("Added default friends: userIDs=%v", req.UserIDs)
	return &admin.AddDefaultFriendResp{}, nil
}
