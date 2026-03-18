package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/constant"
	"fdim/pkg/mcontext"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type KickGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKickGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickGroupMemberLogic {
	return &KickGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *KickGroupMemberLogic) KickGroupMember(req *user.KickGroupMemberReq) (*user.KickGroupMemberResp, error) {
	resp := &user.KickGroupMemberResp{}

	// ������֤
	if len(req.KickedUserIDs) == 0 {
		return nil, fmt.Errorf("kickedUserIDs is empty")
	}

	// ���Ⱥ���Ƿ����
	_, err := l.svcCtx.GroupDB.TakeGroup(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	opUserID := mcontext.GetOpUserID(l.ctx)
	// ���������Ƿ���Ҫ�߳����б���
	for _, userID := range req.KickedUserIDs {
		if userID == opUserID {
			return nil, fmt.Errorf("cannot kick yourself")
		}
	}

	// ���Ⱥ��
	owner, err := l.svcCtx.GroupDB.TakeGroupOwner(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	// ������Ⱥ��
	for _, userID := range req.KickedUserIDs {
		if userID == owner.UserID {
			return nil, fmt.Errorf("cannot kick group owner")
		}
	}

	// Ȩ�޼��
	isAdmin := authverify.IsAdmin(l.ctx)
	if !isAdmin {
		// ���������Ƿ���Ⱥ��Ա
		opMember, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("operator not in group")
		}

		// ���Ҫ�߳��ĳ�Ա
		members, err := l.svcCtx.GroupDB.FindGroupMembers(l.ctx, req.GroupID, req.KickedUserIDs)
		if err != nil {
			return nil, err
		}
		if len(members) != len(req.KickedUserIDs) {
			return nil, fmt.Errorf("some users not found in group")
		}

		// ���ݲ����˽�ɫ���Ȩ��
		switch opMember.RoleLevel {
		case constant.GroupOwner:
			// Ⱥ���������κ��ˣ������Լ������������飩
		case constant.GroupAdmin:
			// ����Աֻ������ͨ��Ա
			for _, member := range members {
				if member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin {
					return nil, fmt.Errorf("admin cannot kick group owner or other admins")
				}
			}
		case constant.GroupOrdinaryUsers:
			// ��ͨ��Ա��������
			return nil, fmt.Errorf("ordinary member cannot kick others")
		default:
			return nil, fmt.Errorf("unknown role level")
		}
	} else {
		// ����Ա�������κ��ˣ�����Ҫ����Ա�Ƿ����
		members, err := l.svcCtx.GroupDB.FindGroupMembers(l.ctx, req.GroupID, req.KickedUserIDs)
		if err != nil {
			return nil, err
		}
		if len(members) != len(req.KickedUserIDs) {
			return nil, fmt.Errorf("some users not found in group")
		}
	}

	// ɾ��Ⱥ��Ա
	if err := l.svcCtx.GroupDB.DeleteGroupMember(l.ctx, req.GroupID, req.KickedUserIDs); err != nil {
		return nil, err
	}

	// ���ͳ�Ա���߳�֪ͨ
	if l.svcCtx.GroupNotification != nil {
		l.svcCtx.GroupNotification.MemberKickedNotification(l.ctx, req.GroupID, opUserID, req.KickedUserIDs)
	}

	// Webhook AfterKickGroupMember �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterKickGroupMemberReq{
			CallbackCommand: webhook.CallbackAfterKickGroupMemberCommand,
			GroupID:         req.GroupID,
			KickedUserIDs:   req.KickedUserIDs,
			OpUserID:        opUserID,
		}
		cbResp := &webhook.CallbackAfterKickGroupMemberResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}
	// TODO: ���ûỰ���кţ�deleteMemberAndSetConversationSeq��

	return resp, nil
}
