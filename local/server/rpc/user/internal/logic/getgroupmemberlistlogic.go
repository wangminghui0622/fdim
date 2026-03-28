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

	// Ȩ�޼�飺��Ҫ��Ⱥ��Ա�����Ա
	if !authverify.IsAdmin(l.ctx) {
		opUserID := mcontext.GetOpUserID(l.ctx)
		_, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}
	}

	// ��ȡ����Ⱥ��Ա
	members, err := l.svcCtx.GroupDB.FindGroupMemberAll(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// 收集所有成员的 userID，批量查询用户最新信息
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

	// 转换为 Protocol Buffer 格式，合并用户表最新 nickname/faceURL
	resp.Members = make([]*sdkws.GroupMemberFullInfo, 0, len(members))
	for _, m := range members {
		pbMember := convert.ModelGroupMemberDB2Pb(m)
		if u, ok := userMap[m.UserID]; ok {
			// 用户表有最新数据时，覆盖群成员快照中的 nickname/faceURL
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
