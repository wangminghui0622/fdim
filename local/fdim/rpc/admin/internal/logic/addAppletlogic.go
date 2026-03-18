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

type AddAppletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAppletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAppletLogic {
	return &AddAppletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAppletLogic) AddApplet(req *admin.AddAppletReq) (*admin.AddAppletResp, error) {
	// 1. 验证参数
	if req.Name == "" {
		return nil, errs.ErrArgs.WrapMsg("name cannot be empty")
	}
	if req.AppID == "" {
		return nil, errs.ErrArgs.WrapMsg("appID cannot be empty")
	}

	// 2. 生成ID（如果未提供）
	appletID := req.Id
	if appletID == "" {
		appletID = uuid.New().String()
	}

	// 3. 构建小程序对象
	applet := &database.Applet{
		ID:         appletID,
		Name:       req.Name,
		AppID:      req.AppID,
		Icon:       req.Icon,
		URL:        req.Url,
		MD5:        req.Md5,
		Size:       req.Size,
		Version:    req.Version,
		Priority:   req.Priority,
		Status:     req.Status,
		CreateTime: time.Now(),
	}

	// 4. 添加小程序
	if err := l.svcCtx.AdminDB.AddApplet(l.ctx, []*database.Applet{applet}); err != nil {
		l.Errorf("AddApplet failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to add applet")
	}

	l.Infof("Added applet: id=%s, name=%s", appletID, req.Name)
	return &admin.AddAppletResp{}, nil
}
