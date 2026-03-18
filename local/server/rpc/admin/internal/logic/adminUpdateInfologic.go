package logic

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/pkg/mcontext"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminUpdateInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateInfoLogic {
	return &AdminUpdateInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminUpdateInfoLogic) AdminUpdateInfo(req *admin.AdminUpdateInfoReq) (*admin.AdminUpdateInfoResp, error) {
	// 1. 从 context 中获取 userID
	userID := mcontext.GetOpUserID(l.ctx)
	if userID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID not found in context")
	}

	// 2. 获取当前管理员信息
	adminInfo, err := l.svcCtx.AdminDB.GetAdminUserID(l.ctx, userID)
	if err != nil {
		l.Errorf("GetAdminUserID failed: %v", err)
		return nil, errs.ErrArgs.WrapMsg("admin not found")
	}

	// 3. 构建更新字段
	update := make(map[string]interface{})
	if req.Account != nil && req.Account.Value != "" {
		update["account"] = req.Account.Value
	}
	if req.Password != nil && req.Password.Value != "" {
		update["password"] = req.Password.Value
	}
	if req.FaceURL != nil {
		update["face_url"] = req.FaceURL.Value
	}
	if req.Nickname != nil && req.Nickname.Value != "" {
		update["nickname"] = req.Nickname.Value
	}
	if req.Level != nil {
		update["level"] = req.Level.Value
	}

	if len(update) == 0 {
		return nil, errs.ErrArgs.WrapMsg("no update info")
	}

	// 4. 更新管理员信息
	if err := l.svcCtx.AdminDB.UpdateAdmin(l.ctx, userID, update); err != nil {
		l.Errorf("UpdateAdmin failed: %v", err)
		return nil, fmt.Errorf("failed to update admin: %w", err)
	}

	// 5. 构建响应
	resp := &admin.AdminUpdateInfoResp{
		UserID: adminInfo.UserID,
	}

	// 设置响应字段（使用更新后的值或原值）
	if req.Nickname != nil {
		resp.Nickname = req.Nickname.Value
	} else {
		resp.Nickname = adminInfo.Nickname
	}

	if req.FaceURL != nil {
		resp.FaceURL = req.FaceURL.Value
	} else {
		resp.FaceURL = adminInfo.FaceURL
	}

	l.Infof("Admin %s updated info successfully", userID)

	return resp, nil
}
