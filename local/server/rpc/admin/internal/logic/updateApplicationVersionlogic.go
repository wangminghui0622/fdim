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
	// 1. ֤
	if req.Id == "" {
		return nil, errs.ErrArgs.WrapMsg("id cannot be empty")
	}

	// 2. ֶ
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
		// Ϊ°汾Ҫͬƽ̨汾Ϊ
		if req.Latest.Value {
			// ȡǰ汾ϢԻȡƽ̨
			// 򻯴ȸµǰ汾Ȼͬƽ̨汾
		}
	}
	if req.Hot != nil {
		update["hot"] = req.Hot.Value
	}

	if len(update) == 0 {
		return nil, errs.ErrArgs.WrapMsg("no update fields provided")
	}

	// 3. Ӧð汾
	if err := l.svcCtx.AdminDB.UpdateApplicationVersion(l.ctx, req.Id, update); err != nil {
		l.Errorf("UpdateApplicationVersion failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to update application version")
	}

	l.Infof("Updated application version: id=%s", req.Id)
	return &admin.UpdateApplicationVersionResp{}, nil
}
