package logic

import (
	"context"
	"crypto/rand"
	"strconv"
	"strings"
	"time"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserAccountLogic {
	return &AddUserAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// genUserID 生成10位数字用户ID
func (l *AddUserAccountLogic) genUserID() string {
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

func (l *AddUserAccountLogic) AddUserAccount(req *chat.AddUserAccountReq) (*chat.AddUserAccountResp, error) {
	// 1. 验证用户信息
	if req.User == nil {
		return nil, errs.ErrArgs.WrapMsg("user info is required")
	}

	// 检查至少有一种登录方式
	if req.User.Email == "" && req.User.PhoneNumber == "" && req.User.Account == "" {
		return nil, errs.ErrArgs.WrapMsg("at least one account type is required")
	}

	// 2. 验证手机号格式
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

	// 3. 检查用户是否已存在
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

	// 4. 生成用户ID
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

	// 5. 创建用户账户
	now := time.Now()
	userAccount := &database.UserAccount{
		UserID:      userID,
		Account:     req.User.Account,
		Password:    req.User.Password,
		AreaCode:    req.User.AreaCode,
		PhoneNumber: req.User.PhoneNumber,
		Email:       req.User.Email,
		CreateTime:  now,
	}
	if err := l.svcCtx.ChatDB.AddUserAccount(l.ctx, []*database.UserAccount{userAccount}); err != nil {
		l.Errorf("AddUserAccount failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to create user account")
	}

	// 6. 创建用户信息
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

	l.Infof("User account added successfully: userID=%s", userID)
	return &chat.AddUserAccountResp{}, nil
}
