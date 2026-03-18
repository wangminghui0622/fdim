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

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// ��ȡ����ID�б�
	friendIDs, err := l.svcCtx.FriendDB.FindFriendUserIDs(l.ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// GetFriendList ֻ���غ���ID�б����������ϸ��Ϣ
	// ��ϸ��Ϣ����ͨ�� GetSpecifiedFriends �� GetPaginationFriends ��ȡ
	resp.FriendIDs = friendIDs
	return resp, nil
}
