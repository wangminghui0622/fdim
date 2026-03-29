package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelUserAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelUserAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserAccountLogic {
	return &DelUserAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelUserAccountLogic) DelUserAccount(req *chat.DelUserAccountReq) (*chat.DelUserAccountResp, error) {
	// 1. ֤
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("user IDs cannot be empty")
	}

	// 2. ɾû˻
	if err := l.svcCtx.ChatDB.DelUserAccount(l.ctx, req.UserIDs); err != nil {
		l.Errorf("DelUserAccount failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to delete user accounts")
	}

	// 3. ɾûϢͨ±Ϊɾֱɾ
	for _, userID := range req.UserIDs {
		if err := l.svcCtx.ChatDB.UpdateUserInfo(l.ctx, userID, map[string]interface{}{
			"deleted": true,
		}); err != nil {
			l.Errorw("UpdateUserInfo failed", logx.Field("userID", userID), logx.Field("error", err))
		}
	}

	l.Infof("Deleted user accounts: count=%d", len(req.UserIDs))
	return &chat.DelUserAccountResp{}, nil
}
