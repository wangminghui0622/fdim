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

type GroupApplicationResponseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupApplicationResponseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupApplicationResponseLogic {
	return &GroupApplicationResponseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupApplicationResponseLogic) GroupApplicationResponse(req *user.GroupApplicationResponseReq) (*user.GroupApplicationResponseResp, error) {
	resp := &user.GroupApplicationResponseResp{}

	// ������֤
	if req.HandleResult != constant.GroupRequestAgree && req.HandleResult != constant.GroupRequestRefuse {
		return nil, fmt.Errorf("invalid handle result: %d", req.HandleResult)
	}

	// Ȩ����֤����Ҫ��Ⱥ�������Ա
	opUserID := mcontext.GetOpUserID(l.ctx)
	if !authverify.IsAdmin(l.ctx) {
		member, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("operator not in group")
		}
		if member.RoleLevel != constant.GroupOwner && member.RoleLevel != constant.GroupAdmin {
			return nil, fmt.Errorf("only group owner or admin can handle application")
		}
	}

	// ���Ⱥ���Ƿ����
	groupInfo, err := l.svcCtx.GroupDB.TakeGroup(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// ���Ⱥ�������Ƿ����
	groupRequest, err := l.svcCtx.GroupDB.TakeGroupRequest(l.ctx, req.GroupID, req.FromUserID)
	if err != nil {
		return nil, err
	}

	// ��������Ƿ��Ѵ���
	if groupRequest.HandleResult != constant.GroupRequestUnhandled {
		return nil, fmt.Errorf("group request already processed")
	}

	// ����û��Ƿ��Ѿ���Ⱥ����
	var inGroup bool
	_, err = l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, req.FromUserID)
	if err == nil {
		inGroup = true
	}

	// TODO: ����û��Ƿ���ڣ���Ҫ userClient.CheckUser��

	var member *model.GroupMember
	// ���ͬ�����û�����Ⱥ���У�����Ⱥ��Ա
	if req.HandleResult == constant.GroupRequestAgree && !inGroup {
		member = &model.GroupMember{
			GroupID:        req.GroupID,
			UserID:         req.FromUserID,
			RoleLevel:      constant.GroupOrdinaryUsers,
			JoinTime:       time.Now(),
			OperatorUserID: opUserID,
			JoinSource:     groupRequest.JoinSource,
			InviterUserID:  groupRequest.InviterUserID,
			MuteEndTime:    time.UnixMilli(0),
		}

		// Webhook BeforeMembersJoinGroup �ص�
		if l.svcCtx.WebhookClient != nil {
			cbReq := &webhook.CallbackBeforeMembersJoinGroupReq{
				CallbackCommand: webhook.CallbackBeforeMembersJoinGroupCommand,
				GroupID:         req.GroupID,
				GroupType:       int32(groupInfo.GroupType),
				MemberUserIDs:   []string{req.FromUserID},
				Ex:              groupInfo.Ex,
			}
			cbResp := &webhook.CallbackBeforeMembersJoinGroupResp{}
			if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
				if err != webhook.ErrCallbackContinue {
					return nil, err
				}
				// ErrCallbackContinue ��ʾ����ִ��
			}
			// ��� webhook �������޸ĺ�ĳ�Ա�б����Ҫ���¹��� member
			if len(cbResp.MemberUserIDs) > 0 {
				// ����򻯴����ʵ��Ӧ�ø��ݷ��صĳ�Ա�б����¹���
			}
		}
	}

	// ����Ⱥ�����루����״̬�����ͬ���򴴽���Ա��
	if err := l.svcCtx.GroupDB.HandlerGroupRequest(l.ctx, req.GroupID, req.FromUserID, req.HandledMsg, req.HandleResult, member); err != nil {
		return nil, err
	}

	// ����֪ͨ
	if req.HandleResult == constant.GroupRequestAgree {
		// ����Ⱥ���������֪ͨ
		if l.svcCtx.GroupNotification != nil {
			l.svcCtx.GroupNotification.GroupApplicationAcceptedNotification(l.ctx, req.GroupID, opUserID, req.FromUserID)
		}
		if member != nil {
			// 为新成员创建群聊会话（与官方一致）
			if l.svcCtx.ConversationClient != nil {
				go func() {
					ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					_, err := l.svcCtx.ConversationClient.CreateGroupChatConversations(ctx2,
						&conversation.CreateGroupChatConversationsReq{
							GroupID: req.GroupID,
							UserIDs: []string{req.FromUserID},
						})
					if err != nil {
						logx.Errorf("CreateGroupChatConversations failed on approve: groupID=%s, error=%v", req.GroupID, err)
					}
				}()
			}

			if groupRequest.InviterUserID == "" {
				// ���ͳ�Ա����֪ͨ��������룩
				if l.svcCtx.GroupNotification != nil {
					l.svcCtx.GroupNotification.MemberEnterNotification(l.ctx, req.GroupID, opUserID, []string{req.FromUserID})
				}
			} else {
				// ����Ⱥ������ͬ����Ա����֪ͨ��������룩
				if l.svcCtx.GroupNotification != nil {
					l.svcCtx.GroupNotification.GroupApplicationAgreeMemberEnterNotification(l.ctx, req.GroupID, opUserID, []string{req.FromUserID})
				}
			}
			// TODO: ���ó�Ա�������кţ�setMemberJoinSeq��
		}
	} else {
		// ����Ⱥ������ܾ�֪ͨ
		if l.svcCtx.GroupNotification != nil {
			l.svcCtx.GroupNotification.GroupApplicationRejectedNotification(l.ctx, req.GroupID, opUserID, req.FromUserID, req.HandledMsg)
		}
	}

	return resp, nil
}
