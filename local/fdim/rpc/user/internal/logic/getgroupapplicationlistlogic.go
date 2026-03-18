package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/mcontext"
	"fdim/pkg/util"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupApplicationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupApplicationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupApplicationListLogic {
	return &GetGroupApplicationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupApplicationListLogic) GetGroupApplicationList(req *user.GetGroupApplicationListReq) (*user.GetGroupApplicationListResp, error) {
	resp := &user.GetGroupApplicationListResp{}

	// ������֤
	if len(req.GroupIDs) == 0 {
		return nil, fmt.Errorf("groupIDs is empty")
	}

	// Ȩ����֤������Ƿ��ǹ���Ա��Ⱥ��/����Ա
	opUserID := mcontext.GetOpUserID(l.ctx)
	if !authverify.IsAdmin(l.ctx) {
		// ���ÿ��Ⱥ���в����û��Ƿ���Ⱥ�������Ա
		for _, groupID := range req.GroupIDs {
			isAdmin, err := l.svcCtx.GroupDB.IsGroupAdmin(l.ctx, groupID, opUserID)
			if err != nil {
				return nil, err
			}
			if !isAdmin {
				return nil, fmt.Errorf("no permission: user %s is not admin or group owner/admin of group %s", opUserID, groupID)
			}
		}
	}

	// ת�� handleResults
	var handleResults []int32
	if len(req.HandleResults) > 0 {
		handleResults = req.HandleResults
	}

	// �����ҳ����
	offset := util.CalculateOffset(req.Pagination.PageNumber, req.Pagination.ShowNumber)
	limit := util.CalculateLimit(req.Pagination.ShowNumber)

	// ��ҳ��ѯȺ������
	total, requests, err := l.svcCtx.GroupDB.PageGroupRequestByGroup(
		l.ctx,
		req.GroupIDs,
		handleResults,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}

	// ת��Ϊ Protocol Buffer ��ʽ
	resp.Total = uint32(total)
	// ����������������ȡ�û���Ⱥ����Ϣ
	getUserInfo := func(userID string) *sdkws.PublicUserInfo {
		user, err := l.svcCtx.UserDB.Take(l.ctx, userID)
		if err != nil || user == nil {
			return nil
		}
		return &sdkws.PublicUserInfo{
			UserID:   user.UserID,
			Nickname: user.Nickname,
			FaceURL:  user.FaceURL,
			Ex:       user.Ex,
		}
	}
	getGroupInfo := func(groupID string) *sdkws.GroupInfo {
		group, err := l.svcCtx.GroupDB.TakeGroup(l.ctx, groupID)
		if err != nil || group == nil {
			return nil
		}
		owner, _ := l.svcCtx.GroupDB.TakeGroupOwner(l.ctx, groupID)
		memberNumMap, _ := l.svcCtx.GroupDB.MapGroupMemberNum(l.ctx, []string{groupID})
		memberNum := memberNumMap[groupID]
		pbGroup := convert.ModelGroupDB2Pb(group)
		if pbGroup != nil {
			pbGroup.OwnerUserID = owner.UserID
			pbGroup.MemberCount = memberNum
		}
		return pbGroup
	}
	resp.GroupRequests = convert.GroupRequestsDB2Pb(requests, getUserInfo, getGroupInfo)

	return resp, nil
}
