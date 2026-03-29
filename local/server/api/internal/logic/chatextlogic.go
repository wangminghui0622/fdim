package logic

import (
	"context"
	"fmt"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/admin"
	"fdim/protocol/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

// ========== FindUserPublicInfo ==========

type FindUserPublicInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserPublicInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserPublicInfoLogic {
	return &FindUserPublicInfoLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FindUserPublicInfoLogic) FindUserPublicInfo(req *types.FindUserPublicInfoReq) (*types.FindUserPublicInfoResp, error) {
	if l.svcCtx.UserDB == nil {
		return &types.FindUserPublicInfoResp{Users: []interface{}{}}, nil
	}
	users, err := l.svcCtx.UserDB.Find(l.ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	for _, u := range users {
		result = append(result, map[string]interface{}{
			"userID":   u.UserID,
			"nickname": u.Nickname,
			"faceURL":  u.FaceURL,
		})
	}
	return &types.FindUserPublicInfoResp{Users: result}, nil
}

// ========== SearchUserPublicInfo ==========

type SearchUserPublicInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserPublicInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserPublicInfoLogic {
	return &SearchUserPublicInfoLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SearchUserPublicInfoLogic) SearchUserPublicInfo(req *types.SearchUserPublicInfoReq) (*types.SearchUserPublicInfoResp, error) {
	if l.svcCtx.UserDB == nil {
		return &types.SearchUserPublicInfoResp{Total: 0, Users: []interface{}{}}, nil
	}
	offset := (req.Pagination.PageNumber - 1) * req.Pagination.ShowNumber
	total, users, err := l.svcCtx.UserDB.SearchUsers(l.ctx, req.Keyword, offset, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	for _, u := range users {
		result = append(result, map[string]interface{}{
			"userID":   u.UserID,
			"nickname": u.Nickname,
			"faceURL":  u.FaceURL,
		})
	}
	return &types.SearchUserPublicInfoResp{Total: int32(total), Users: result}, nil
}

// ========== GetTokenForVideoMeeting ==========

type GetTokenForVideoMeetingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenForVideoMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenForVideoMeetingLogic {
	return &GetTokenForVideoMeetingLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetTokenForVideoMeetingLogic) GetTokenForVideoMeeting(req *types.GetTokenForVideoMeetingReq) (*types.GetTokenForVideoMeetingResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, fmt.Errorf("chat service not available")
	}
	rpcResp, err := l.svcCtx.ChatClient.GetTokenForVideoMeeting(l.ctx, &chat.GetTokenForVideoMeetingReq{
		Room:     req.Room,
		Identity: req.Identity,
	})
	if err != nil {
		return nil, err
	}
	return &types.GetTokenForVideoMeetingResp{
		ServerUrl: rpcResp.ServerUrl,
		Token:     rpcResp.Token,
	}, nil
}

// ========== FindApplet ==========

type FindAppletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAppletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAppletLogic {
	return &FindAppletLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FindAppletLogic) FindApplet() (*types.FindAppletResp, error) {
	if l.svcCtx.AdminClient == nil {
		return &types.FindAppletResp{Applets: []types.AppletInfo{}}, nil
	}
	rpcResp, err := l.svcCtx.AdminClient.FindApplet(l.ctx, &admin.FindAppletReq{})
	if err != nil {
		return nil, err
	}
	applets := make([]types.AppletInfo, 0, len(rpcResp.Applets))
	for _, a := range rpcResp.Applets {
		applets = append(applets, types.AppletInfo{
			ID:         a.Id,
			Name:       a.Name,
			AppID:      a.AppID,
			Icon:       a.Icon,
			Url:        a.Url,
			Md5:        a.Md5,
			Size:       a.Size,
			Version:    a.Version,
			Priority:   int32(a.Priority),
			Status:     int32(a.Status),
			CreateTime: a.CreateTime,
		})
	}
	return &types.FindAppletResp{Applets: applets}, nil
}

// ========== LatestApplicationVersion ==========

type LatestApplicationVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLatestApplicationVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LatestApplicationVersionLogic {
	return &LatestApplicationVersionLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *LatestApplicationVersionLogic) LatestApplicationVersion(req *types.LatestApplicationVersionReq) (*types.LatestApplicationVersionResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcResp, err := l.svcCtx.AdminClient.LatestApplicationVersion(l.ctx, &admin.LatestApplicationVersionReq{
		Platform: req.Platform,
	})
	if err != nil {
		return &types.LatestApplicationVersionResp{Version: nil}, nil
	}
	if rpcResp.Version == nil {
		return &types.LatestApplicationVersionResp{Version: nil}, nil
	}
	v := rpcResp.Version
	return &types.LatestApplicationVersionResp{
		Version: &types.ApplicationVersion{
			ID:         v.Id,
			Platform:   v.Platform,
			Version:    v.Version,
			Url:        v.Url,
			Text:       v.Text,
			IsForce:    v.Force,
			Latest:     v.Latest,
			Hot:        v.Hot,
			CreateTime: v.CreateTime,
		},
	}, nil
}

// ========== PageApplicationVersion ==========

type PageApplicationVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageApplicationVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageApplicationVersionLogic {
	return &PageApplicationVersionLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *PageApplicationVersionLogic) PageApplicationVersion(req *types.PageApplicationVersionReq) (*types.PageApplicationVersionResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	platforms := []string{}
	if req.Platform != "" {
		platforms = append(platforms, req.Platform)
	}
	rpcResp, err := l.svcCtx.AdminClient.PageApplicationVersion(l.ctx, &admin.PageApplicationVersionReq{
		Platform:   platforms,
		Pagination: toPagination(req.Pagination),
	})
	if err != nil {
		return nil, err
	}
	versions := make([]types.ApplicationVersion, 0, len(rpcResp.Versions))
	for _, v := range rpcResp.Versions {
		versions = append(versions, types.ApplicationVersion{
			ID:         v.Id,
			Platform:   v.Platform,
			Version:    v.Version,
			Url:        v.Url,
			Text:       v.Text,
			IsForce:    v.Force,
			Latest:     v.Latest,
			Hot:        v.Hot,
			CreateTime: v.CreateTime,
		})
	}
	return &types.PageApplicationVersionResp{Total: rpcResp.Total, Versions: versions}, nil
}

// ========== FDIMCallback ==========

type FDIMCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFDIMCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FDIMCallbackLogic {
	return &FDIMCallbackLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FDIMCallbackLogic) FDIMCallback(req *types.FDIMCallbackReq) (*types.FDIMCallbackResp, error) {
	l.Infof("FDIM callback: command=%s", req.Command)
	// Process callback from FDIM server (message hooks, user hooks, etc.)
	return &types.FDIMCallbackResp{}, nil
}
