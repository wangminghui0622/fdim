package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/pkg/model"
	"fdim/pkg/msgprocessor"
	"fdim/protocol/sdkws"
	"fdim/rpc/msg/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PullMessageBySeqsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPullMessageBySeqsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullMessageBySeqsLogic {
	return &PullMessageBySeqsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PullMessageBySeqsLogic) PullMessageBySeqs(req *sdkws.PullMessageBySeqsReq) (*sdkws.PullMessageBySeqsResp, error) {
	if l.svcCtx.MsgCache == nil {
		return nil, errs.ErrInternalServer.WrapMsg("message cache not initialized")
	}

	resp := &sdkws.PullMessageBySeqsResp{
		Msgs:             make(map[string]*sdkws.PullMsgs),
		NotificationMsgs: make(map[string]*sdkws.PullMsgs),
	}

	for _, rng := range req.SeqRanges {
		if rng.ConversationID == "" {
			continue
		}
		// 构造需要拉取的 seq 列表
		var seqs []int64
		begin := rng.Begin
		end := rng.End
		if end < begin {
			begin, end = end, begin
		}
		var count int64
		for s := begin; s <= end; s++ {
			seqs = append(seqs, s)
			count++
			if rng.Num > 0 && count >= rng.Num {
				break
			}
		}
		if len(seqs) == 0 {
			continue
		}

		// 从缓存中获取消息
		msgDocs, err := l.svcCtx.MsgCache.GetMessagesBySeq(l.ctx, rng.ConversationID, seqs)
		if err != nil {
			l.Errorw("GetMessagesBySeq failed", logx.Field("conversationID", rng.ConversationID), logx.Field("error", err))
			continue
		}

		if len(msgDocs) == 0 {
			continue
		}

		// 转换为 sdkws.MsgData
		pullMsgs := &sdkws.PullMsgs{Msgs: make([]*sdkws.MsgData, 0, len(msgDocs))}
		var lastSeq int64
		for _, doc := range msgDocs {
			md := docToMsgData(doc)
			pullMsgs.Msgs = append(pullMsgs.Msgs, md)
			if md.Seq > lastSeq {
				lastSeq = md.Seq
			}
		}
		pullMsgs.IsEnd = true
		pullMsgs.EndSeq = lastSeq

		// 与官方一致：n_ 前缀的会话放入 NotificationMsgs，其他放入 Msgs
		if msgprocessor.IsNotification(rng.ConversationID) {
			resp.NotificationMsgs[rng.ConversationID] = pullMsgs
		} else {
			resp.Msgs[rng.ConversationID] = pullMsgs
		}
	}

	return resp, nil
}

// docToMsgData 将数据库文档转换为 sdkws.MsgData
func docToMsgData(doc *model.MsgDoc) *sdkws.MsgData {
	if doc == nil {
		return nil
	}
	return &sdkws.MsgData{
		SendID:           doc.SendID,
		RecvID:           doc.RecvID,
		GroupID:          doc.GroupID,
		ClientMsgID:      doc.ClientMsgID,
		ServerMsgID:      doc.ServerMsgID,
		SenderPlatformID: doc.SenderPlatformID,
		SenderNickname:   doc.SenderNickname,
		SenderFaceURL:    doc.SenderFaceURL,
		SessionType:      doc.SessionType,
		MsgFrom:          doc.MsgFrom,
		ContentType:      doc.ContentType,
		Content:          doc.Content,
		Seq:              doc.Seq,
		SendTime:         doc.SendTime,
		CreateTime:       doc.CreateTime,  // 已经是 int64 时间戳
		Status:           doc.Status,
		IsRead:           doc.IsRead,      // 【修复】返回已读状态
		Options:          doc.Options,
		AtUserIDList:     doc.AtUserIDs,
		AttachedInfo:     doc.AttachedInfo,
		Ex:               doc.Ex,
		// OfflinePushInfo 暂不反填，客户端通常不依赖该字段拉历史
	}
}
