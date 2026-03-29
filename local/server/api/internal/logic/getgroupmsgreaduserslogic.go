package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
	"fdim/protocol/user"
)

type GetGroupMsgReadUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupMsgReadUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMsgReadUsersLogic {
	return &GetGroupMsgReadUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupMsgReadUsersLogic) GetGroupMsgReadUsers(req *types.GetGroupMsgReadUsersReq) (*types.GetGroupMsgReadUsersResp, error) {
	resp := &types.GetGroupMsgReadUsersResp{
		ReadUsers:   []types.GroupMsgReadUser{},
		UnreadUsers: []types.GroupMsgReadUser{},
	}

	// 1. 获取群成员列表
	membersResp, err := l.svcCtx.GroupClient.GetGroupMemberList(l.ctx, &user.GetGroupMemberListReq{
		GroupID: req.GroupID,
		Pagination: &sdkws.RequestPagination{
			PageNumber: 1,
			ShowNumber: 1000,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(membersResp.Members) == 0 {
		return resp, nil
	}

	// 2. 构建成员映射和 userID 列表
	memberMap := make(map[string]*types.GroupMsgReadUser)
	memberUserIDs := make([]string, len(membersResp.Members))
	for i, m := range membersResp.Members {
		memberUserIDs[i] = m.UserID
		memberMap[m.UserID] = &types.GroupMsgReadUser{
			UserID:   m.UserID,
			Nickname: m.Nickname,
			FaceURL:  m.FaceURL,
		}
	}

	// 3. 通过 RPC 调用获取已读/未读成员列表
	rpcResp, err := l.svcCtx.MsgClient.GetGroupMsgReadUsers(l.ctx, &msg.GetGroupMsgReadUsersReq{
		ConversationID: req.ConversationID,
		Seq:            req.Seq,
		MemberUserIDs:  memberUserIDs,
	})
	if err != nil {
		return nil, err
	}

	// 4. 构建响应
	for _, userID := range rpcResp.ReadUserIDs {
		if u, ok := memberMap[userID]; ok {
			resp.ReadUsers = append(resp.ReadUsers, *u)
		}
	}
	for _, userID := range rpcResp.UnreadUserIDs {
		if u, ok := memberMap[userID]; ok {
			resp.UnreadUsers = append(resp.UnreadUsers, *u)
		}
	}

	return resp, nil
}
