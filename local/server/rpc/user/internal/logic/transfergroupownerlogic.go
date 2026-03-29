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

type TransferGroupOwnerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransferGroupOwnerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferGroupOwnerLogic {
	return &TransferGroupOwnerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TransferGroupOwnerLogic) TransferGroupOwner(req *user.TransferGroupOwnerReq) (*user.TransferGroupOwnerResp, error) {
	resp := &user.TransferGroupOwnerResp{}

	// ????????????
	_, err := l.svcCtx.GroupDB.TakeGroup(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// ????????????????????????
	opUserID := mcontext.GetOpUserID(l.ctx)
	owner, err := l.svcCtx.GroupDB.TakeGroupOwner(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	if !authverify.IsAdmin(l.ctx) {
		if owner.UserID != opUserID {
			return nil, fmt.Errorf("only group owner can transfer ownership")
		}
	}

	// ???????
	if req.OldOwnerUserID == "" {
		req.OldOwnerUserID = owner.UserID
	}
	if req.OldOwnerUserID != owner.UserID {
		return nil, fmt.Errorf("old owner user id mismatch")
	}
	if req.NewOwnerUserID == req.OldOwnerUserID {
		return nil, fmt.Errorf("new owner cannot be same as old owner")
	}

	// ?????????????????
	_, err = l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, req.NewOwnerUserID)
	if err != nil {
		return nil, fmt.Errorf("new owner not in group")
	}

	// ???????????????
	oldOwnerData := map[string]interface{}{
		"role_level": constant.GroupOrdinaryUsers,
	}
	if err := l.svcCtx.GroupDB.UpdateGroupMemberMap(l.ctx, req.GroupID, req.OldOwnerUserID, oldOwnerData); err != nil {
		return nil, err
	}

	// ?????????????
	newOwnerData := map[string]interface{}{
		"role_level": constant.GroupOwner,
	}
	if err := l.svcCtx.GroupDB.UpdateGroupMemberMap(l.ctx, req.GroupID, req.NewOwnerUserID, newOwnerData); err != nil {
		return nil, err
	}

	// ????????????
	if l.svcCtx.GroupNotification != nil {
		l.svcCtx.GroupNotification.GroupOwnerTransferredNotification(l.ctx, req.GroupID, req.OldOwnerUserID, req.NewOwnerUserID)
	}

	// Webhook AfterTransferGroupOwner ???
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterTransferGroupOwnerReq{
			CallbackCommand: webhook.CallbackAfterTransferGroupOwnerCommand,
			GroupID:         req.GroupID,
			OldOwnerUserID:  req.OldOwnerUserID,
			NewOwnerUserID:  req.NewOwnerUserID,
		}
		cbResp := &webhook.CallbackAfterTransferGroupOwnerResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
