package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenForVideoMeetingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenForVideoMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenForVideoMeetingLogic {
	return &GetTokenForVideoMeetingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTokenForVideoMeetingLogic) GetTokenForVideoMeeting(req *chat.GetTokenForVideoMeetingReq) (*chat.GetTokenForVideoMeetingResp, error) {
	// 1. 验证参数
	if req.Room == "" {
		return nil, errs.ErrArgs.WrapMsg("room cannot be empty")
	}
	if req.Identity == "" {
		return nil, errs.ErrArgs.WrapMsg("identity cannot be empty")
	}

	// 2. 简化实现：基于配置生成签名 Token（避免纯占位符）
	serverURL := l.svcCtx.Config.RTC.ServerURL
	if serverURL == "" {
		serverURL = "https://video-meeting.example.com"
	}
	salt := l.svcCtx.Config.RTC.TokenSalt
	if salt == "" {
		salt = l.svcCtx.Config.Secret
	}
	issuedAt := time.Now().Unix()
	payload := fmt.Sprintf("room=%s&identity=%s&iat=%d", req.Room, req.Identity, issuedAt)
	mac := hmac.New(sha256.New, []byte(salt))
	_, _ = mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))
	token := payload + "&sig=" + sig

	l.Infof("Generated video meeting token: room=%s, identity=%s", req.Room, req.Identity)
	return &chat.GetTokenForVideoMeetingResp{
		ServerUrl: serverURL,
		Token:     token,
	}, nil
}
