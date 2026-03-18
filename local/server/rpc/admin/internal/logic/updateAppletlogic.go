package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppletLogic {
	return &UpdateAppletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppletLogic) UpdateApplet(req *admin.UpdateAppletReq) (*admin.UpdateAppletResp, error) {
	// 1. 验证参数
	if req.Id == "" {
		return nil, errs.ErrArgs.WrapMsg("id cannot be empty")
	}

	// 2. 构建更新字段
	update := make(map[string]interface{})
	if req.Name != nil {
		update["name"] = req.Name.Value
	}
	if req.AppID != nil {
		update["app_id"] = req.AppID.Value
	}
	if req.Icon != nil {
		update["icon"] = req.Icon.Value
	}
	if req.Url != nil {
		update["url"] = req.Url.Value
	}
	if req.Md5 != nil {
		update["md5"] = req.Md5.Value
	}
	if req.Size != nil {
		update["size"] = req.Size.Value
	}
	if req.Version != nil {
		update["version"] = req.Version.Value
	}
	if req.Priority != nil {
		update["priority"] = req.Priority.Value
	}
	if req.Status != nil {
		update["status"] = req.Status.Value
	}

	if len(update) == 0 {
		return nil, errs.ErrArgs.WrapMsg("no update fields provided")
	}

	// 3. 更新小程序
	if err := l.svcCtx.AdminDB.UpdateApplet(l.ctx, req.Id, update); err != nil {
		l.Errorf("UpdateApplet failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to update applet")
	}

	l.Infof("Updated applet: id=%s", req.Id)
	return &admin.UpdateAppletResp{}, nil
}
