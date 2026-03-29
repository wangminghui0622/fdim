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
	// 1. У
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID cannot be empty")
	}

	// 2. ûΪ൱ע˺ţʵ֣
	// Ѿֱӷ
	if _, err := l.svcCtx.AdminDB.GetBlockInfo(l.ctx, req.UserID); err == nil {
		// ѾǷ״̬Ϊעɹ
		return &admin.CancellationUserResp{}, nil
	}

	// 3. ӷ¼
	block := &database.BlockUser{
		UserID:         req.UserID,
		Reason:         req.Reason,
		OperatorUserID: "system", // ԰ ctx ȡʵ
		CreateTime:     time.Now(),
	}
	if err := l.svcCtx.AdminDB.AddBlockUser(l.ctx, []*database.BlockUser{block}); err != nil {
		l.Errorf("AddBlockUser failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to cancellation user (block user)")
	}

	// 4. ʹû token ʧЧǿ
	if err := l.svcCtx.AdminDB.InvalidateToken(l.ctx, req.UserID); err != nil {
		l.Errorw("InvalidateToken failed", logx.Field("userID", req.UserID), logx.Field("error", err))
	}

	l.Infof("CancellationUser success: userID=%s, reason=%s", req.UserID, req.Reason)
	return &admin.CancellationUserResp{}, nil
}
