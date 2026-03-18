package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/authverify"
	"fdim/pkg/constant"
	"fdim/pkg/mcontext"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupInfoLogic {
	return &SetGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetGroupInfoLogic) SetGroupInfo(req *user.SetGroupInfoReq) (*user.SetGroupInfoResp, error) {
	resp := &user.SetGroupInfoResp{}

	// ������֤
	if req.GroupInfoForSet == nil {
		return nil, fmt.Errorf("groupInfoForSet is empty")
	}
	if req.GroupInfoForSet.GroupID == "" {
		return nil, fmt.Errorf("groupID is empty")
	}

	// Ȩ����֤����Ҫ��Ⱥ�������Ա
	opUserID := mcontext.GetOpUserID(l.ctx)
	if !authverify.IsAdmin(l.ctx) {
		member, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupInfoForSet.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}
		if member.RoleLevel != constant.GroupOwner && member.RoleLevel != constant.GroupAdmin {
			return nil, fmt.Errorf("only group owner or admin can set group info")
		}
	}

	// ���Ⱥ���Ƿ����
	_, err := l.svcCtx.GroupDB.TakeGroup(l.ctx, req.GroupInfoForSet.GroupID)
	if err != nil {
		return nil, err
	}

	// ׼�� Webhook �ص�����
	var groupName, introduction, faceURL, ex *string
	if req.GroupInfoForSet.GroupName != "" {
		groupName = &req.GroupInfoForSet.GroupName
	}
	if req.GroupInfoForSet.Introduction != "" {
		introduction = &req.GroupInfoForSet.Introduction
	}
	if req.GroupInfoForSet.FaceURL != "" {
		faceURL = &req.GroupInfoForSet.FaceURL
	}
	if req.GroupInfoForSet.Ex != nil {
		ex = &req.GroupInfoForSet.Ex.Value
	}

	// Webhook BeforeSetGroupInfo �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackBeforeSetGroupInfoReq{
			CallbackCommand: webhook.CallbackBeforeSetGroupInfoCommand,
			GroupID:         req.GroupInfoForSet.GroupID,
			GroupName:       groupName,
			Introduction:    introduction,
			FaceURL:         faceURL,
			Ex:              ex,
		}
		cbResp := &webhook.CallbackBeforeSetGroupInfoResp{}
		if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
			if err != webhook.ErrCallbackContinue {
				return nil, err
			}
			// ErrCallbackContinue ��ʾ����ִ��
		}
		// ��� webhook �������޸ĺ��ֵ��ʹ����
		if cbResp.GroupName != nil {
			req.GroupInfoForSet.GroupName = *cbResp.GroupName
		}
		if cbResp.Introduction != nil {
			req.GroupInfoForSet.Introduction = *cbResp.Introduction
		}
		if cbResp.FaceURL != nil {
			req.GroupInfoForSet.FaceURL = *cbResp.FaceURL
		}
		if cbResp.Ex != nil && req.GroupInfoForSet.Ex != nil {
			req.GroupInfoForSet.Ex.Value = *cbResp.Ex
		}
	}

	// ������������
	data := make(map[string]interface{})
	if req.GroupInfoForSet.GroupName != "" {
		data["group_name"] = req.GroupInfoForSet.GroupName
	}
	if req.GroupInfoForSet.Notification != "" {
		data["notification"] = req.GroupInfoForSet.Notification
		data["notification_update_time"] = time.Now()
		data["notification_user_id"] = opUserID
	}
	if req.GroupInfoForSet.Introduction != "" {
		data["introduction"] = req.GroupInfoForSet.Introduction
	}
	if req.GroupInfoForSet.FaceURL != "" {
		data["face_url"] = req.GroupInfoForSet.FaceURL
	}
	if req.GroupInfoForSet.Ex != nil {
		data["ex"] = req.GroupInfoForSet.Ex.Value
	}
	if req.GroupInfoForSet.NeedVerification != nil {
		data["need_verification"] = req.GroupInfoForSet.NeedVerification.Value
	}
	if req.GroupInfoForSet.LookMemberInfo != nil {
		data["look_member_info"] = req.GroupInfoForSet.LookMemberInfo.Value
	}
	if req.GroupInfoForSet.ApplyMemberFriend != nil {
		data["apply_member_friend"] = req.GroupInfoForSet.ApplyMemberFriend.Value
	}

	if len(data) > 0 {
		if err := l.svcCtx.GroupDB.UpdateGroupMap(l.ctx, req.GroupInfoForSet.GroupID, data); err != nil {
			return nil, err
		}
	}

	// ����Ⱥ����Ϣ����֪ͨ
	if l.svcCtx.GroupNotification != nil {
		l.svcCtx.GroupNotification.GroupInfoSetNotification(l.ctx, req.GroupInfoForSet.GroupID, opUserID)
	}

	// Webhook AfterSetGroupInfo �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterSetGroupInfoReq{
			CallbackCommand: webhook.CallbackAfterSetGroupInfoCommand,
			GroupID:         req.GroupInfoForSet.GroupID,
			GroupName:       req.GroupInfoForSet.GroupName,
			Introduction:    req.GroupInfoForSet.Introduction,
			FaceURL:         req.GroupInfoForSet.FaceURL,
			Ex:              "",
		}
		if req.GroupInfoForSet.Ex != nil {
			cbReq.Ex = req.GroupInfoForSet.Ex.Value
		}
		cbResp := &webhook.CallbackAfterSetGroupInfoResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
