package logic

import (
	"context"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/chat"

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

func (l *SendVerifyCodeLogic) SendVerifyCode(req *types.SendVerifyCodeReq) (*types.SendVerifyCodeResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, errs.ErrInternalServer.WrapMsg("chat service not available")
	}

	_, err := l.svcCtx.ChatClient.SendVerifyCode(l.ctx, &chat.SendVerifyCodeReq{
		UsedFor:     req.UsedFor,
		PhoneNumber: req.PhoneNumber,
		AreaCode:    req.AreaCode,
		Email:       req.Email,
		Ip:          req.Ip,
		Platform:    req.Platform,
		DeviceID:    req.DeviceID,
	})
	if err != nil {
		l.Errorf("SendVerifyCode failed: %v", err)
		return nil, err
	}

	return &types.SendVerifyCodeResp{}, nil
}
