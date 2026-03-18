package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PageApplicationVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageApplicationVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageApplicationVersionLogic {
	return &PageApplicationVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageApplicationVersionLogic) PageApplicationVersion(req *admin.PageApplicationVersionReq) (*admin.PageApplicationVersionResp, error) {
	// 1. 验证分页参数
	if req.Pagination == nil || req.Pagination.PageNumber <= 0 || req.Pagination.ShowNumber <= 0 {
		return nil, errs.ErrArgs.WrapMsg("invalid pagination parameters")
	}

	// 2. 分页获取应用版本
	total, versions, err := l.svcCtx.AdminDB.PageApplicationVersion(l.ctx, req.Platform, req.Pagination)
	if err != nil {
		l.Errorf("PageApplicationVersion failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to page application versions")
	}

	// 3. 转换为响应格式
	results := make([]*admin.ApplicationVersion, 0, len(versions))
	for _, version := range versions {
		results = append(results, &admin.ApplicationVersion{
			Id:         version.ID,
			Platform:   version.Platform,
			Version:    version.Version,
			Url:        version.URL,
			Text:       version.Text,
			Force:      version.Force,
			Latest:     version.Latest,
			Hot:        version.Hot,
			CreateTime: version.CreateTime.UnixMilli(),
		})
	}

	return &admin.PageApplicationVersionResp{
		Total:    int64(total),
		Versions: results,
	}, nil
}
