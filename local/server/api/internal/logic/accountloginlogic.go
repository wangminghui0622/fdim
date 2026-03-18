package logic

import (
	"context"
	"strings"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/auth"
	"fdim/protocol/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccountLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccountLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountLoginLogic {
	return &AccountLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AccountLoginLogic) AccountLogin(req *types.AccountLoginReq) (*types.AccountLoginResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, errs.ErrInternalServer.WrapMsg("chat service not available")
	}

	if req.PhoneNumber == "" && req.Email == "" && req.Account == "" {
		return nil, errs.ErrArgs.WrapMsg("phone number, email or account required")
	}

	// 手机号登录时默认区号 +86
	if req.PhoneNumber != "" && req.AreaCode == "" {
		req.AreaCode = "+86"
	}

	deviceID := req.DeviceID
	if deviceID == "" {
		deviceID = "unknown"
	}

	ip := req.Ip
	if ip == "" {
		ip = "0.0.0.0"
	}

	chatResp, err := l.svcCtx.ChatClient.Login(l.ctx, &chat.LoginReq{
		Account:     req.Account,
		AreaCode:    req.AreaCode,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    req.Password,
		Platform:    req.Platform,
		DeviceID:    deviceID,
		Ip:          ip,
	})
	if err != nil {
		l.Errorf("chat.Login failed: %v", err)
		errMsg := err.Error()
		if strings.Contains(errMsg, "用户不存在") || strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "not exist") {
			return nil, errs.ErrRecordNotFound.WrapMsg("账号不存在")
		}
		if strings.Contains(errMsg, "密码错误") || strings.Contains(errMsg, "password") {
			return nil, errs.ErrArgs.WrapMsg("密码错误")
		}
		return nil, errs.ErrInternalServer.WrapMsg("登录失败，请稍后重试")
	}

	tokenResp, err := l.svcCtx.AuthClient.GetUserToken(l.ctx, &auth.GetUserTokenReq{
		UserID:     chatResp.UserID,
		PlatformID: req.Platform,
	})
	if err != nil {
		l.Errorf("GetUserToken failed: userID=%s, err=%v", chatResp.UserID, err)
		return nil, errs.ErrInternalServer.WrapMsg("failed to get token")
	}

	return &types.AccountLoginResp{
		UserID:    chatResp.UserID,
		ImToken:   tokenResp.Token,
		ChatToken: chatResp.ChatToken,
	}, nil
}
