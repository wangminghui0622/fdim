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

type OpenIMCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOpenIMCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenIMCallbackLogic {
	return &OpenIMCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OpenIMCallbackLogic) OpenIMCallback(req *chat.OpenIMCallbackReq) (*chat.OpenIMCallbackResp, error) {
	l.Infof("OpenIM callback received: command=%s, body=%s", req.Command, req.Body)

	switch req.Command {
	case constantpb.CallbackBeforeAddFriendCommand:
		// 与原 chat 行为保持一致：解析回调体，根据被添加方的配置决定是否允许加好友。
		// 目前 go-im 尚未完整引入 attribute 表，这里采用“统一允许”的策略，并保留扩展点。
		var body map[string]any
		if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
			l.Errorf("failed to unmarshal beforeAddFriend callback body: %v", err)
			return nil, errs.WrapMsg(err, "invalid callback body")
		}
		// 如需细化限制，可在此处根据被添加用户的配置做拦截。
		return &chat.OpenIMCallbackResp{}, nil
	default:
		// 未实现的 command：与 open-im-server/chat 的容错行为对齐，避免回调链路被中断
		l.Errorf("unsupported OpenIM callback command (ignored): %s", req.Command)
		return &chat.OpenIMCallbackResp{}, nil
	}
}
