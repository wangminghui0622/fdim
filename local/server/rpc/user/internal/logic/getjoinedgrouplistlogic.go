package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/mcontext"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetJoinedGroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetJoinedGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJoinedGroupListLogic {
	return &GetJoinedGroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetJoinedGroupListLogic) GetJoinedGroupList(req *user.GetJoinedGroupListReq) (*user.GetJoinedGroupListResp, error) {
	resp := &user.GetJoinedGroupListResp{}

	// ??????
	opUserID := mcontext.GetOpUserID(l.ctx)
	userID := req.FromUserID
	if userID == "" {
		userID = opUserID
	} else {
		if err := authverify.CheckAccess(l.ctx, userID); err != nil {
			return nil, err
		}
	}

	// ???????????????????????
	members, err := l.svcCtx.GroupDB.FindGroupMemberByUserID(l.ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		resp.Total = 0
		resp.Groups = []*sdkws.GroupInfo{}
		return resp, nil
	}

	// ??????ID??
	groupIDs := make([]string, len(members))
	for i, m := range members {
		groupIDs[i] = m.GroupID
	}

	// ?????????
	groups, err := l.svcCtx.GroupDB.FindGroup(l.ctx, groupIDs)
	if err != nil {
		return nil, err
	}

	// ????????????
	memberNums, err := l.svcCtx.GroupDB.MapGroupMemberNum(l.ctx, groupIDs)
	if err != nil {
		return nil, err
	}

	// ?????????
	owners, err := l.svcCtx.GroupDB.FindGroupsOwner(l.ctx, groupIDs)
	if err != nil {
		return nil, err
	}

	// ??????????
	ownerMap := make(map[string]string)
	for _, owner := range owners {
		ownerMap[owner.GroupID] = owner.UserID
	}

	// תΪ Protocol Buffer ʽ
	resp.Groups = make([]*sdkws.GroupInfo, 0, len(groups))
	for _, g := range groups {
		pbGroup := convert.ModelGroupDB2Pb(g)
		if pbGroup != nil {
			pbGroup.MemberCount = memberNums[g.GroupID]
			pbGroup.OwnerUserID = ownerMap[g.GroupID]
		}
		resp.Groups = append(resp.Groups, pbGroup)
	}

	resp.Total = uint32(len(resp.Groups))

	return resp, nil
}
