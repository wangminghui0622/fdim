package logic

import (
	"context"

	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckUserExistLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckUserExistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUserExistLogic {
	return &CheckUserExistLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckUserExistLogic) CheckUserExist(req *chat.CheckUserExistReq) (*chat.CheckUserExistResp, error) {
	if req.User == nil {
		return &chat.CheckUserExistResp{IsRegistered: false}, nil
	}

	// 检查手机号
	if req.User.PhoneNumber != "" && req.User.AreaCode != "" {
		userAccount, err := l.svcCtx.ChatDB.GetUserAccountByPhone(l.ctx, req.User.AreaCode, req.User.PhoneNumber)
		if err == nil {
			return &chat.CheckUserExistResp{
				Userid:       userAccount.UserID,
				IsRegistered: true,
			}, nil
		}
	}

	// 检查邮箱
	if req.User.Email != "" {
		userAccount, err := l.svcCtx.ChatDB.GetUserAccountByEmail(l.ctx, req.User.Email)
		if err == nil {
			return &chat.CheckUserExistResp{
				Userid:       userAccount.UserID,
				IsRegistered: true,
			}, nil
		}
	}

	// 检查账户
	if req.User.Account != "" {
		userAccount, err := l.svcCtx.ChatDB.GetUserAccountByAccount(l.ctx, req.User.Account)
		if err == nil {
			return &chat.CheckUserExistResp{
				Userid:       userAccount.UserID,
				IsRegistered: true,
			}, nil
		}
	}

	return &chat.CheckUserExistResp{IsRegistered: false}, nil
}
