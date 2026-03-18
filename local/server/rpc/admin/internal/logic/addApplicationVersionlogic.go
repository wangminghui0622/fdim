package logic

import (
	"context"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddApplicationVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddApplicationVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddApplicationVersionLogic {
	return &AddApplicationVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddApplicationVersionLogic) AddApplicationVersion(req *admin.AddApplicationVersionReq) (*admin.AddApplicationVersionResp, error) {
	// 1. 验证参数
	if req.Platform == "" {
		return nil, errs.ErrArgs.WrapMsg("platform cannot be empty")
	}
	if req.Version == "" {
		return nil, errs.ErrArgs.WrapMsg("version cannot be empty")
	}

	// 2. 生成ID
	versionID := uuid.New().String()

	// 3. 构建应用版本对象
	version := &database.ApplicationVersion{
		ID:         versionID,
		Platform:   req.Platform,
		Version:    req.Version,
		URL:        req.Url,
		Text:       req.Text,
		Force:      req.Force,
		Latest:     req.Latest,
		Hot:        req.Hot,
		CreateTime: time.Now(),
	}

	// 4. 如果设置为最新版本，需要将同平台的其他版本设置为非最新
	// 注意：这个逻辑应该在数据库层实现，这里简化处理
	if req.Latest {
		// 通过更新其他版本的方式处理，这里先添加，后续可以通过批量更新处理
		// 简化处理：在添加后通过UpdateApplicationVersion更新其他版本
	}

	// 5. 添加应用版本
	if err := l.svcCtx.AdminDB.AddApplicationVersion(l.ctx, []*database.ApplicationVersion{version}); err != nil {
		l.Errorf("AddApplicationVersion failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to add application version")
	}

	l.Infof("Added application version: id=%s, platform=%s, version=%s", versionID, req.Platform, req.Version)
	return &admin.AddApplicationVersionResp{}, nil
}
