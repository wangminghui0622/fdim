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

type DismissGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDismissGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DismissGroupLogic {
	return &DismissGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DismissGroupLogic) DismissGroup(req *user.DismissGroupReq) (*user.DismissGroupResp, error) {
	resp := &user.DismissGroupResp{}

	// ���Ⱥ���Ƿ����
	groupInfo, err := l.svcCtx.GroupDB.TakeGroup(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// ����Ƿ���Ⱥ����ֻ��Ⱥ�������Ա���Խ�ɢȺ�飩
	owner, err := l.svcCtx.GroupDB.TakeGroupOwner(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	opUserID := mcontext.GetOpUserID(l.ctx)
	if !authverify.IsAdmin(l.ctx) {
		if owner.UserID != opUserID {
			return nil, fmt.Errorf("only group owner can dismiss group")
		}
	}

	// ���Ⱥ��״̬
	if !req.DeleteMember && groupInfo.Status == constant.GroupStatusDismissed {
		return nil, fmt.Errorf("group is already dismissed")
	}

	// ����Ⱥ��״̬Ϊ�ѽ�ɢ
	data := map[string]interface{}{
		"status": constant.GroupStatusDismissed,
	}
	if err := l.svcCtx.GroupDB.UpdateGroupMap(l.ctx, req.GroupID, data); err != nil {
		return nil, err
	}

	// ���DeleteMemberΪtrue��ɾ������Ⱥ��Ա
	if req.DeleteMember {
		members, err := l.svcCtx.GroupDB.FindGroupMemberAll(l.ctx, req.GroupID)
		if err == nil && len(members) > 0 {
			userIDs := make([]string, len(members))
			for i, m := range members {
				userIDs[i] = m.UserID
			}
			if err := l.svcCtx.GroupDB.DeleteGroupMember(l.ctx, req.GroupID, userIDs); err != nil {
				return nil, err
			}
		}
	}

	// ����Ⱥ���ɢ֪ͨ
	if l.svcCtx.GroupNotification != nil {
		l.svcCtx.GroupNotification.GroupDismissedNotification(l.ctx, req.GroupID, opUserID)
	}

	// Webhook AfterDismissGroup �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterDismissGroupReq{
			CallbackCommand: webhook.CallbackAfterDismissGroupCommand,
			GroupID:         req.GroupID,
			OpUserID:        opUserID,
		}
		cbResp := &webhook.CallbackAfterDismissGroupResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
