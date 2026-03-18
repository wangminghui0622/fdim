package logic

import (
	"context"
	"fmt"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/admin"
	"fdim/protocol/chat"
	"fdim/protocol/sdkws"
	"fdim/protocol/wrapperspb"

	"github.com/zeromicro/go-zero/core/logx"
)

// toPagination converts API pagination to RPC pagination
func toPagination(p types.Pagination) *sdkws.RequestPagination {
	return &sdkws.RequestPagination{
		PageNumber: int32(p.PageNumber),
		ShowNumber: int32(p.ShowNumber),
	}
}

// ========== AdminLogin ==========

type AdminLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLoginLogic {
	return &AdminLoginLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AdminLoginLogic) AdminLogin(req *types.AdminLoginReq) (*types.AdminLoginResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcResp, err := l.svcCtx.AdminClient.Login(l.ctx, &admin.LoginReq{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &types.AdminLoginResp{
		AdminToken:  rpcResp.AdminToken,
		AdminUserID: rpcResp.AdminUserID,
		Nickname:    rpcResp.Nickname,
		FaceURL:     rpcResp.FaceURL,
		Level:       rpcResp.Level,
	}, nil
}

// ========== AdminUpdateInfo ==========

type AdminUpdateInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminUpdateInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateInfoLogic {
	return &AdminUpdateInfoLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AdminUpdateInfoLogic) AdminUpdateInfo(req *types.AdminUpdateInfoReq) (*types.AdminUpdateInfoResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcReq := &admin.AdminUpdateInfoReq{}
	if req.Account != nil {
		rpcReq.Account = wrapperspb.String(*req.Account)
	}
	if req.Password != nil {
		rpcReq.Password = wrapperspb.String(*req.Password)
	}
	if req.FaceURL != nil {
		rpcReq.FaceURL = wrapperspb.String(*req.FaceURL)
	}
	if req.Nickname != nil {
		rpcReq.Nickname = wrapperspb.String(*req.Nickname)
	}
	_, err := l.svcCtx.AdminClient.AdminUpdateInfo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.AdminUpdateInfoResp{}, nil
}

// ========== AdminInfo ==========

type AdminInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminInfoLogic {
	return &AdminInfoLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AdminInfoLogic) AdminInfo() (*types.AdminInfoResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcResp, err := l.svcCtx.AdminClient.GetAdminInfo(l.ctx, &admin.GetAdminInfoReq{})
	if err != nil {
		return nil, err
	}
	return &types.AdminInfoResp{
		Account:  rpcResp.Account,
		Nickname: rpcResp.Nickname,
		FaceURL:  rpcResp.FaceURL,
		Level:    rpcResp.Level,
	}, nil
}

// ========== ChangeAdminPassword ==========

type ChangeAdminPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeAdminPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeAdminPasswordLogic {
	return &ChangeAdminPasswordLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ChangeAdminPasswordLogic) ChangeAdminPassword(req *types.ChangeAdminPasswordReq) (*types.ChangeAdminPasswordResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.ChangeAdminPassword(l.ctx, &admin.ChangeAdminPasswordReq{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	return &types.ChangeAdminPasswordResp{}, nil
}

// ========== AddAdminAccount ==========

type AddAdminAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAdminAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAdminAccountLogic {
	return &AddAdminAccountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddAdminAccountLogic) AddAdminAccount(req *types.AddAdminAccountReq) (*types.AddAdminAccountResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.AddAdminAccount(l.ctx, &admin.AddAdminAccountReq{
		Account:  req.Account,
		Password: req.Password,
		FaceURL:  req.FaceURL,
		Nickname: req.Nickname,
	})
	if err != nil {
		return nil, err
	}
	return &types.AddAdminAccountResp{}, nil
}

// ========== DelAdminAccount ==========

type DelAdminAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelAdminAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelAdminAccountLogic {
	return &DelAdminAccountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DelAdminAccountLogic) DelAdminAccount(req *types.DelAdminAccountReq) (*types.DelAdminAccountResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.DelAdminAccount(l.ctx, &admin.DelAdminAccountReq{
		UserIDs: req.UserIDs,
	})
	if err != nil {
		return nil, err
	}
	return &types.DelAdminAccountResp{}, nil
}

// ========== SearchAdminAccount ==========

type SearchAdminAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchAdminAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAdminAccountLogic {
	return &SearchAdminAccountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SearchAdminAccountLogic) SearchAdminAccount(req *types.SearchAdminAccountReq) (*types.SearchAdminAccountResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcResp, err := l.svcCtx.AdminClient.SearchAdminAccount(l.ctx, &admin.SearchAdminAccountReq{
		Pagination: toPagination(req.Pagination),
	})
	if err != nil {
		return nil, err
	}
	admins := make([]types.AdminAccountInfo, 0, len(rpcResp.AdminAccounts))
	for _, a := range rpcResp.AdminAccounts {
		admins = append(admins, types.AdminAccountInfo{
			UserID:   a.UserID,
			Account:  a.Account,
			Nickname: a.Nickname,
			FaceURL:  a.FaceURL,
			Level:    a.Level,
		})
	}
	return &types.SearchAdminAccountResp{Total: int64(rpcResp.Total), Admins: admins}, nil
}

// ========== AddUserAccount ==========

type AddUserAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserAccountLogic {
	return &AddUserAccountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddUserAccountLogic) AddUserAccount(req *types.AddUserAccountReq) (*types.AddUserAccountResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, fmt.Errorf("chat service not available")
	}
	// Register user via chat RPC
	rpcResp, err := l.svcCtx.ChatClient.RegisterUser(l.ctx, &chat.RegisterUserReq{
		User: &chat.RegisterUserInfo{
			Nickname:    req.User.Nickname,
			FaceURL:     req.User.FaceURL,
			PhoneNumber: req.User.PhoneNumber,
			Email:       req.User.Email,
			Account:     req.User.Account,
			Password:    req.User.Password,
			AreaCode:    req.User.AreaCode,
		},
	})
	if err != nil {
		return nil, err
	}
	return &types.AddUserAccountResp{UserID: rpcResp.UserID}, nil
}

// ========== ImportUserByJson ==========

type ImportUserByJsonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewImportUserByJsonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportUserByJsonLogic {
	return &ImportUserByJsonLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ImportUserByJsonLogic) ImportUserByJson(req *types.ImportUserByJsonReq) (*types.ImportUserByJsonResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, fmt.Errorf("chat service not available")
	}
	userIDs := make([]string, 0, len(req.Users))
	for _, u := range req.Users {
		rpcResp, err := l.svcCtx.ChatClient.RegisterUser(l.ctx, &chat.RegisterUserReq{
			User: &chat.RegisterUserInfo{
				Nickname:    u.Nickname,
				FaceURL:     u.FaceURL,
				PhoneNumber: u.PhoneNumber,
				Email:       u.Email,
				Account:     u.Account,
				Password:    u.Password,
				AreaCode:    u.AreaCode,
			},
		})
		if err != nil {
			l.Errorf("ImportUserByJson: failed to register user %s: %v", u.Nickname, err)
			continue
		}
		userIDs = append(userIDs, rpcResp.UserID)
	}
	return &types.ImportUserByJsonResp{UserIDs: userIDs}, nil
}

// ========== GetAllowRegister ==========

type GetAllowRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllowRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllowRegisterLogic {
	return &GetAllowRegisterLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetAllowRegisterLogic) GetAllowRegister() (*types.GetAllowRegisterResp, error) {
	if l.svcCtx.AdminClient == nil {
		return &types.GetAllowRegisterResp{AllowRegister: true}, nil
	}
	rpcResp, err := l.svcCtx.AdminClient.GetClientConfig(l.ctx, &admin.GetClientConfigReq{})
	if err != nil {
		return &types.GetAllowRegisterResp{AllowRegister: true}, nil
	}
	allowRegister := true
	if v, ok := rpcResp.Config["allowRegister"]; ok && v == "false" {
		allowRegister = false
	}
	return &types.GetAllowRegisterResp{AllowRegister: allowRegister}, nil
}

// ========== SetAllowRegister ==========

type SetAllowRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetAllowRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetAllowRegisterLogic {
	return &SetAllowRegisterLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SetAllowRegisterLogic) SetAllowRegister(req *types.SetAllowRegisterReq) (*types.SetAllowRegisterResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	val := "true"
	if !req.AllowRegister {
		val = "false"
	}
	_, err := l.svcCtx.AdminClient.SetClientConfig(l.ctx, &admin.SetClientConfigReq{
		Config: map[string]string{"allowRegister": val},
	})
	if err != nil {
		return nil, err
	}
	return &types.SetAllowRegisterResp{}, nil
}

// ========== Default Friend ==========

type AddDefaultFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddDefaultFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDefaultFriendLogic {
	return &AddDefaultFriendLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddDefaultFriendLogic) AddDefaultFriend(req *types.AddDefaultFriendReq) (*types.AddDefaultFriendResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.AddDefaultFriend(l.ctx, &admin.AddDefaultFriendReq{UserIDs: req.UserIDs})
	if err != nil {
		return nil, err
	}
	return &types.AddDefaultFriendResp{}, nil
}

type DelDefaultFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelDefaultFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelDefaultFriendLogic {
	return &DelDefaultFriendLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DelDefaultFriendLogic) DelDefaultFriend(req *types.DelDefaultFriendReq) (*types.DelDefaultFriendResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.DelDefaultFriend(l.ctx, &admin.DelDefaultFriendReq{UserIDs: req.UserIDs})
	if err != nil {
		return nil, err
	}
	return &types.DelDefaultFriendResp{}, nil
}

type FindDefaultFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindDefaultFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindDefaultFriendLogic {
	return &FindDefaultFriendLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FindDefaultFriendLogic) FindDefaultFriend() (*types.FindDefaultFriendResp, error) {
	if l.svcCtx.AdminClient == nil {
		return &types.FindDefaultFriendResp{UserIDs: []string{}}, nil
	}
	rpcResp, err := l.svcCtx.AdminClient.FindDefaultFriend(l.ctx, &admin.FindDefaultFriendReq{})
	if err != nil {
		return nil, err
	}
	return &types.FindDefaultFriendResp{UserIDs: rpcResp.UserIDs}, nil
}

type SearchDefaultFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchDefaultFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchDefaultFriendLogic {
	return &SearchDefaultFriendLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SearchDefaultFriendLogic) SearchDefaultFriend(req *types.SearchDefaultFriendReq) (*types.SearchDefaultFriendResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcResp, err := l.svcCtx.AdminClient.SearchDefaultFriend(l.ctx, &admin.SearchDefaultFriendReq{
		Keyword:    req.Keyword,
		Pagination: toPagination(req.Pagination),
	})
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, 0, len(rpcResp.Users))
	for _, u := range rpcResp.Users {
		userIDs = append(userIDs, u.UserID)
	}
	return &types.SearchDefaultFriendResp{Total: int64(rpcResp.Total), UserIDs: userIDs}, nil
}

// ========== Default Group ==========

type AddDefaultGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddDefaultGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDefaultGroupLogic {
	return &AddDefaultGroupLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddDefaultGroupLogic) AddDefaultGroup(req *types.AddDefaultGroupReq) (*types.AddDefaultGroupResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.AddDefaultGroup(l.ctx, &admin.AddDefaultGroupReq{GroupIDs: req.GroupIDs})
	if err != nil {
		return nil, err
	}
	return &types.AddDefaultGroupResp{}, nil
}

type DelDefaultGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelDefaultGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelDefaultGroupLogic {
	return &DelDefaultGroupLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DelDefaultGroupLogic) DelDefaultGroup(req *types.DelDefaultGroupReq) (*types.DelDefaultGroupResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.DelDefaultGroup(l.ctx, &admin.DelDefaultGroupReq{GroupIDs: req.GroupIDs})
	if err != nil {
		return nil, err
	}
	return &types.DelDefaultGroupResp{}, nil
}

type FindDefaultGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindDefaultGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindDefaultGroupLogic {
	return &FindDefaultGroupLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *FindDefaultGroupLogic) FindDefaultGroup() (*types.FindDefaultGroupResp, error) {
	if l.svcCtx.AdminClient == nil {
		return &types.FindDefaultGroupResp{GroupIDs: []string{}}, nil
	}
	rpcResp, err := l.svcCtx.AdminClient.FindDefaultGroup(l.ctx, &admin.FindDefaultGroupReq{})
	if err != nil {
		return nil, err
	}
	return &types.FindDefaultGroupResp{GroupIDs: rpcResp.GroupIDs}, nil
}

type SearchDefaultGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchDefaultGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchDefaultGroupLogic {
	return &SearchDefaultGroupLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SearchDefaultGroupLogic) SearchDefaultGroup(req *types.SearchDefaultGroupReq) (*types.SearchDefaultGroupResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcResp, err := l.svcCtx.AdminClient.SearchDefaultGroup(l.ctx, &admin.SearchDefaultGroupReq{
		Keyword:    req.Keyword,
		Pagination: toPagination(req.Pagination),
	})
	if err != nil {
		return nil, err
	}
	return &types.SearchDefaultGroupResp{Total: int64(rpcResp.Total), GroupIDs: rpcResp.GroupIDs}, nil
}

// ========== Invitation Code ==========

type AddInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddInvitationCodeLogic {
	return &AddInvitationCodeLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddInvitationCodeLogic) AddInvitationCode(req *types.AddInvitationCodeReq) (*types.AddInvitationCodeResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.AddInvitationCode(l.ctx, &admin.AddInvitationCodeReq{Codes: req.Codes})
	if err != nil {
		return nil, err
	}
	return &types.AddInvitationCodeResp{}, nil
}

type GenInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenInvitationCodeLogic {
	return &GenInvitationCodeLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GenInvitationCodeLogic) GenInvitationCode(req *types.GenInvitationCodeReq) (*types.GenInvitationCodeResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.GenInvitationCode(l.ctx, &admin.GenInvitationCodeReq{
		Len: req.Len,
		Num: req.Num,
	})
	if err != nil {
		return nil, err
	}
	// GenInvitationCodeResp has no Codes field in proto; codes are stored server-side
	return &types.GenInvitationCodeResp{Codes: []string{}}, nil
}

type DelInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelInvitationCodeLogic {
	return &DelInvitationCodeLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DelInvitationCodeLogic) DelInvitationCode(req *types.DelInvitationCodeReq) (*types.DelInvitationCodeResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.DelInvitationCode(l.ctx, &admin.DelInvitationCodeReq{Codes: req.Codes})
	if err != nil {
		return nil, err
	}
	return &types.DelInvitationCodeResp{}, nil
}

type SearchInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchInvitationCodeLogic {
	return &SearchInvitationCodeLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SearchInvitationCodeLogic) SearchInvitationCode(req *types.SearchInvitationCodeReq) (*types.SearchInvitationCodeResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	var status int32
	if req.Status != nil {
		status = *req.Status
	}
	rpcResp, err := l.svcCtx.AdminClient.SearchInvitationCode(l.ctx, &admin.SearchInvitationCodeReq{
		Status:     status,
		Keyword:    req.Keyword,
		Pagination: toPagination(req.Pagination),
	})
	if err != nil {
		return nil, err
	}
	codes := make([]types.InvitationCodeInfo, 0, len(rpcResp.List))
	for _, c := range rpcResp.List {
		codes = append(codes, types.InvitationCodeInfo{
			Code:       c.InvitationCode,
			CreateTime: c.CreateTime,
			UsedUserID: c.UsedUserID,
		})
	}
	return &types.SearchInvitationCodeResp{Total: int64(rpcResp.Total), Codes: codes}, nil
}

// ========== IP Forbidden ==========

type AddIPForbiddenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddIPForbiddenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddIPForbiddenLogic {
	return &AddIPForbiddenLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddIPForbiddenLogic) AddIPForbidden(req *types.AddIPForbiddenReq) (*types.AddIPForbiddenResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	forbiddens := make([]*admin.IPForbiddenAdd, 0, len(req.Ips))
	for _, ip := range req.Ips {
		forbiddens = append(forbiddens, &admin.IPForbiddenAdd{Ip: ip, LimitRegister: true, LimitLogin: true})
	}
	_, err := l.svcCtx.AdminClient.AddIPForbidden(l.ctx, &admin.AddIPForbiddenReq{Forbiddens: forbiddens})
	if err != nil {
		return nil, err
	}
	return &types.AddIPForbiddenResp{}, nil
}

type DelIPForbiddenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelIPForbiddenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelIPForbiddenLogic {
	return &DelIPForbiddenLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DelIPForbiddenLogic) DelIPForbidden(req *types.DelIPForbiddenReq) (*types.DelIPForbiddenResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.DelIPForbidden(l.ctx, &admin.DelIPForbiddenReq{Ips: req.Ips})
	if err != nil {
		return nil, err
	}
	return &types.DelIPForbiddenResp{}, nil
}

type SearchIPForbiddenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchIPForbiddenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchIPForbiddenLogic {
	return &SearchIPForbiddenLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SearchIPForbiddenLogic) SearchIPForbidden(req *types.SearchIPForbiddenReq) (*types.SearchIPForbiddenResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcResp, err := l.svcCtx.AdminClient.SearchIPForbidden(l.ctx, &admin.SearchIPForbiddenReq{
		Keyword:    req.Keyword,
		Pagination: toPagination(req.Pagination),
	})
	if err != nil {
		return nil, err
	}
	ips := make([]types.IPForbiddenInfo, 0, len(rpcResp.Forbiddens))
	for _, f := range rpcResp.Forbiddens {
		ips = append(ips, types.IPForbiddenInfo{
			Ip:            f.Ip,
			LimitLogin:    f.LimitLogin,
			LimitRegister: f.LimitRegister,
			CreateTime:    f.CreateTime,
		})
	}
	return &types.SearchIPForbiddenResp{Total: int64(rpcResp.Total), Ips: ips}, nil
}

// ========== User IP Limit ==========

type AddUserIPLimitLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserIPLimitLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserIPLimitLoginLogic {
	return &AddUserIPLimitLoginLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddUserIPLimitLoginLogic) AddUserIPLimitLogin(req *types.AddUserIPLimitLoginReq) (*types.AddUserIPLimitLoginResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	limits := make([]*admin.UserIPLimitLogin, 0, len(req.Ips))
	for _, ip := range req.Ips {
		limits = append(limits, &admin.UserIPLimitLogin{UserID: req.UserID, Ip: ip})
	}
	_, err := l.svcCtx.AdminClient.AddUserIPLimitLogin(l.ctx, &admin.AddUserIPLimitLoginReq{Limits: limits})
	if err != nil {
		return nil, err
	}
	return &types.AddUserIPLimitLoginResp{}, nil
}

type DelUserIPLimitLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelUserIPLimitLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserIPLimitLoginLogic {
	return &DelUserIPLimitLoginLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DelUserIPLimitLoginLogic) DelUserIPLimitLogin(req *types.DelUserIPLimitLoginReq) (*types.DelUserIPLimitLoginResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	limits := make([]*admin.UserIPLimitLogin, 0, len(req.Ips))
	for _, ip := range req.Ips {
		limits = append(limits, &admin.UserIPLimitLogin{UserID: req.UserID, Ip: ip})
	}
	_, err := l.svcCtx.AdminClient.DelUserIPLimitLogin(l.ctx, &admin.DelUserIPLimitLoginReq{Limits: limits})
	if err != nil {
		return nil, err
	}
	return &types.DelUserIPLimitLoginResp{}, nil
}

type SearchUserIPLimitLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserIPLimitLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserIPLimitLoginLogic {
	return &SearchUserIPLimitLoginLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SearchUserIPLimitLoginLogic) SearchUserIPLimitLogin(req *types.SearchUserIPLimitLoginReq) (*types.SearchUserIPLimitLoginResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcResp, err := l.svcCtx.AdminClient.SearchUserIPLimitLogin(l.ctx, &admin.SearchUserIPLimitLoginReq{
		Keyword:    req.Keyword,
		Pagination: toPagination(req.Pagination),
	})
	if err != nil {
		return nil, err
	}
	limits := make([]types.UserIPLimitInfo, 0, len(rpcResp.Limits))
	for _, lm := range rpcResp.Limits {
		limits = append(limits, types.UserIPLimitInfo{
			UserID:     lm.UserID,
			Ip:         lm.Ip,
			CreateTime: lm.CreateTime,
		})
	}
	return &types.SearchUserIPLimitLoginResp{Total: int64(rpcResp.Total), Limits: limits}, nil
}

// ========== Applet CRUD ==========

type AddAppletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAppletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAppletLogic {
	return &AddAppletLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddAppletLogic) AddApplet(req *types.AddAppletReq) (*types.AddAppletResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.AddApplet(l.ctx, &admin.AddAppletReq{
		Name:     req.Name,
		AppID:    req.AppID,
		Icon:     req.Icon,
		Url:      req.Url,
		Md5:      req.Md5,
		Size:     req.Size,
		Version:  req.Version,
		Priority: uint32(req.Priority),
	})
	if err != nil {
		return nil, err
	}
	return &types.AddAppletResp{}, nil
}

type DelAppletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelAppletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelAppletLogic {
	return &DelAppletLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DelAppletLogic) DelApplet(req *types.DelAppletReq) (*types.DelAppletResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.DelApplet(l.ctx, &admin.DelAppletReq{AppletIds: req.AppletIDs})
	if err != nil {
		return nil, err
	}
	return &types.DelAppletResp{}, nil
}

type UpdateAppletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppletLogic {
	return &UpdateAppletLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateAppletLogic) UpdateApplet(req *types.UpdateAppletReq) (*types.UpdateAppletResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcReq := &admin.UpdateAppletReq{Id: req.ID}
	if req.Name != nil {
		rpcReq.Name = wrapperspb.String(*req.Name)
	}
	if req.AppID != nil {
		rpcReq.AppID = wrapperspb.String(*req.AppID)
	}
	if req.Icon != nil {
		rpcReq.Icon = wrapperspb.String(*req.Icon)
	}
	if req.Url != nil {
		rpcReq.Url = wrapperspb.String(*req.Url)
	}
	if req.Md5 != nil {
		rpcReq.Md5 = wrapperspb.String(*req.Md5)
	}
	if req.Version != nil {
		rpcReq.Version = wrapperspb.String(*req.Version)
	}
	_, err := l.svcCtx.AdminClient.UpdateApplet(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.UpdateAppletResp{}, nil
}

type SearchAppletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchAppletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAppletLogic {
	return &SearchAppletLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SearchAppletLogic) SearchApplet(req *types.SearchAppletReq) (*types.SearchAppletResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcResp, err := l.svcCtx.AdminClient.SearchApplet(l.ctx, &admin.SearchAppletReq{
		Keyword:    req.Keyword,
		Pagination: toPagination(req.Pagination),
	})
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
	return &types.SearchAppletResp{Total: int64(rpcResp.Total), Applets: applets}, nil
}

// ========== Block User ==========

type BlockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockUserLogic {
	return &BlockUserLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *BlockUserLogic) BlockUser(req *types.BlockUserReq) (*types.BlockUserResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.BlockUser(l.ctx, &admin.BlockUserReq{
		UserID: req.UserID,
		Reason: req.Reason,
	})
	if err != nil {
		return nil, err
	}
	return &types.BlockUserResp{}, nil
}

type UnblockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnblockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnblockUserLogic {
	return &UnblockUserLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UnblockUserLogic) UnblockUser(req *types.UnblockUserReq) (*types.UnblockUserResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.UnblockUser(l.ctx, &admin.UnblockUserReq{UserIDs: req.UserIDs})
	if err != nil {
		return nil, err
	}
	return &types.UnblockUserResp{}, nil
}

type SearchBlockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchBlockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchBlockUserLogic {
	return &SearchBlockUserLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SearchBlockUserLogic) SearchBlockUser(req *types.SearchBlockUserReq) (*types.SearchBlockUserResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcResp, err := l.svcCtx.AdminClient.SearchBlockUser(l.ctx, &admin.SearchBlockUserReq{
		Keyword:    req.Keyword,
		Pagination: toPagination(req.Pagination),
	})
	if err != nil {
		return nil, err
	}
	users := make([]types.BlockUserInfo, 0, len(rpcResp.Users))
	for _, u := range rpcResp.Users {
		users = append(users, types.BlockUserInfo{
			UserID:     u.UserID,
			Reason:     u.Reason,
			OpUserID:   u.OpUserID,
			CreateTime: u.CreateTime,
		})
	}
	return &types.SearchBlockUserResp{Total: int64(rpcResp.Total), Users: users}, nil
}

// ========== Admin Reset User Password ==========

type AdminResetUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminResetUserPasswordLogic {
	return &AdminResetUserPasswordLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AdminResetUserPasswordLogic) AdminResetUserPassword(req *types.AdminResetUserPasswordReq) (*types.AdminResetUserPasswordResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, fmt.Errorf("chat service not available")
	}
	_, err := l.svcCtx.ChatClient.ChangePassword(l.ctx, &chat.ChangePasswordReq{
		UserID:      req.UserID,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	return &types.AdminResetUserPasswordResp{}, nil
}

// ========== Admin Client Config ==========

type AdminSetClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminSetClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminSetClientConfigLogic {
	return &AdminSetClientConfigLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AdminSetClientConfigLogic) SetClientConfig(req *types.SetClientConfigReq) (*types.SetClientConfigResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.SetClientConfig(l.ctx, &admin.SetClientConfigReq{Config: req.Config})
	if err != nil {
		return nil, err
	}
	return &types.SetClientConfigResp{}, nil
}

type AdminDelClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminDelClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDelClientConfigLogic {
	return &AdminDelClientConfigLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AdminDelClientConfigLogic) DelClientConfig(req *types.DelClientConfigReq) (*types.DelClientConfigResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.DelClientConfig(l.ctx, &admin.DelClientConfigReq{Keys: req.Keys})
	if err != nil {
		return nil, err
	}
	return &types.DelClientConfigResp{}, nil
}

// ========== Admin Statistics ==========

type NewUserCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNewUserCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewUserCountLogic {
	return &NewUserCountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *NewUserCountLogic) NewUserCount(req *types.NewUserCountReq) (*types.NewUserCountResp, error) {
	// TODO: implement when chat RPC adds NewUserCount method
	return &types.NewUserCountResp{Total: 0, DateCounts: []types.DateCount{}}, nil
}

type LoginUserCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginUserCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginUserCountLogic {
	return &LoginUserCountLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *LoginUserCountLogic) LoginUserCount(req *types.LoginUserCountReq) (*types.LoginUserCountResp, error) {
	if l.svcCtx.ChatClient == nil {
		return &types.LoginUserCountResp{Total: 0, DateCounts: []types.DateCount{}}, nil
	}
	rpcResp, err := l.svcCtx.ChatClient.UserLoginCount(l.ctx, &chat.UserLoginCountReq{
		Start: req.Start,
		End:   req.End,
	})
	if err != nil {
		return &types.LoginUserCountResp{Total: 0, DateCounts: []types.DateCount{}}, nil
	}
	dateCounts := make([]types.DateCount, 0, len(rpcResp.Count))
	for date, count := range rpcResp.Count {
		dateCounts = append(dateCounts, types.DateCount{Date: date, Count: count})
	}
	return &types.LoginUserCountResp{Total: rpcResp.LoginCount, DateCounts: dateCounts}, nil
}

// ========== Application Version Management ==========

type AddApplicationVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddApplicationVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddApplicationVersionLogic {
	return &AddApplicationVersionLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddApplicationVersionLogic) AddApplicationVersion(req *types.AddApplicationVersionReq) (*types.AddApplicationVersionResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.AddApplicationVersion(l.ctx, &admin.AddApplicationVersionReq{
		Platform: req.Platform,
		Version:  req.Version,
		Url:      req.Url,
		Text:     req.Text,
		Force:    req.IsForce,
		Latest:   req.Latest,
		Hot:      req.Hot,
	})
	if err != nil {
		return nil, err
	}
	return &types.AddApplicationVersionResp{}, nil
}

type UpdateApplicationVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateApplicationVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApplicationVersionLogic {
	return &UpdateApplicationVersionLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateApplicationVersionLogic) UpdateApplicationVersion(req *types.UpdateApplicationVersionReq) (*types.UpdateApplicationVersionResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	rpcReq := &admin.UpdateApplicationVersionReq{Id: req.ID}
	if req.Platform != nil {
		rpcReq.Platform = wrapperspb.String(*req.Platform)
	}
	if req.Version != nil {
		rpcReq.Version = wrapperspb.String(*req.Version)
	}
	if req.Url != nil {
		rpcReq.Url = wrapperspb.String(*req.Url)
	}
	if req.Text != nil {
		rpcReq.Text = wrapperspb.String(*req.Text)
	}
	if req.IsForce != nil {
		rpcReq.Force = wrapperspb.Bool(*req.IsForce)
	}
	if req.Latest != nil {
		rpcReq.Latest = wrapperspb.Bool(*req.Latest)
	}
	if req.Hot != nil {
		rpcReq.Hot = wrapperspb.Bool(*req.Hot)
	}
	_, err := l.svcCtx.AdminClient.UpdateApplicationVersion(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.UpdateApplicationVersionResp{}, nil
}

type DeleteApplicationVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteApplicationVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApplicationVersionLogic {
	return &DeleteApplicationVersionLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteApplicationVersionLogic) DeleteApplicationVersion(req *types.DeleteApplicationVersionReq) (*types.DeleteApplicationVersionResp, error) {
	if l.svcCtx.AdminClient == nil {
		return nil, fmt.Errorf("admin service not available")
	}
	_, err := l.svcCtx.AdminClient.DeleteApplicationVersion(l.ctx, &admin.DeleteApplicationVersionReq{Id: []string{req.ID}})
	if err != nil {
		return nil, err
	}
	return &types.DeleteApplicationVersionResp{}, nil
}
