package logic

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterUserLogic {
	return &RegisterUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// genUserID 生成10位数字用户ID
func (l *RegisterUserLogic) genUserID() string {
	const length = 10
	data := make([]byte, length)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		if i == 0 {
			data[i] = chars[1:][data[i]%9] // 第一位不能是0
		} else {
			data[i] = chars[data[i]%10]
		}
	}
	return string(data)
}

// BuildCredentialPhone 构建手机号凭证
func BuildCredentialPhone(areaCode, phone string) string {
	return areaCode + " " + phone
}

func (l *RegisterUserLogic) RegisterUser(req *chat.RegisterUserReq) (*chat.RegisterUserResp, error) {
	// 1. 检查注册是否被禁止
	if l.svcCtx.AdminRpc != nil {
		_, err := l.svcCtx.AdminRpc.CheckRegisterForbidden(l.ctx, &admin.CheckRegisterForbiddenReq{
			Ip: req.Ip,
		})
		if err != nil {
			l.Errorf("CheckRegisterForbidden failed: %v", err)
			return nil, errs.WrapMsg(err, "register forbidden")
		}
	}

	// 2. 验证用户信息
	if req.User == nil {
		return nil, errs.ErrArgs.WrapMsg("user info is required")
	}

	// 检查至少有一种登录方式
	if req.User.Email == "" && req.User.PhoneNumber == "" && req.User.Account == "" {
		return nil, errs.ErrArgs.WrapMsg("at least one account type is required")
	}

	// 3. 验证手机号格式
	if req.User.PhoneNumber != "" {
		if req.User.AreaCode == "" {
			return nil, errs.ErrArgs.WrapMsg("area code is required for phone number")
		}
		if !strings.HasPrefix(req.User.AreaCode, "+") {
			req.User.AreaCode = "+" + req.User.AreaCode
		}
		if _, err := strconv.ParseUint(req.User.AreaCode[1:], 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("area code must be number")
		}
		if _, err := strconv.ParseUint(req.User.PhoneNumber, 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("phone number must be number")
		}
	}

	// 4. 检查用户是否已存在
	if req.User.PhoneNumber != "" {
		_, err := l.svcCtx.ChatDB.GetUserAccountByPhone(l.ctx, req.User.AreaCode, req.User.PhoneNumber)
		if err == nil {
			return nil, errs.ErrArgs.WrapMsg("phone number already registered")
		}
	}
	if req.User.Account != "" {
		_, err := l.svcCtx.ChatDB.GetUserAccountByAccount(l.ctx, req.User.Account)
		if err == nil {
			return nil, errs.ErrArgs.WrapMsg("account already registered")
		}
	}
	if req.User.Email != "" {
		_, err := l.svcCtx.ChatDB.GetUserAccountByEmail(l.ctx, req.User.Email)
		if err == nil {
			return nil, errs.ErrArgs.WrapMsg("email already registered")
		}
	}

	// 5. 验证验证码（非管理员注册需要）
	if req.VerifyCode != "" {
		// 开发阶段：固定测试验证码 123456 直接通过
		if req.VerifyCode == "123456" {
			l.Infof("Using test verify code for registration")
		} else if req.User.PhoneNumber != "" {
			verifyCode, err := l.svcCtx.ChatDB.FindVerifyCode(l.ctx, req.User.PhoneNumber, req.User.AreaCode, req.VerifyCode)
			if err != nil {
				return nil, errs.ErrArgs.WrapMsg("invalid verify code")
			}
			if time.Now().After(verifyCode.ExpireTime) {
				return nil, errs.ErrArgs.WrapMsg("verify code expired")
			}
			_ = l.svcCtx.ChatDB.DelVerifyCode(l.ctx, req.User.PhoneNumber, req.User.AreaCode)
		} else if req.User.Email != "" {
			verifyCode, err := l.svcCtx.ChatDB.FindVerifyCodeByEmail(l.ctx, req.User.Email, req.VerifyCode)
			if err != nil {
				return nil, errs.ErrArgs.WrapMsg("invalid verify code")
			}
			if time.Now().After(verifyCode.ExpireTime) {
				return nil, errs.ErrArgs.WrapMsg("verify code expired")
			}
			_ = l.svcCtx.ChatDB.DelVerifyCodeByEmail(l.ctx, req.User.Email)
		}
	}

	// 6. 检查邀请码（如果需要，这里只验证，不立即使用）
	// 实际使用会在用户创建成功后
	if req.InvitationCode != "" && l.svcCtx.AdminRpc != nil {
		// 这里可以添加邀请码验证逻辑，但不在注册前使用
		// 因为UseInvitationCode会标记为已使用
	}

	// 7. 生成用户ID
	userID := req.User.UserID
	if userID == "" {
		for i := 0; i < 20; i++ {
			userID = l.genUserID()
			_, err := l.svcCtx.ChatDB.GetUserAccountByUserID(l.ctx, userID)
			if err != nil {
				// 用户ID不存在，可以使用
				break
			}
			if i == 19 {
				return nil, errs.ErrInternalServer.WrapMsg("failed to generate user ID")
			}
		}
	} else {
		// 检查指定的用户ID是否已存在
		_, err := l.svcCtx.ChatDB.GetUserAccountByUserID(l.ctx, userID)
		if err == nil {
			return nil, errs.ErrArgs.WrapMsg("user ID already exists")
		}
	}

	// 8. 加密密码
	hashedPassword := req.User.Password
	if hashedPassword != "" {
		hash := sha256.Sum256([]byte(hashedPassword))
		hashedPassword = hex.EncodeToString(hash[:])
	}

	// 9. 创建用户账户
	now := time.Now()
	userAccount := &database.UserAccount{
		UserID:      userID,
		Account:     req.User.Account,
		Password:    hashedPassword,
		AreaCode:    req.User.AreaCode,
		PhoneNumber: req.User.PhoneNumber,
		Email:       req.User.Email,
		CreateTime:  now,
	}
	if err := l.svcCtx.ChatDB.AddUserAccount(l.ctx, []*database.UserAccount{userAccount}); err != nil {
		l.Errorf("AddUserAccount failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to create user account")
	}

	// 10. 创建用户信息
	userInfo := &database.UserFullInfo{
		UserID:           userID,
		Account:          req.User.Account,
		PhoneNumber:      req.User.PhoneNumber,
		AreaCode:         req.User.AreaCode,
		Email:            req.User.Email,
		Nickname:         req.User.Nickname,
		FaceURL:          req.User.FaceURL,
		Gender:           req.User.Gender,
		Level:            1,
		Birth:            req.User.Birth,
		AllowAddFriend:   1,
		AllowBeep:        1,
		AllowVibration:   1,
		GlobalRecvMsgOpt: 0,
		RegisterType:     req.User.RegisterType,
		CreateTime:       now,
	}
	if err := l.svcCtx.ChatDB.AddUserInfo(l.ctx, userInfo); err != nil {
		l.Errorf("AddUserInfo failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to create user info")
	}

	// 11. 使用邀请码（在用户创建成功后使用）
	if req.InvitationCode != "" && l.svcCtx.AdminRpc != nil {
		_, err := l.svcCtx.AdminRpc.UseInvitationCode(l.ctx, &admin.UseInvitationCodeReq{
			Code: req.InvitationCode,
		})
		if err != nil {
			l.Errorw("UseInvitationCode failed", logx.Field("error", err))
		}
	}

	// 12. 与官方一致：注册只注册，不自动签发 token
	// 登录与注册分离，token 只在登录成功后创建
	resp := &chat.RegisterUserResp{
		UserID: userID,
	}

	l.Infof("User registered successfully: userID=%s", userID)
	return resp, nil
}
