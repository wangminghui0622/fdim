package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendApplyUnhandledCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendApplyUnhandledCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendApplyUnhandledCountLogic {
	return &GetFriendApplyUnhandledCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendApplyUnhandledCountLogic) GetFriendApplyUnhandledCount(req *user.GetSelfUnhandledApplyCountReq) (*user.GetSelfUnhandledApplyCountResp, error) {
	resp := &user.GetSelfUnhandledApplyCountResp{}

	// ??????
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// ???????????????????
	count, err := l.svcCtx.FriendDB.GetUnhandledFriendRequestCount(l.ctx, req.UserID, req.Time)
	if err != nil {
		return nil, err
	}

	resp.Count = int64(count)
	return resp, nil
}
