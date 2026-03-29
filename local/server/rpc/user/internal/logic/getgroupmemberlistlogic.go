package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/mcontext"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMemberListLogic) GetGroupMemberList(req *user.GetGroupMemberListReq) (*user.GetGroupMemberListResp, error) {
	resp := &user.GetGroupMemberListResp{}

	// ???????????????????
	if !authverify.IsAdmin(l.ctx) {
		opUserID := mcontext.GetOpUserID(l.ctx)
		_, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}
	}

	// ???????????
	members, err := l.svcCtx.GroupDB.FindGroupMemberAll(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// ռгԱ userIDѯûϢ
	userIDs := make([]string, 0, len(members))
	for _, m := range members {
		userIDs = append(userIDs, m.UserID)
	}
	userMap := make(map[string]*sdkws.UserInfo)
	if users, err := l.svcCtx.UserDB.Find(l.ctx, userIDs); err == nil {
		for _, u := range users {
			userMap[u.UserID] = convert.ModelUserDB2Pb(u)
		}
	}

	// תΪ Protocol Buffer ʽϲû nickname/faceURL
	resp.Members = make([]*sdkws.GroupMemberFullInfo, 0, len(members))
	for _, m := range members {
		pbMember := convert.ModelGroupMemberDB2Pb(m)
		if u, ok := userMap[m.UserID]; ok {
			// ûʱȺԱе nickname/faceURL
			if u.Nickname != "" {
				pbMember.Nickname = u.Nickname
			}
			if u.FaceURL != "" {
				pbMember.FaceURL = u.FaceURL
			}
		}
		resp.Members = append(resp.Members, pbMember)
	}

	resp.Total = uint32(len(members))
	return resp, nil
}
