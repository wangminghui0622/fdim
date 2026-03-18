package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserAccountLogic {
	return &FindUserAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserAccountLogic) FindUserAccount(req *chat.FindUserAccountReq) (*chat.FindUserAccountResp, error) {
	// 1. 验证参数
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("user IDs cannot be empty")
	}

	// 2. 查找用户信息
	userInfos, err := l.svcCtx.ChatDB.FindUserFullInfo(l.ctx, req.UserIDs)
	if err != nil {
		l.Errorf("FindUserFullInfo failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to find user info")
	}

	// 3. 构建用户账户映射
	userAccountMap := make(map[string]string)
	for _, userInfo := range userInfos {
		userAccountMap[userInfo.UserID] = userInfo.Account
	}

	// 4. 如果某些用户ID没有找到账户信息，尝试从账户表查找
	for _, userID := range req.UserIDs {
		if _, exists := userAccountMap[userID]; !exists {
			account, err := l.svcCtx.ChatDB.GetUserAccountByUserID(l.ctx, userID)
			if err == nil {
				userAccountMap[userID] = account.Account
			}
		}
	}

	return &chat.FindUserAccountResp{
		UserAccountMap: userAccountMap,
	}, nil
}
