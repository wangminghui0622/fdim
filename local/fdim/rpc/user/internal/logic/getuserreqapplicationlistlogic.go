package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/util"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserReqApplicationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserReqApplicationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserReqApplicationListLogic {
	return &GetUserReqApplicationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserReqApplicationListLogic) GetUserReqApplicationList(req *user.GetUserReqApplicationListReq) (*user.GetUserReqApplicationListResp, error) {
	resp := &user.GetUserReqApplicationListResp{}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// ת�� handleResults
	var handleResults []int32
	if len(req.HandleResults) > 0 {
		handleResults = req.HandleResults
	}

	// �����ҳ����
	offset := util.CalculateOffset(req.Pagination.PageNumber, req.Pagination.ShowNumber)
	limit := util.CalculateLimit(req.Pagination.ShowNumber)

	// ��ҳ��ѯ�û����͵�Ⱥ������
	total, requests, err := l.svcCtx.GroupDB.PageGroupRequest(
		l.ctx,
		req.UserID,
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
