package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type BlockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockUserLogic {
	return &BlockUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BlockUserLogic) BlockUser(req *admin.BlockUserReq) (*admin.BlockUserResp, error) {
	// 1. 验证参数
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID cannot be empty")
	}

	// 2. 获取操作者ID（从context中获取）
	opUserID := l.getUserIDFromContext()
	if opUserID == "" {
		opUserID = "system" // 默认系统操作
	}

	// 3. 检查用户是否已被封禁
	_, err := l.svcCtx.AdminDB.GetBlockInfo(l.ctx, req.UserID)
	if err == nil {
		return nil, fmt.Errorf("user already blocked")
	}

	// 4. 添加封禁记录
	block := &database.BlockUser{
		UserID:         req.UserID,
		Reason:         req.Reason,
		OperatorUserID: opUserID,
		CreateTime:     time.Now(),
	}

	if err := l.svcCtx.AdminDB.AddBlockUser(l.ctx, []*database.BlockUser{block}); err != nil {
		l.Errorf("AddBlockUser failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to block user")
	}

	l.Infof("Blocked user: userID=%s, reason=%s, opUserID=%s", req.UserID, req.Reason, opUserID)
	return &admin.BlockUserResp{}, nil
}

func (l *BlockUserLogic) getUserIDFromContext() string {
	if md, ok := metadata.FromIncomingContext(l.ctx); ok {
		if userIDs := md.Get("userID"); len(userIDs) > 0 {
			return userIDs[0]
		}
		if userIDs := md.Get("UserID"); len(userIDs) > 0 {
			return userIDs[0]
		}
	}
	if userID, ok := l.ctx.Value("userID").(string); ok {
		return userID
	}
	if userID, ok := l.ctx.Value("UserID").(string); ok {
		return userID
	}
	return ""
}
