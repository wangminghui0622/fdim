package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/pkg/mcontext"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAdminInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAdminInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdminInfoLogic {
	return &GetAdminInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAdminInfoLogic) GetAdminInfo(req *admin.GetAdminInfoReq) (*admin.GetAdminInfoResp, error) {
	// 1. 从 context 中获取 userID（通常由中间件从 token 中解析并放入 context）
	userID := mcontext.GetOpUserID(l.ctx)
	if userID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID not found in context")
	}

	// 2. 从数据库获取管理员信息
	adminInfo, err := l.svcCtx.AdminDB.GetAdminUserID(l.ctx, userID)
	if err != nil {
		l.Errorf("GetAdminUserID failed: %v", err)
		return nil, errs.ErrArgs.WrapMsg("admin not found")
	}

	// 3. 返回管理员信息
	return &admin.GetAdminInfoResp{
		UserID:     adminInfo.UserID,
		Account:    adminInfo.Account,
		Nickname:   adminInfo.Nickname,
		FaceURL:    adminInfo.FaceURL,
		Level:      adminInfo.Level,
		CreateTime: adminInfo.CreateTime,
	}, nil
}
