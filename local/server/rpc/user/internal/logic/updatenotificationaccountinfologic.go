package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNotificationAccountInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateNotificationAccountInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotificationAccountInfoLogic {
	return &UpdateNotificationAccountInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateNotificationAccountInfoLogic) UpdateNotificationAccountInfo(req *user.UpdateNotificationAccountInfoReq) (*user.UpdateNotificationAccountInfoResp, error) {
	// 权限验证：需要管理员权限
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	// 检查用户是否存?
	if _, err := l.svcCtx.UserDB.FindWithError(l.ctx, []string{req.UserID}); err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 构建更新数据
	updateData := make(map[string]interface{})
	if req.NickName != "" {
		updateData["nickname"] = req.NickName
	}
	if req.FaceURL != "" {
		updateData["face_url"] = req.FaceURL
	}

	// 更新用户信息
	if err := l.svcCtx.UserDB.UpdateByMap(l.ctx, req.UserID, updateData); err != nil {
		return nil, err
	}

	return &user.UpdateNotificationAccountInfoResp{}, nil
}
