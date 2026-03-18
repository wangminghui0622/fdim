package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoExLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoExLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoExLogic {
	return &UpdateUserInfoExLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoExLogic) UpdateUserInfoEx(req *user.UpdateUserInfoExReq) (*user.UpdateUserInfoExResp, error) {
	resp := &user.UpdateUserInfoExResp{}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.UserInfo.UserID); err != nil {
		return nil, err
	}

	// ������֤
	if req.UserInfo.UserID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	// ����û��Ƿ����
	user, err := l.svcCtx.UserDB.Take(l.ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: %s", req.UserInfo.UserID)
	}

	// Webhook Before �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackBeforeUpdateUserInfoExReq{
			CallbackCommand: webhook.CallbackBeforeUpdateUserInfoExCommand,
			UserID:          req.UserInfo.UserID,
		}
		if req.UserInfo.FaceURL != nil {
			cbReq.FaceURL = &req.UserInfo.FaceURL.Value
		}
		if req.UserInfo.Nickname != nil {
			cbReq.Nickname = &req.UserInfo.Nickname.Value
		}
		if req.UserInfo.Ex != nil {
			cbReq.Ex = &req.UserInfo.Ex.Value
		}
		cbResp := &webhook.CallbackBeforeUpdateUserInfoExResp{}
		if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
			if err != webhook.ErrCallbackContinue {
				return nil, err
			}
			// ErrCallbackContinue ��ʾ����ִ��
		}
		// �����������ݣ���� webhook �������޸ģ�
		if cbResp.FaceURL != nil && req.UserInfo.FaceURL != nil {
			req.UserInfo.FaceURL.Value = *cbResp.FaceURL
		}
		if cbResp.Nickname != nil && req.UserInfo.Nickname != nil {
			req.UserInfo.Nickname.Value = *cbResp.Nickname
		}
		if cbResp.Ex != nil && req.UserInfo.Ex != nil {
			req.UserInfo.Ex.Value = *cbResp.Ex
		}
	}

	// ת���������û���Ϣ
	data := convert.UserPb2DBMapEx(req.UserInfo)
	if len(data) == 0 {
		return resp, nil // û����Ҫ���µ��ֶ�
	}

	if err := l.svcCtx.UserDB.UpdateByMap(l.ctx, req.UserInfo.UserID, data); err != nil {
		return nil, err
	}

	// �����û���Ϣ����֪ͨ��֪ͨ�û��Լ���
	if l.svcCtx.UserNotification != nil {
		l.svcCtx.UserNotification.UserInfoUpdatedNotification(l.ctx, req.UserInfo.UserID, req.UserInfo.UserID)
	}

	// Webhook After �ص�
	if l.svcCtx.WebhookClient != nil {
		faceURL := ""
		nickname := ""
		ex := ""
		if req.UserInfo.FaceURL != nil {
			faceURL = req.UserInfo.FaceURL.Value
		}
		if req.UserInfo.Nickname != nil {
			nickname = req.UserInfo.Nickname.Value
		}
		if req.UserInfo.Ex != nil {
			ex = req.UserInfo.Ex.Value
		}
		cbReq := &webhook.CallbackAfterUpdateUserInfoExReq{
			CallbackCommand: webhook.CallbackAfterUpdateUserInfoExCommand,
			UserID:          req.UserInfo.UserID,
			FaceURL:         faceURL,
			Nickname:        nickname,
			Ex:              ex,
		}
		cbResp := &webhook.CallbackAfterUpdateUserInfoExResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
