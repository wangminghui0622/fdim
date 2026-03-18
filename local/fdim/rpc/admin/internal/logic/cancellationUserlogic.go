package logic

import (
	"context"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CancellationUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancellationUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancellationUserLogic {
	return &CancellationUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancellationUserLogic) CancellationUser(req *admin.CancellationUserReq) (*admin.CancellationUserResp, error) {
	// 1. 参数校验
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID cannot be empty")
	}

	// 2. 标记用户为封禁，相当于注销账号（简化实现）
	// 如果已经封禁则直接返回
	if _, err := l.svcCtx.AdminDB.GetBlockInfo(l.ctx, req.UserID); err == nil {
		// 已经是封禁状态，视为注销成功
		return &admin.CancellationUserResp{}, nil
	}

	// 3. 添加封禁记录
	block := &database.BlockUser{
		UserID:         req.UserID,
		Reason:         req.Reason,
		OperatorUserID: "system", // 可以按需从 ctx 中提取真实操作者
		CreateTime:     time.Now(),
	}
	if err := l.svcCtx.AdminDB.AddBlockUser(l.ctx, []*database.BlockUser{block}); err != nil {
		l.Errorf("AddBlockUser failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to cancellation user (block user)")
	}

	// 4. 使该用户所有 token 失效，强制下线
	if err := l.svcCtx.AdminDB.InvalidateToken(l.ctx, req.UserID); err != nil {
		l.Errorw("InvalidateToken failed", logx.Field("userID", req.UserID), logx.Field("error", err))
	}

	l.Infof("CancellationUser success: userID=%s, reason=%s", req.UserID, req.Reason)
	return &admin.CancellationUserResp{}, nil
}
