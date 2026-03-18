package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *user.UpdateUserInfoReq) (*user.UpdateUserInfoResp, error) {
	resp := &user.UpdateUserInfoResp{}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.UserInfo.UserID); err != nil {
		return nil, err
	}

	// Webhook Before �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackBeforeUpdateUserInfoReq{
			CallbackCommand: webhook.CallbackBeforeUpdateUserInfoCommand,
			UserID:          req.UserInfo.UserID,
		}
		if req.UserInfo.FaceURL != "" {
			cbReq.FaceURL = &req.UserInfo.FaceURL
		}
		if req.UserInfo.Nickname != "" {
			cbReq.Nickname = &req.UserInfo.Nickname
		}
		if req.UserInfo.Ex != "" {
			cbReq.Ex = &req.UserInfo.Ex
		}
		cbResp := &webhook.CallbackBeforeUpdateUserInfoResp{}
		if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
			if err != webhook.ErrCallbackContinue {
				return nil, err
			}
			// ErrCallbackContinue ��ʾ����ִ��
		}
		// �����������ݣ���� webhook �������޸ģ�
		if cbResp.FaceURL != nil {
			req.UserInfo.FaceURL = *cbResp.FaceURL
		}
		if cbResp.Nickname != nil {
			req.UserInfo.Nickname = *cbResp.Nickname
		}
		if cbResp.Ex != nil {
			req.UserInfo.Ex = *cbResp.Ex
		}
	}

	data := convert.UserPb2DBMap(req.UserInfo)
	if err := l.svcCtx.UserDB.UpdateByMap(l.ctx, req.UserInfo.UserID, data); err != nil {
		return nil, err
	}

	// �����û���Ϣ����֪ͨ��֪ͨ�û��Լ���
	if l.svcCtx.UserNotification != nil {
		l.svcCtx.UserNotification.UserInfoUpdatedNotification(l.ctx, req.UserInfo.UserID, req.UserInfo.UserID)
	}

	// Webhook After �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterUpdateUserInfoReq{
			CallbackCommand: webhook.CallbackAfterUpdateUserInfoCommand,
			UserID:          req.UserInfo.UserID,
			FaceURL:         req.UserInfo.FaceURL,
			Nickname:        req.UserInfo.Nickname,
			Ex:              req.UserInfo.Ex,
		}
		cbResp := &webhook.CallbackAfterUpdateUserInfoResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
