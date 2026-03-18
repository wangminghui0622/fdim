package logic

import (
	"context"
	"time"

	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetServerTimeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetServerTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerTimeLogic {
	return &GetServerTimeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetServerTimeLogic) GetServerTime(req *msg.GetServerTimeReq) (*msg.GetServerTimeResp, error) {
	// 返回当前服务器时间
	return &msg.GetServerTimeResp{
		ServerTime: time.Now().Unix(),
	}, nil
}
