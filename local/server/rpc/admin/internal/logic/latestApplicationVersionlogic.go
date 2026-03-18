package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type LatestApplicationVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLatestApplicationVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LatestApplicationVersionLogic {
	return &LatestApplicationVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LatestApplicationVersionLogic) LatestApplicationVersion(req *admin.LatestApplicationVersionReq) (*admin.LatestApplicationVersionResp, error) {
	// 1. 验证参数
	if req.Platform == "" {
		return nil, errs.ErrArgs.WrapMsg("platform cannot be empty")
	}

	// 2. 获取最新版本
	version, err := l.svcCtx.AdminDB.LatestVersion(l.ctx, req.Platform)
	if err != nil {
		l.Errorf("LatestVersion failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to get latest version")
	}

	// 3. 转换为响应格式
	return &admin.LatestApplicationVersionResp{
		Version: &admin.ApplicationVersion{
			Id:         version.ID,
			Platform:   version.Platform,
			Version:    version.Version,
			Url:        version.URL,
			Text:       version.Text,
			Force:      version.Force,
			Latest:     version.Latest,
			Hot:        version.Hot,
			CreateTime: version.CreateTime.UnixMilli(),
		},
	}, nil
}
