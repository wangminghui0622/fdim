package logic

import (
	"context"
	"time"

	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserLoginCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginCountLogic {
	return &UserLoginCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserLoginCountLogic) UserLoginCount(req *chat.UserLoginCountReq) (*chat.UserLoginCountResp, error) {
	// 1. ֤
	if req.Start == 0 {
		req.Start = time.Now().AddDate(0, 0, -7).Unix() // Ĭ7
	}
	if req.End == 0 {
		req.End = time.Now().Unix()
	}

	// 2. ͳƵ¼
	loginCount, err := l.svcCtx.ChatDB.UserLoginCount(l.ctx, req.Start, req.End)
	if err != nil {
		l.Errorf("UserLoginCount failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to count user login")
	}

	// 3. ʵ֣ص¼δ¼Ҫû
	// 򻯴ֻص¼
	return &chat.UserLoginCountResp{
		LoginCount: loginCount,
		Count:      make(map[string]int64),
	}, nil
}
