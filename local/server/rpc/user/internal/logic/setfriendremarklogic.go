package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetFriendRemarkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetFriendRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetFriendRemarkLogic {
	return &SetFriendRemarkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetFriendRemarkLogic) SetFriendRemark(req *user.SetFriendRemarkReq) (*user.SetFriendRemarkResp, error) {
	resp := &user.SetFriendRemarkResp{}

	// ??????
	if err := authverify.CheckAccess(l.ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	// ???????????
	_, err := l.svcCtx.FriendDB.FindFriendsWithError(l.ctx, req.OwnerUserID, []string{req.FriendUserID})
	if err != nil {
		return nil, err
	}

	// Webhook BeforeSetFriendRemark ???
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackBeforeSetFriendRemarkReq{
			CallbackCommand: webhook.CallbackBeforeSetFriendRemarkCommand,
			OwnerUserID:     req.OwnerUserID,
			FriendUserID:    req.FriendUserID,
			Remark:          req.Remark,
		}
		cbResp := &webhook.CallbackBeforeSetFriendRemarkResp{}
		if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
			if err != webhook.ErrCallbackContinue {
				return nil, err
			}
			// ErrCallbackContinue ??????????
		}
		// ??? webhook ????????????????????
		if cbResp.Remark != nil {
			req.Remark = *cbResp.Remark
		}
	}

	// ?????
	if err := l.svcCtx.FriendDB.UpdateRemark(l.ctx, req.OwnerUserID, req.FriendUserID, req.Remark); err != nil {
		return nil, err
	}

	// ???????????????
	if l.svcCtx.FriendNotification != nil {
		l.svcCtx.FriendNotification.FriendRemarkSetNotification(l.ctx, req.OwnerUserID, req.FriendUserID)
	}

	// Webhook AfterSetFriendRemark ???
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterSetFriendRemarkReq{
			CallbackCommand: webhook.CallbackAfterSetFriendRemarkCommand,
			OwnerUserID:     req.OwnerUserID,
			FriendUserID:    req.FriendUserID,
			Remark:          req.Remark,
		}
		cbResp := &webhook.CallbackAfterSetFriendRemarkResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
