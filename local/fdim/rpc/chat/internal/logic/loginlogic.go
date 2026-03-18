package logic

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(req *chat.LoginReq) (*chat.LoginResp, error) {
	if req.Password == "" && req.VerifyCode == "" {
		return nil, errs.ErrArgs.WrapMsg("password or verify code must be set")
	}

	var userAccount *database.UserAccount
	var err error

	switch {
	case req.PhoneNumber != "":
		if req.AreaCode == "" {
			return nil, errs.ErrArgs.WrapMsg("area code is required for phone number")
		}
		if !strings.HasPrefix(req.AreaCode, "+") {
			req.AreaCode = "+" + req.AreaCode
		}
		if _, err := strconv.ParseUint(req.AreaCode[1:], 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("area code must be number")
		}
		userAccount, err = l.svcCtx.ChatDB.GetUserAccountByPhone(l.ctx, req.AreaCode, req.PhoneNumber)
	case req.Email != "":
		userAccount, err = l.svcCtx.ChatDB.GetUserAccountByEmail(l.ctx, req.Email)
	default:
		return nil, errs.ErrArgs.WrapMsg("phone number or email must be set")
	}

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || strings.Contains(err.Error(), "not found") {
			return nil, errs.ErrArgs.WrapMsg("user not found")
		}
		return nil, errs.ErrArgs.WrapMsg("login failed")
	}

	if l.svcCtx.AdminRpc != nil {
		_, err := l.svcCtx.AdminRpc.CheckLoginForbidden(l.ctx, &admin.CheckLoginForbiddenReq{
			UserID: userAccount.UserID,
			Ip:     req.Ip,
		})
		if err != nil {
			l.Errorf("CheckLoginForbidden failed: %v", err)
			return nil, errs.WrapMsg(err, "login forbidden")
		}
	}

	if req.Password != "" {
		hash := sha256.Sum256([]byte(req.Password))
		hashedPassword := hex.EncodeToString(hash[:])
		if userAccount.Password != hashedPassword {
			return nil, errs.ErrArgs.WrapMsg("invalid password")
		}
	} else {
		if req.PhoneNumber != "" {
			verifyCode, err := l.svcCtx.ChatDB.FindVerifyCode(l.ctx, req.PhoneNumber, req.AreaCode, req.VerifyCode)
			if err != nil {
				return nil, errs.ErrArgs.WrapMsg("invalid verify code")
			}
			if time.Now().After(verifyCode.ExpireTime) {
				return nil, errs.ErrArgs.WrapMsg("verify code expired")
			}
			_ = l.svcCtx.ChatDB.DelVerifyCode(l.ctx, req.PhoneNumber, req.AreaCode)
		} else if req.Email != "" {
			verifyCode, err := l.svcCtx.ChatDB.FindVerifyCodeByEmail(l.ctx, req.Email, req.VerifyCode)
			if err != nil {
				return nil, errs.ErrArgs.WrapMsg("invalid verify code")
			}
			if time.Now().After(verifyCode.ExpireTime) {
				return nil, errs.ErrArgs.WrapMsg("verify code expired")
			}
			_ = l.svcCtx.ChatDB.DelVerifyCodeByEmail(l.ctx, req.Email)
		}
	}

	var chatToken string
	if l.svcCtx.AdminRpc != nil {
		tokenResp, err := l.svcCtx.AdminRpc.CreateToken(l.ctx, &admin.CreateTokenReq{
			UserID:   userAccount.UserID,
			UserType: 1,
		})
		if err != nil {
			l.Errorf("CreateToken failed: %v", err)
			return nil, errs.WrapMsg(err, "failed to create token")
		}
		chatToken = tokenResp.Token
	}

	loginRecord := &database.UserLoginRecord{
		UserID:    userAccount.UserID,
		LoginTime: time.Now(),
		IP:        req.Ip,
		DeviceID:  req.DeviceID,
		Platform:  req.Platform,
	}
	if err := l.svcCtx.ChatDB.AddUserLoginRecord(l.ctx, loginRecord); err != nil {
		l.Errorw("AddUserLoginRecord failed", logx.Field("error", err))
	}

	l.Infof("User logged in successfully: userID=%s", userAccount.UserID)
	return &chat.LoginResp{
		UserID:    userAccount.UserID,
		ChatToken: chatToken,
	}, nil
}
