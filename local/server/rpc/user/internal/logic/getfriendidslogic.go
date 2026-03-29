package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendIDsLogic {
	return &GetFriendIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendIDsLogic) GetFriendIDs(req *user.GetFriendIDsReq) (*user.GetFriendIDsResp, error) {
	resp := &user.GetFriendIDsResp{}

	// ??????
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	friendIDs, err := l.svcCtx.FriendDB.FindFriendUserIDs(l.ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	resp.FriendIDs = friendIDs
	return resp, nil
}
