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

type GetGroupMembersInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMembersInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMembersInfoLogic {
	return &GetGroupMembersInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMembersInfoLogic) GetGroupMembersInfo(req *user.GetGroupMembersInfoReq) (*user.GetGroupMembersInfoResp, error) {
	resp := &user.GetGroupMembersInfoResp{}

	// ???????
	if req.GroupID == "" {
		return nil, fmt.Errorf("groupID is empty")
	}
	if len(req.UserIDs) == 0 {
		return nil, fmt.Errorf("userIDs is empty")
	}

	// ???????????????????
	if !authverify.IsAdmin(l.ctx) {
		opUserID := mcontext.GetOpUserID(l.ctx)
		_, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}
	}

	// ????????
	members, err := l.svcCtx.GroupDB.FindGroupMembers(l.ctx, req.GroupID, req.UserIDs)
	if err != nil {
		return nil, err
	}

	// תΪ Protocol Buffer ʽ
	resp.Members = make([]*sdkws.GroupMemberFullInfo, 0, len(members))
	for _, m := range members {
		pbMember := convert.ModelGroupMemberDB2Pb(m)
		resp.Members = append(resp.Members, pbMember)
	}

	return resp, nil
}
