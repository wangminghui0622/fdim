package logic

import (
	"context"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/auth"
	"fdim/protocol/chat"
	"fdim/protocol/sdkws"
	"fdim/protocol/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccountRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccountRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountRegisterLogic {
	return &AccountRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AccountRegisterLogic) AccountRegister(req *types.AccountRegisterReq) (*types.AccountRegisterResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, errs.ErrInternalServer.WrapMsg("chat service not available")
	}

	if req.User == nil {
		return nil, errs.ErrArgs.WrapMsg("field \"user\" is required")
	}

	u := req.User
	if u.PhoneNumber == "" && u.Email == "" && u.Account == "" {
		return nil, errs.ErrArgs.WrapMsg("phone number, email or account required")
	}
	if u.Nickname == "" {
		return nil, errs.ErrArgs.WrapMsg("nickname is required")
	}
	if len(u.Password) < 6 {
		return nil, errs.ErrArgs.WrapMsg("password too short")
	}

	// 手机号注册时默认区号 +86
	if u.PhoneNumber != "" && u.AreaCode == "" {
		u.AreaCode = "+86"
	}

	deviceID := req.DeviceID
	if deviceID == "" {
		deviceID = "unknown"
	}

	chatResp, err := l.svcCtx.ChatClient.RegisterUser(l.ctx, &chat.RegisterUserReq{
		VerifyCode: req.VerifyCode,
		Platform:   req.Platform,
		DeviceID:   deviceID,
		AutoLogin:  req.AutoLogin,
		User: &chat.RegisterUserInfo{
			Nickname:    u.Nickname,
			FaceURL:     u.FaceURL,
			PhoneNumber: u.PhoneNumber,
			AreaCode:    u.AreaCode,
			Email:       u.Email,
			Account:     u.Account,
			Password:    u.Password,
		},
	})
	if err != nil {
		l.Errorf("chat.RegisterUser failed: %v", err)
		return nil, err
	}

	_, err = l.svcCtx.UserClient.UserRegister(l.ctx, &user.UserRegisterReq{
		Users: []*sdkws.UserInfo{
			{
				UserID:   chatResp.UserID,
				Nickname: u.Nickname,
				FaceURL:  u.FaceURL,
			},
		},
	})
	if err != nil {
		l.Errorf("user.UserRegister failed: userID=%s, err=%v", chatResp.UserID, err)
	}

	tokenResp, err := l.svcCtx.AuthClient.GetUserToken(l.ctx, &auth.GetUserTokenReq{
		UserID:     chatResp.UserID,
		PlatformID: req.Platform,
	})
	if err != nil {
		l.Errorf("GetUserToken failed: userID=%s, err=%v", chatResp.UserID, err)
		return &types.AccountRegisterResp{
			UserID: chatResp.UserID,
		}, nil
	}

	return &types.AccountRegisterResp{
		UserID:    chatResp.UserID,
		ImToken:   tokenResp.Token,
		ChatToken: chatResp.ChatToken,
	}, nil
}
