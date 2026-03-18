package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/authverify"
	"fdim/pkg/constant"
	"fdim/pkg/mcontext"
	"fdim/pkg/model"
	"fdim/pkg/webhook"
	"fdim/protocol/conversation"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type InviteUserToGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInviteUserToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteUserToGroupLogic {
	return &InviteUserToGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InviteUserToGroupLogic) InviteUserToGroup(req *user.InviteUserToGroupReq) (*user.InviteUserToGroupResp, error) {
	resp := &user.InviteUserToGroupResp{}

	// ������֤
	if len(req.InvitedUserIDs) == 0 {
		return nil, fmt.Errorf("invitedUserIDs is empty")
	}

	// ���Ⱥ���Ƿ����
	groupInfo, err := l.svcCtx.GroupDB.TakeGroup(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// Ȩ����֤����Ҫ��Ⱥ�������Ա
	opUserID := mcontext.GetOpUserID(l.ctx)
	if !authverify.IsAdmin(l.ctx) {
		member, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}
		if member.RoleLevel != constant.GroupOwner && member.RoleLevel != constant.GroupAdmin {
			return nil, fmt.Errorf("only group owner or admin can invite users")
		}
	}

	// TODO: ����û��Ƿ���ڣ���Ҫ userClient.CheckUser��

	// ����Ⱥ����֤Ҫ�������ֱ�Ӽ��뻹�Ǵ�������
	const AllNeedVerification = 1
	const JoinByInvitation = 1

	// �����Ҫ��֤�Ҳ��ǹ���Ա����������
	if groupInfo.NeedVerification == AllNeedVerification && !authverify.IsAdmin(l.ctx) {
		// ���������Ƿ���Ⱥ�������Ա
		member, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil || (member.RoleLevel != constant.GroupOwner && member.RoleLevel != constant.GroupAdmin) {
			// �������������ֱ�Ӽ���
			var requests []*model.GroupRequest
			for _, userID := range req.InvitedUserIDs {
				// ����Ƿ��Ѿ���Ⱥ��Ա
				_, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, userID)
				if err == nil {
					continue // �Ѿ��ǳ�Ա������
				}
				requests = append(requests, &model.GroupRequest{
					UserID:        userID,
					GroupID:       req.GroupID,
					JoinSource:    JoinByInvitation,
					InviterUserID: opUserID,
					ReqTime:       time.Now(),
					HandledTime:   time.Unix(0, 0),
					HandleResult:  constant.GroupRequestUnhandled,
				})
			}
			if len(requests) > 0 {
				if err := l.svcCtx.GroupDB.CreateGroupRequest(l.ctx, requests); err != nil {
					return nil, err
				}
				// ���ͼ���Ⱥ������֪ͨ
				if l.svcCtx.GroupNotification != nil {
					for _, req := range requests {
						l.svcCtx.GroupNotification.JoinGroupApplicationNotification(l.ctx, req.GroupID, req.UserID)
					}
				}
			}
			return resp, nil
		}
	}

	// Webhook BeforeInviteJoinGroup �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackBeforeInviteJoinGroupReq{
			CallbackCommand: webhook.CallbackBeforeInviteJoinGroupCommand,
			GroupID:         req.GroupID,
			InvitedUserIDs:  req.InvitedUserIDs,
			Reason:          req.Reason,
		}
		cbResp := &webhook.CallbackBeforeInviteJoinGroupResp{}
		if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
			if err != webhook.ErrCallbackContinue {
				return nil, err
			}
			// ErrCallbackContinue ��ʾ����ִ��
		}
		// ��� webhook �����˾ܾ��ĳ�Ա�б���������б����Ƴ�
		if len(cbResp.RefusedMembersAccount) > 0 {
			refusedMap := make(map[string]bool)
			for _, refusedID := range cbResp.RefusedMembersAccount {
				refusedMap[refusedID] = true
			}
			filteredUserIDs := make([]string, 0, len(req.InvitedUserIDs))
			for _, userID := range req.InvitedUserIDs {
				if !refusedMap[userID] {
					filteredUserIDs = append(filteredUserIDs, userID)
				}
			}
			req.InvitedUserIDs = filteredUserIDs
		}
	}

	// ֱ�Ӽ���
	now := time.Now()
	var groupMembers []*model.GroupMember
	for _, userID := range req.InvitedUserIDs {
		// ����Ƿ��Ѿ���Ⱥ��Ա
		_, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, userID)
		if err == nil {
			continue // �Ѿ��ǳ�Ա������
		}

		groupMember := &model.GroupMember{
			GroupID:        req.GroupID,
			UserID:         userID,
			RoleLevel:      constant.GroupOrdinaryUsers,
			JoinTime:       now,
			OperatorUserID: opUserID,
			JoinSource:     JoinByInvitation,
			InviterUserID:  opUserID,
			MuteEndTime:    time.UnixMilli(0),
		}
		groupMembers = append(groupMembers, groupMember)
	}

	if len(groupMembers) > 0 {
		// Webhook BeforeMembersJoinGroup �ص�
		if l.svcCtx.WebhookClient != nil {
			memberUserIDs := make([]string, 0, len(groupMembers))
			for _, m := range groupMembers {
				memberUserIDs = append(memberUserIDs, m.UserID)
			}
			cbReq := &webhook.CallbackBeforeMembersJoinGroupReq{
				CallbackCommand: webhook.CallbackBeforeMembersJoinGroupCommand,
				GroupID:         req.GroupID,
				GroupType:       int32(groupInfo.GroupType),
				MemberUserIDs:   memberUserIDs,
				Ex:              groupInfo.Ex,
			}
			cbResp := &webhook.CallbackBeforeMembersJoinGroupResp{}
			if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
				if err != webhook.ErrCallbackContinue {
					return nil, err
				}
				// ErrCallbackContinue ��ʾ����ִ��
			}
			// ��� webhook �������޸ĺ�ĳ�Ա�б����Ҫ���¹��� groupMembers
			if len(cbResp.MemberUserIDs) > 0 {
				// ����򻯴����ʵ��Ӧ�ø��ݷ��صĳ�Ա�б����¹���
			}
		}

		if err := l.svcCtx.GroupDB.CreateGroup(l.ctx, nil, groupMembers); err != nil {
			return nil, err
		}

		newMemberIDs := make([]string, 0, len(groupMembers))
		for _, m := range groupMembers {
			newMemberIDs = append(newMemberIDs, m.UserID)
		}

		// 为新成员创建群聊会话（与官方一致）
		if l.svcCtx.ConversationClient != nil && len(newMemberIDs) > 0 {
			go func() {
				ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, err := l.svcCtx.ConversationClient.CreateGroupChatConversations(ctx2,
					&conversation.CreateGroupChatConversationsReq{
						GroupID: req.GroupID,
						UserIDs: newMemberIDs,
					})
				if err != nil {
					logx.Errorf("CreateGroupChatConversations failed on invite: groupID=%s, error=%v", req.GroupID, err)
				}
			}()
		}

		// ���ͳ�Ա����֪ͨ��������룩
		if l.svcCtx.GroupNotification != nil {
			l.svcCtx.GroupNotification.MemberEnterNotification(l.ctx, req.GroupID, opUserID, newMemberIDs)
		}
		// TODO: ���ó�Ա�������кţ�setMemberJoinSeq��
	}

	return resp, nil
}
