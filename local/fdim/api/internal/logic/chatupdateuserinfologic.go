package logic

import (
	"context"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/protocol/wrapperspb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatUpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatUpdateUserInfoLogic {
	return &ChatUpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChatUpdateUserInfoLogic) ChatUpdateUserInfo(req *types.ChatUpdateUserInfoReq) (*types.ChatUpdateUserInfoResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, errs.ErrInternalServer.WrapMsg("chat service not available")
	}

	rpcReq := &chat.UpdateUserInfoReq{
		UserID: req.UserID,
	}
	if req.Account != nil {
		rpcReq.Account = wrapperspb.StringPtr(req.Account)
	}
	if req.PhoneNumber != nil {
		rpcReq.PhoneNumber = wrapperspb.StringPtr(req.PhoneNumber)
	}
	if req.AreaCode != nil {
		rpcReq.AreaCode = wrapperspb.StringPtr(req.AreaCode)
	}
	if req.Email != nil {
		rpcReq.Email = wrapperspb.StringPtr(req.Email)
	}
	if req.Nickname != nil {
		rpcReq.Nickname = wrapperspb.StringPtr(req.Nickname)
	}
	if req.FaceURL != nil {
		rpcReq.FaceURL = wrapperspb.StringPtr(req.FaceURL)
	}
	if req.Gender != nil {
		rpcReq.Gender = wrapperspb.Int32Ptr(req.Gender)
	}
	if req.Birth != nil {
		rpcReq.Birth = wrapperspb.Int64Ptr(req.Birth)
	}

	_, err := l.svcCtx.ChatClient.UpdateUserInfo(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("ChatUpdateUserInfo failed: %v", err)
		return nil, err
	}

	return &types.ChatUpdateUserInfoResp{}, nil
}
