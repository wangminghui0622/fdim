package logic

import (
	"context"
	"fdim/Infrastructure_service/msggateway/internal/svc"

	"fdim/protocol/msggateway"
	"github.com/zeromicro/go-zero/core/logx"
)

type SuperGroupOnlineBatchPushOneMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSuperGroupOnlineBatchPushOneMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SuperGroupOnlineBatchPushOneMsgLogic {
	return &SuperGroupOnlineBatchPushOneMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SuperGroupOnlineBatchPushOneMsgLogic) SuperGroupOnlineBatchPushOneMsg(req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	// 超级群批量推送与普通批量推送逻辑相同
	// 可以在这里添加超级群特定的逻辑，比如权限检查等
	logic := NewOnlineBatchPushOneMsgLogic(l.ctx, l.svcCtx)
	return logic.OnlineBatchPushOneMsg(req)
}
