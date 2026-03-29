package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/mcontext"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFullGroupMemberUserIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFullGroupMemberUserIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFullGroupMemberUserIDsLogic {
	return &GetFullGroupMemberUserIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFullGroupMemberUserIDsLogic) GetFullGroupMemberUserIDs(req *user.GetFullGroupMemberUserIDsReq) (*user.GetFullGroupMemberUserIDsResp, error) {
	resp := &user.GetFullGroupMemberUserIDsResp{}

	// ???????????????????
	if !authverify.IsAdmin(l.ctx) {
		opUserID := mcontext.GetOpUserID(l.ctx)
		_, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}
	}

	// ???????
	if req.GroupID == "" {
		return nil, fmt.Errorf("groupID is empty")
	}

	// ???????????
	members, err := l.svcCtx.GroupDB.FindGroupMemberAll(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// ??????ID??
	userIDs := make([]string, 0, len(members))
	for _, m := range members {
		userIDs = append(userIDs, m.UserID)
	}

	resp.UserIDs = userIDs

	return resp, nil
}
