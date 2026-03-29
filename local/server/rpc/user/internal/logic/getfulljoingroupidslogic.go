package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/mcontext"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFullJoinGroupIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFullJoinGroupIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFullJoinGroupIDsLogic {
	return &GetFullJoinGroupIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFullJoinGroupIDsLogic) GetFullJoinGroupIDs(req *user.GetFullJoinGroupIDsReq) (*user.GetFullJoinGroupIDsResp, error) {
	resp := &user.GetFullJoinGroupIDsResp{}

	// ??????
	opUserID := mcontext.GetOpUserID(l.ctx)
	if req.UserID == "" {
		req.UserID = opUserID
	} else {
		if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
			return nil, err
		}
	}

	// ???????????????????????
	members, err := l.svcCtx.GroupDB.FindGroupMemberByUserID(l.ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// ??????ID??
	groupIDs := make([]string, 0, len(members))
	for _, m := range members {
		groupIDs = append(groupIDs, m.GroupID)
	}

	resp.GroupIDs = groupIDs

	return resp, nil
}
