package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendListLogic) GetFriendList(req *user.GetFriendIDsReq) (*user.GetFriendIDsResp, error) {
	resp := &user.GetFriendIDsResp{}

	// ??????
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// ???????ID??
	friendIDs, err := l.svcCtx.FriendDB.FindFriendUserIDs(l.ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// GetFriendList ????????ID???????????????
	// ????????????? GetSpecifiedFriends ?? GetPaginationFriends ???
	resp.FriendIDs = friendIDs
	return resp, nil
}
