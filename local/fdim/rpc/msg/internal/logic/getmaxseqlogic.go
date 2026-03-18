package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/pkg/util/conversationutil"
	"fdim/protocol/conversation"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMaxSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMaxSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMaxSeqLogic {
	return &GetMaxSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMaxSeqLogic) GetMaxSeq(req *sdkws.GetMaxSeqReq) (*sdkws.GetMaxSeqResp, error) {
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	resp := &sdkws.GetMaxSeqResp{
		MaxSeqs: make(map[string]int64),
		MinSeqs: make(map[string]int64),
	}

	// 1. 获取用户所有会话 ID
	var conversationIDs []string
	if l.svcCtx.ConversationClient != nil {
		convResp, err := l.svcCtx.ConversationClient.GetConversationIDs(l.ctx, &conversation.GetConversationIDsReq{
			UserID: req.UserID,
		})
		if err != nil {
			l.Errorw("GetConversationIDs failed", logx.Field("error", err))
		} else {
			conversationIDs = convResp.ConversationIDs
		}
	}

	// 2. 为每个会话获取 maxSeq 和 minSeq（包括对应的通知会话）
	for _, convID := range conversationIDs {
		// 普通会话 maxSeq
		maxSeq, err := l.svcCtx.MsgCache.GetMaxSeq(l.ctx, convID)
		if err != nil {
			l.Errorw("GetMaxSeq failed", logx.Field("conversationID", convID), logx.Field("error", err))
			continue
		}
		if maxSeq > 0 {
			resp.MaxSeqs[convID] = maxSeq
		}

		// 普通会话 minSeq
		minSeq, err := l.svcCtx.MsgCache.GetMinSeq(l.ctx, convID)
		if err != nil {
			l.Errorw("GetMinSeq failed", logx.Field("conversationID", convID), logx.Field("error", err))
		}
		if minSeq > 0 {
			resp.MinSeqs[convID] = minSeq
		}

		// 通知会话 maxSeq（n_ 前缀）
		notifConvID := conversationutil.GetNotificationConversationIDByConversationID(convID)
		if notifConvID != "" {
			nMaxSeq, err := l.svcCtx.MsgCache.GetMaxSeq(l.ctx, notifConvID)
			if err != nil {
				continue
			}
			if nMaxSeq > 0 {
				resp.MaxSeqs[notifConvID] = nMaxSeq
			}
			nMinSeq, err := l.svcCtx.MsgCache.GetMinSeq(l.ctx, notifConvID)
			if err != nil {
				continue
			}
			if nMinSeq > 0 {
				resp.MinSeqs[notifConvID] = nMinSeq
			}
		}
	}

	return resp, nil
}
