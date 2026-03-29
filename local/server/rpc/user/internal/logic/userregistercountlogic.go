package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRegisterCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterCountLogic {
	return &UserRegisterCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRegisterCountLogic) UserRegisterCount(req *user.UserRegisterCountReq) (*user.UserRegisterCountResp, error) {
	resp := &user.UserRegisterCountResp{}

	// Ȩ֤ҪԱȨ
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	// ͳû
	total, err := l.svcCtx.UserDB.CountUsers(l.ctx)
	if err != nil {
		return nil, err
	}

	// ͳָʱ䷶Χ֮ǰû
	var before int64
	if req.Start > 0 {
		before, err = l.svcCtx.UserDB.CountUsersBeforeTime(l.ctx, req.Start)
		if err != nil {
			return nil, err
		}
	}

	resp.Total = int64(total)
	resp.Before = int64(before)

	return resp, nil
}
