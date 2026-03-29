package logic

import (
	"context"
	"fmt"

	"fdim/Infrastructure_service/push/internal/svc"
	"fdim/protocol/msggateway"
	"fdim/protocol/push"
	"github.com/zeromicro/go-zero/core/logx"
)

type PushMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushMsgLogic {
	return &PushMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushMsgLogic) PushMsg(req *push.PushMsgReq) (*push.PushMsgResp, error) {
	// 1. 验证推送消息的合法?
	if req.MsgData == nil {
		return nil, fmt.Errorf("msgData is empty")
	}
	if req.ConversationID == "" {
		return nil, fmt.Errorf("conversationID is empty")
	}
	if len(req.UserIDs) == 0 {
		return &push.PushMsgResp{}, nil
	}

	// 2. 调用 MessageGateway 进行在线推?
	if l.svcCtx.MessageGatewayClient != nil {
		// 使用 OnlineBatchPushOneMsg 进行批量推?
		_, err := l.svcCtx.MessageGatewayClient.OnlineBatchPushOneMsg(l.ctx, &msggateway.OnlineBatchPushOneMsgReq{
			MsgData:       req.MsgData,
			PushToUserIDs: req.UserIDs,
		})
		if err != nil {
			l.Errorf("Failed to push message online: %v", err)
			return nil, fmt.Errorf("failed to push message: %w", err)
		}
		l.Infof("Pushed message to %d users: conversationID=%s", len(req.UserIDs), req.ConversationID)
	}

	// 3. 返回推送结?
	return &push.PushMsgResp{}, nil
}
