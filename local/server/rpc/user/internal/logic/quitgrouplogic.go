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

type QuitGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuitGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuitGroupLogic {
	return &QuitGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QuitGroupLogic) QuitGroup(req *user.QuitGroupReq) (*user.QuitGroupResp, error) {
	resp := &user.QuitGroupResp{}

	// ????????????????UserID????????????????????????ID
	opUserID := mcontext.GetOpUserID(l.ctx)
	if req.UserID == "" {
		if opUserID == "" {
			return nil, fmt.Errorf("userID is empty")
		}
		req.UserID = opUserID
	} else {
		if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
			return nil, err
		}
	}

	// ????????????
	member, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not in group")
	}

	// ????????????????????????????????
	if member.RoleLevel == constant.GroupOwner {
		return nil, fmt.Errorf("group owner cannot quit, please dismiss group")
	}

	// ???????
	if err := l.svcCtx.GroupDB.DeleteGroupMember(l.ctx, req.GroupID, []string{req.UserID}); err != nil {
		return nil, err
	}

	// ???????????
	if l.svcCtx.GroupNotification != nil {
		l.svcCtx.GroupNotification.MemberQuitNotification(l.ctx, req.GroupID, req.UserID)
	}

	// Webhook AfterQuitGroup ???
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterQuitGroupReq{
			CallbackCommand: webhook.CallbackAfterQuitGroupCommand,
			GroupID:         req.GroupID,
			UserID:          req.UserID,
		}
		cbResp := &webhook.CallbackAfterQuitGroupResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}
	// TODO: ??????????deleteMemberAndSetConversationSeq??

	return resp, nil
}
