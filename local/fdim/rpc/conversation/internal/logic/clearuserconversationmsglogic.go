package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/errs"
	"fdim/protocol/conversation"
	"fdim/protocol/msg"
	"fdim/rpc/conversation/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ClearUserConversationMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearUserConversationMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearUserConversationMsgLogic {
	return &ClearUserConversationMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ClearUserConversationMsgLogic) ClearUserConversationMsg(req *conversation.ClearUserConversationMsgReq) (*conversation.ClearUserConversationMsgResp, error) {
	resp := &conversation.ClearUserConversationMsgResp{}

	if req.Timestamp <= 0 {
		return nil, errs.ErrArgs.WrapMsg("timestamp must be positive")
	}

	// 获取需要清除消息的会话
	conversations, err := l.svcCtx.ConversationDB.GetConversationIDsNeedDestruct(l.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations need clear msg: %w", err)
	}

	count := 0
	latestMsgDestructTime := time.UnixMilli(req.Timestamp)
	for _, conv := range conversations {
		if !conv.IsMsgDestruct || conv.MsgDestructTime == 0 {
			continue
		}

		// 从 Msg 模块获取该时间之前的最后一条消息 seq
		var minSeq int64
		if l.svcCtx.MsgClient != nil {
			seqResp, err := l.svcCtx.MsgClient.GetLastMessageSeqByTime(l.ctx, &msg.GetLastMessageSeqByTimeReq{
				ConversationID: conv.ConversationID,
				Time:           req.Timestamp,
			})
			if err != nil {
				l.Errorf("failed to get last message seq by time: %v", err)
			} else if seqResp != nil {
				minSeq = seqResp.Seq
			}
		}

		// 更新 latest_msg_destruct_time 和 min_seq
		args := map[string]interface{}{
			"latest_msg_destruct_time": latestMsgDestructTime,
			"min_seq":                  minSeq,
		}
		_, err := l.svcCtx.ConversationDB.UpdateByMap(l.ctx, []string{conv.OwnerUserID}, conv.ConversationID, args)
		if err != nil {
			l.Errorf("failed to update conversation latest_msg_destruct_time: %v", err)
			continue
		}
		count++
	}

	resp.Count = int32(count)
	return resp, nil
}
