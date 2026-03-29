package logic

import (
	"context"
	"encoding/json"

	"fdim/pkg/errs"
	"fdim/protocol/chat"
	constantpb "fdim/protocol/constant"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FDIMCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFDIMCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FDIMCallbackLogic {
	return &FDIMCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FDIMCallbackLogic) FDIMCallback(req *chat.FDIMCallbackReq) (*chat.FDIMCallbackResp, error) {
	l.Infof("FDIM callback received: command=%s, body=%s", req.Command, req.Body)

	switch req.Command {
	case constantpb.CallbackBeforeAddFriendCommand:
		var body map[string]any
		if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
			l.Errorf("failed to unmarshal beforeAddFriend callback body: %v", err)
			return nil, errs.WrapMsg(err, "invalid callback body")
		}
		return &chat.FDIMCallbackResp{}, nil
	default:
		l.Errorf("unsupported FDIM callback command (ignored): %s", req.Command)
		return &chat.FDIMCallbackResp{}, nil
	}
}
