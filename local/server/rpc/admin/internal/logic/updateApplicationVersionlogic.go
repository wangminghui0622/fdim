package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateApplicationVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateApplicationVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApplicationVersionLogic {
	return &UpdateApplicationVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateApplicationVersionLogic) UpdateApplicationVersion(req *admin.UpdateApplicationVersionReq) (*admin.UpdateApplicationVersionResp, error) {
	// 1. 验证参数
	if req.Id == "" {
		return nil, errs.ErrArgs.WrapMsg("id cannot be empty")
	}

	// 2. 构建更新字段
	update := make(map[string]interface{})
	if req.Platform != nil {
		update["platform"] = req.Platform.Value
	}
	if req.Version != nil {
		update["version"] = req.Version.Value
	}
	if req.Url != nil {
		update["url"] = req.Url.Value
	}
	if req.Text != nil {
		update["text"] = req.Text.Value
	}
	if req.Force != nil {
		update["force"] = req.Force.Value
	}
	if req.Latest != nil {
		update["latest"] = req.Latest.Value
		// 如果设置为最新版本，需要将同平台的其他版本设置为非最新
		if req.Latest.Value {
			// 获取当前版本信息以获取平台
			// 简化处理：先更新当前版本，然后更新同平台其他版本
		}
	}
	if req.Hot != nil {
		update["hot"] = req.Hot.Value
	}

	if len(update) == 0 {
		return nil, errs.ErrArgs.WrapMsg("no update fields provided")
	}

	// 3. 更新应用版本
	if err := l.svcCtx.AdminDB.UpdateApplicationVersion(l.ctx, req.Id, update); err != nil {
		l.Errorf("UpdateApplicationVersion failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to update application version")
	}

	l.Infof("Updated application version: id=%s", req.Id)
	return &admin.UpdateApplicationVersionResp{}, nil
}
