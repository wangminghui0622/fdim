package logic

import (
	"context"

	"fdim/pkg/convert"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupsInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupsInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupsInfoLogic {
	return &GetGroupsInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupsInfoLogic) GetGroupsInfo(req *user.GetGroupsInfoReq) (*user.GetGroupsInfoResp, error) {
	resp := &user.GetGroupsInfoResp{}

	// ???????
	groups, err := l.svcCtx.GroupDB.FindGroup(l.ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}

	// ????????????
	memberNums, err := l.svcCtx.GroupDB.MapGroupMemberNum(l.ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}

	// ?????????
	groupIDs := make([]string, len(groups))
	for i, g := range groups {
		groupIDs[i] = g.GroupID
	}
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
	resp.GroupInfos = make([]*sdkws.GroupInfo, 0, len(groups))
	for _, g := range groups {
		pbGroup := convert.ModelGroupDB2Pb(g)
		if pbGroup != nil {
			pbGroup.MemberCount = memberNums[g.GroupID]
			pbGroup.OwnerUserID = ownerMap[g.GroupID]
		}
		resp.GroupInfos = append(resp.GroupInfos, pbGroup)
	}

	return resp, nil
}
