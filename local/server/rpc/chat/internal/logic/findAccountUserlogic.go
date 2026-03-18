package logic

import (
	"context"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindAccountUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAccountUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAccountUserLogic {
	return &FindAccountUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindAccountUserLogic) FindAccountUser(req *chat.FindAccountUserReq) (*chat.FindAccountUserResp, error) {
	// 1. 验证参数
	if len(req.Accounts) == 0 {
		return nil, errs.ErrArgs.WrapMsg("accounts cannot be empty")
	}

	// 2. 构建查询账户列表
	accounts := make([]*database.UserAccount, 0, len(req.Accounts))
	for _, account := range req.Accounts {
		accounts = append(accounts, &database.UserAccount{
			Account: account,
		})
	}

	// 3. 查找用户账户
	userAccounts, err := l.svcCtx.ChatDB.FindUserAccount(l.ctx, accounts)
	if err != nil {
		l.Errorf("FindUserAccount failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to find user accounts")
	}

	// 4. 构建账户到用户ID的映射
	accountUserMap := make(map[string]string)
	for _, userAccount := range userAccounts {
		if userAccount.Account != "" {
			accountUserMap[userAccount.Account] = userAccount.UserID
		}
	}

	return &chat.FindAccountUserResp{
		AccountUserMap: accountUserMap,
	}, nil
}
