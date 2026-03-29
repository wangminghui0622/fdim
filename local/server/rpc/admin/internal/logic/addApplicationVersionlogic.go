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
	// 1. ֤
	if req.Platform == "" {
		return nil, errs.ErrArgs.WrapMsg("platform cannot be empty")
	}
	if req.Version == "" {
		return nil, errs.ErrArgs.WrapMsg("version cannot be empty")
	}

	// 2. ID
	versionID := uuid.New().String()

	// 3. Ӧð汾
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

	// 4. Ϊ°汾Ҫͬƽ̨汾Ϊ
	// ע⣺߼Ӧݿʵ֣򻯴
	if req.Latest {
		// ͨ汾ķʽӣͨ´
		// 򻯴ӺͨUpdateApplicationVersion汾
	}

	// 5. Ӧð汾
	if err := l.svcCtx.AdminDB.AddApplicationVersion(l.ctx, []*database.ApplicationVersion{version}); err != nil {
		l.Errorf("AddApplicationVersion failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to add application version")
	}

	l.Infof("Added application version: id=%s, platform=%s, version=%s", versionID, req.Platform, req.Version)
	return &admin.AddApplicationVersionResp{}, nil
}
