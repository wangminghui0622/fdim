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

type GetIncrementalJoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIncrementalJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalJoinGroupLogic {
	return &GetIncrementalJoinGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetIncrementalJoinGroupLogic) GetIncrementalJoinGroup(req *user.GetIncrementalJoinGroupReq) (*user.GetIncrementalJoinGroupResp, error) {
	resp := &user.GetIncrementalJoinGroupResp{}

	// Ȩ����֤
	opUserID := mcontext.GetOpUserID(l.ctx)
	if req.UserID == "" {
		req.UserID = opUserID
	} else {
		if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
			return nil, err
		}
	}

	// ������֤
	if req.UserID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	// TODO: ʵ�����������Ⱥ��ͬ���߼�
	// ������Ҫ���ݰ汾�Ż�ȡ������Ⱥ����
	// ��Ҫ�汾����ģ��֧�֣�
	// 1. ���ݰ汾�Ų�ѯ�汾��־
	// 2. ��ȡ���������¡�ɾ����Ⱥ���б�
	// 3. ���ذ汾��Ϣ�ͱ���б�

	// ��ʱʵ�֣�����汾��Ϊ0��δ�ṩ�����������б��Full=true��
	// ע�⣺����������ͬ����Ҫ�汾����ģ��֧��
	if req.Version == 0 && req.VersionID == "" {
		// ��ȡ�û����������Ⱥ��
		members, err := l.svcCtx.GroupDB.FindGroupMemberByUserID(l.ctx, req.UserID)
		if err != nil {
			return nil, err
		}

		// ��ȡȺ��ID�б�
		groupIDs := make([]string, 0, len(members))
		for _, m := range members {
			groupIDs = append(groupIDs, m.GroupID)
		}

		// ��ȡȺ����Ϣ
		if len(groupIDs) > 0 {
			groups, err := l.svcCtx.GroupDB.FindGroup(l.ctx, groupIDs)
			if err != nil {
				return nil, err
			}

			// ��ȡȺ����Ϣ
			owners, err := l.svcCtx.GroupDB.FindGroupsOwner(l.ctx, groupIDs)
			if err != nil {
				return nil, err
			}

			// ��ȡ��Ա����
			memberNums, err := l.svcCtx.GroupDB.MapGroupMemberNum(l.ctx, groupIDs)
			if err != nil {
				return nil, err
			}

			// ����Ⱥ��ӳ��
			ownerMap := make(map[string]string)
			for _, owner := range owners {
				ownerMap[owner.GroupID] = owner.UserID
			}

			// 转换为 GroupInfo
			groupInfos := make([]*sdkws.GroupInfo, 0, len(groups))
			for _, g := range groups {
				pbGroup := convert.ModelGroupDB2Pb(g)
				if pbGroup != nil {
					pbGroup.OwnerUserID = ownerMap[g.GroupID]
					pbGroup.MemberCount = memberNums[g.GroupID]
				}
				groupInfos = append(groupInfos, pbGroup)
			}

			resp.Insert = groupInfos
		} else {
			resp.Insert = []*sdkws.GroupInfo{}
		}
		resp.Full = true
	} else {
		// ����ͬ�������ؿ��б����Ҫ�汾����֧�֣�
		// TODO: ʵ�ְ汾�����߼�
		// 1. ���ݰ汾�Ų�ѯ�汾��־
		// 2. ��ȡ���������¡�ɾ����Ⱥ���б�
		// 3. ���ذ汾��Ϣ�ͱ���б�
		resp.Insert = []*sdkws.GroupInfo{}
		resp.Update = []*sdkws.GroupInfo{}
		resp.Delete = []string{}
		resp.Full = false
	}

	// TODO: �Ӱ汾����ģ���ȡ�汾��Ϣ
	resp.Version = 0
	resp.VersionID = ""

	return resp, nil
}
