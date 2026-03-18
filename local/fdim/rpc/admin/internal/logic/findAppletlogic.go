package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/protocol/sdkws"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindAppletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAppletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAppletLogic {
	return &FindAppletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindAppletLogic) FindApplet(req *admin.FindAppletReq) (*admin.FindAppletResp, error) {
	// FindApplet 返回所有已上架的小程序（status=1）
	applets, err := l.svcCtx.AdminDB.FindOnShelfApplet(l.ctx)
	if err != nil {
		l.Errorf("FindOnShelfApplet failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to find on shelf applets")
	}

	// 转换为响应格式
	results := make([]*sdkws.ChatAppletInfo, 0, len(applets))
	for _, applet := range applets {
		results = append(results, &sdkws.ChatAppletInfo{
			Id:         applet.ID,
			Name:       applet.Name,
			AppID:      applet.AppID,
			Icon:       applet.Icon,
			Url:        applet.URL,
			Md5:        applet.MD5,
			Size:       applet.Size,
			Version:    applet.Version,
			Priority:   applet.Priority,
			Status:     applet.Status,
			CreateTime: applet.CreateTime.UnixMilli(),
		})
	}

	return &admin.FindAppletResp{
		Applets: results,
	}, nil
}
