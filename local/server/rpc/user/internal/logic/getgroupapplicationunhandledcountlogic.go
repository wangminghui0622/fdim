package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupApplicationUnhandledCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupApplicationUnhandledCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupApplicationUnhandledCountLogic {
	return &GetGroupApplicationUnhandledCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupApplicationUnhandledCountLogic) GetGroupApplicationUnhandledCount(req *user.GetGroupApplicationUnhandledCountReq) (*user.GetGroupApplicationUnhandledCountResp, error) {
	resp := &user.GetGroupApplicationUnhandledCountResp{}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// ��ȡ�û����������Ⱥ���Ա��¼
	members, err := l.svcCtx.GroupDB.FindGroupMemberByUserID(l.ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// ���û�м����κ�Ⱥ�飬����0
	if len(members) == 0 {
		resp.Count = 0
		return resp, nil
	}

	// ��ȡȺ��ID�б�
	groupIDs := make([]string, len(members))
	for i, m := range members {
		groupIDs[i] = m.GroupID
	}

	// ��ȡδ�����Ⱥ����������
	count, err := l.svcCtx.GroupDB.GetUnhandledGroupRequestCount(l.ctx, groupIDs, req.Time)
	if err != nil {
		return nil, err
	}

	resp.Count = int64(count)
	return resp, nil
}
