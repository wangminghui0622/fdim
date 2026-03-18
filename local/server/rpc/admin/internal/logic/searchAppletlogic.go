package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/protocol/sdkws"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchAppletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchAppletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAppletLogic {
	return &SearchAppletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchAppletLogic) SearchApplet(req *admin.SearchAppletReq) (*admin.SearchAppletResp, error) {
	// 1. 验证分页参数
	if req.Pagination == nil || req.Pagination.PageNumber <= 0 || req.Pagination.ShowNumber <= 0 {
		return nil, errs.ErrArgs.WrapMsg("invalid pagination parameters")
	}

	// 2. 搜索小程序
	total, applets, err := l.svcCtx.AdminDB.SearchApplet(l.ctx, req.Keyword, req.Pagination)
	if err != nil {
		l.Errorf("SearchApplet failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to search applets")
	}

	// 3. 转换为响应格式
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

	return &admin.SearchAppletResp{
		Total:   uint32(total),
		Applets: results,
	}, nil
}
