package logic

import (
	"context"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserFullInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserFullInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserFullInfoLogic {
	return &FindUserFullInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserFullInfoLogic) FindUserFullInfo(req *types.FindUserFullInfoReq) (*types.FindUserFullInfoResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, errs.ErrInternalServer.WrapMsg("chat service not available")
	}

	resp, err := l.svcCtx.ChatClient.FindUserFullInfo(l.ctx, &chat.FindUserFullInfoReq{
		UserIDs: req.UserIDs,
	})
	if err != nil {
		l.Errorf("FindUserFullInfo failed: %v", err)
		return nil, err
	}

	return &types.FindUserFullInfoResp{Users: resp.Users}, nil
}
