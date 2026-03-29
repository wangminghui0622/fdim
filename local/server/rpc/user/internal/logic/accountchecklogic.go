package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AccountCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccountCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountCheckLogic {
	return &AccountCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AccountCheckLogic) AccountCheck(req *user.AccountCheckReq) (*user.AccountCheckResp, error) {
	resp := &user.AccountCheckResp{}

	// Ȩ֤ҪԱȨ
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	// ֤
	if len(req.CheckUserIDs) == 0 {
		return nil, fmt.Errorf("checkUserIDs is empty")
	}

	// û
	users, err := l.svcCtx.UserDB.Find(l.ctx, req.CheckUserIDs)
	if err != nil {
		return nil, err
	}

	// ûIDӳ
	userIDMap := make(map[string]bool)
	for _, u := range users {
		userIDMap[u.UserID] = true
	}

	// Ӧ
	for _, userID := range req.CheckUserIDs {
		status := &user.AccountCheckRespSingleUserStatus{
			UserID: userID,
		}
		if userIDMap[userID] {
			status.AccountStatus = 1 // Registered
		} else {
			status.AccountStatus = 0 // UnRegistered
		}
		resp.Results = append(resp.Results, status)
	}

	return resp, nil
}
