package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/util"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendsLogic {
	return &UpdateFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFriendsLogic) UpdateFriends(req *user.UpdateFriendsReq) (*user.UpdateFriendsResp, error) {
	resp := &user.UpdateFriendsResp{}

	// ������֤
	if len(req.FriendUserIDs) == 0 {
		return nil, fmt.Errorf("friendIDList is empty")
	}

	// ��� friendUserIDs �Ƿ��ظ�
	if util.HasDuplicate(req.FriendUserIDs) {
		return nil, fmt.Errorf("friendIDList repeated")
	}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	// ����Ƿ�Ϊ����
	_, err := l.svcCtx.FriendDB.FindFriendsWithError(l.ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}

	// ������������
	data := make(map[string]interface{})
	if req.IsPinned != nil {
		data["is_pinned"] = req.IsPinned.Value
	}
	if req.Remark != nil {
		data["remark"] = req.Remark.Value
	}
	if req.Ex != nil {
		data["ex"] = req.Ex.Value
	}

	if len(data) > 0 {
		if err := l.svcCtx.FriendDB.UpdateFriends(l.ctx, req.OwnerUserID, req.FriendUserIDs, data); err != nil {
			return nil, err
		}
	}

	// ���ͺ�����Ϣ����֪ͨ
	if l.svcCtx.FriendNotification != nil {
		l.svcCtx.FriendNotification.FriendsInfoUpdateNotification(l.ctx, req.OwnerUserID, req.FriendUserIDs)
	}

	return resp, nil
}
