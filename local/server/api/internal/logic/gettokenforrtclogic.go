package logic

import (
	"context"
	"errors"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenForRTCLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTokenForRTCLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenForRTCLogic {
	return &GetTokenForRTCLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenForRTCLogic) GetTokenForRTC(req *types.GetTokenForRTCReq) (*types.GetTokenForRTCResp, error) {
	if l.svcCtx.LiveKit == nil {
		return nil, errors.New("LiveKit is not configured")
	}
	if req.Room == "" || req.Identity == "" {
		return nil, errors.New("room and identity are required")
	}

	token, err := l.svcCtx.LiveKit.GetToken(req.Room, req.Identity)
	if err != nil {
		return nil, err
	}

	return &types.GetTokenForRTCResp{
		ServerUrl: l.svcCtx.LiveKit.GetURL(),
		Token:     token,
	}, nil
}
