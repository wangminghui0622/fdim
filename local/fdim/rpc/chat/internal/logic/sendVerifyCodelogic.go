package logic

import (
	"context"
	"crypto/rand"
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

type SendVerifyCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendVerifyCodeLogic {
	return &SendVerifyCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// genVerifyCode 生成6位数字验证码
func (l *SendVerifyCodeLogic) genVerifyCode() string {
	const length = 6
	data := make([]byte, length)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		data[i] = chars[data[i]%10]
	}
	return string(data)
}

func (l *SendVerifyCodeLogic) SendVerifyCode(req *chat.SendVerifyCodeReq) (*chat.SendVerifyCodeResp, error) {
	// usedFor: 1=注册, 2=登录, 3=重置密码
	const (
		VerificationCodeForRegister      = 1
		VerificationCodeForLogin         = 2
		VerificationCodeForResetPassword = 3
	)

	// 1. 基础参数校验：至少一个联系方式
	if req.Email == "" && (req.AreaCode == "" || req.PhoneNumber == "") {
		return nil, errs.ErrArgs.WrapMsg("email or phone must be set")
	}

	// 2. 根据用途进行不同的验证
	switch req.UsedFor {
	case VerificationCodeForRegister:
		// 检查注册是否被禁止（IP 维度）
		if l.svcCtx.AdminRpc != nil {
			_, err := l.svcCtx.AdminRpc.CheckRegisterForbidden(l.ctx, &admin.CheckRegisterForbiddenReq{
				Ip: req.Ip,
			})
			if err != nil {
				return nil, errs.WrapMsg(err, "register forbidden")
			}
		}

		// 校验手机号或邮箱格式
		if req.Email == "" {
			if req.AreaCode == "" || req.PhoneNumber == "" {
				return nil, errs.ErrArgs.WrapMsg("area code or phone number is empty")
			}
			if !strings.HasPrefix(req.AreaCode, "+") {
				req.AreaCode = "+" + req.AreaCode
			}
			if _, err := strconv.ParseUint(req.AreaCode[1:], 10, 64); err != nil {
				return nil, errs.ErrArgs.WrapMsg("area code must be number")
			}
			if _, err := strconv.ParseUint(req.PhoneNumber, 10, 64); err != nil {
				return nil, errs.ErrArgs.WrapMsg("phone number must be number")
			}
		}
	case VerificationCodeForLogin, VerificationCodeForResetPassword:
		// 登录/重置密码时，账号必须已存在
		if req.Email == "" {
			if req.AreaCode == "" || req.PhoneNumber == "" {
				return nil, errs.ErrArgs.WrapMsg("area code or phone number is empty")
			}
			if _, err := l.svcCtx.ChatDB.GetUserAccountByPhone(l.ctx, req.AreaCode, req.PhoneNumber); err != nil {
				return nil, errs.ErrArgs.WrapMsg("phone unregistered")
			}
		} else {
			if _, err := l.svcCtx.ChatDB.GetUserAccountByEmail(l.ctx, req.Email); err != nil {
				return nil, errs.ErrArgs.WrapMsg("email unregistered")
			}
		}
	default:
		return nil, errs.ErrArgs.WrapMsg("unknown usedFor value")
	}

	// 3. 简单的发送频率限制：同一账号 1 分钟内只允许发送一条验证码
	now := time.Now()
	if req.Email != "" {
		if last, err := l.svcCtx.ChatDB.GetLastVerifyCodeByEmail(l.ctx, req.Email); err == nil && last != nil {
			if now.Sub(last.CreateTime) < time.Minute {
				return nil, errs.ErrArgs.WrapMsg("verify code send too frequently")
			}
		}
	} else {
		if last, err := l.svcCtx.ChatDB.GetLastVerifyCode(l.ctx, req.PhoneNumber, req.AreaCode); err == nil && last != nil {
			if now.Sub(last.CreateTime) < time.Minute {
				return nil, errs.ErrArgs.WrapMsg("verify code send too frequently")
			}
		}
	}

	// 4. 生成验证码
	code := l.genVerifyCode()

	// 5. 保存验证码到数据库
	verifyCode := &database.VerifyCode{
		PhoneNumber: req.PhoneNumber,
		AreaCode:    req.AreaCode,
		Email:       req.Email,
		Code:        code,
		CreateTime:  now,
		ExpireTime:  now.Add(10 * time.Minute), // 10分钟过期
	}
	if err := l.svcCtx.ChatDB.AddVerifyCode(l.ctx, verifyCode); err != nil {
		l.Errorf("AddVerifyCode failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to save verify code")
	}

	// 6. 发送验证码（当前实现：只记录日志，不实际发送短信/邮件）
	if req.Email != "" {
		if l.svcCtx.Mailer != nil {
			if err := l.svcCtx.Mailer.SendMail(l.ctx, req.Email, code); err != nil {
				l.Errorf("SendMail failed: %v", err)
				return nil, errs.WrapMsg(err, "failed to send email verify code")
			}
			l.Infof("Verify code sent to email: %s", req.Email)
		} else {
			l.Infof("Verify code generated (mailer disabled): email=%s, code=%s", req.Email, code)
		}
	} else {
		if l.svcCtx.SMS != nil {
			if err := l.svcCtx.SMS.SendCode(l.ctx, req.AreaCode, req.PhoneNumber, code); err != nil {
				l.Errorf("SendCode failed: %v", err)
				return nil, errs.WrapMsg(err, "failed to send sms verify code")
			}
			l.Infof("Verify code sent to phone: %s %s", req.AreaCode, req.PhoneNumber)
		} else {
			l.Infof("Verify code generated (sms disabled): phone=%s %s, code=%s", req.AreaCode, req.PhoneNumber, code)
		}
	}

	return &chat.SendVerifyCodeResp{}, nil
}
