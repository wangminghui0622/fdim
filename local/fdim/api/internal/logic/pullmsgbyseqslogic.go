package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/sdkws"
)

type PullMsgBySeqsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPullMsgBySeqsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullMsgBySeqsLogic {
	return &PullMsgBySeqsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PullMsgBySeqsLogic) PullMsgBySeqs(req *types.PullMsgBySeqsReq) (resp *types.PullMsgBySeqsResp, err error) {
	// 转换请求参数
	var seqRanges []*sdkws.SeqRange

	if req.ConversationID != "" && len(req.Seqs) >= 2 {
		minSeq := req.Seqs[0]
		maxSeq := req.Seqs[0]
		for _, s := range req.Seqs[1:] {
			if s < minSeq {
				minSeq = s
			}
			if s > maxSeq {
				maxSeq = s
			}
		}
		seqRanges = append(seqRanges, &sdkws.SeqRange{
			ConversationID: req.ConversationID,
			Begin:          int64(minSeq),
			End:            int64(maxSeq),
		})
	} else if req.ConversationID != "" && len(req.Seqs) == 1 {
		seqRanges = append(seqRanges, &sdkws.SeqRange{
			ConversationID: req.ConversationID,
			Begin:          int64(req.Seqs[0]),
			End:            int64(req.Seqs[0]),
		})
	} else {
		// 兼容旧逻辑：无 conversationID 时逐个 seq 拉取
		for _, seq := range req.Seqs {
			seqRanges = append(seqRanges, &sdkws.SeqRange{
				Begin: int64(seq),
				End:   int64(seq),
			})
		}
	}

	rpcReq := &sdkws.PullMessageBySeqsReq{
		UserID:    req.UserID,
		SeqRanges: seqRanges,
		Order:     0, // 默认顺序
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.MsgClient.PullMessageBySeqs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var msgs []interface{}
	for _, msg := range rpcResp.Msgs {
		msgs = append(msgs, msg)
	}

	resp = &types.PullMsgBySeqsResp{
		Msgs: msgs,
	}

	return resp, nil
}
