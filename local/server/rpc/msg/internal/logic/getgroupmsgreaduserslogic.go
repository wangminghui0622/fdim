package logic

import (
	"context"

	"fdim/protocol/msg"
	"fdim/rpc/msg/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMsgReadUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMsgReadUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMsgReadUsersLogic {
	return &GetGroupMsgReadUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupMsgReadUsers 获取群消息已读/未读成员列表
func (l *GetGroupMsgReadUsersLogic) GetGroupMsgReadUsers(req *msg.GetGroupMsgReadUsersReq) (*msg.GetGroupMsgReadUsersResp, error) {
	if req.ConversationID == "" || len(req.MemberUserIDs) == 0 {
		return &msg.GetGroupMsgReadUsersResp{}, nil
	}

	readUserIDs, unreadUserIDs, err := l.svcCtx.MsgCache.GetGroupMsgReadUsers(
		l.ctx,
		req.ConversationID,
		req.Seq,
		req.MemberUserIDs,
	)
	if err != nil {
		l.Errorf("GetGroupMsgReadUsers failed: %v", err)
		return nil, err
	}

	return &msg.GetGroupMsgReadUsersResp{
		ReadUserIDs:   readUserIDs,
		UnreadUserIDs: unreadUserIDs,
	}, nil
}
