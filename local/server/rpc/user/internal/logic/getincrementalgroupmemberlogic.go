package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/mcontext"
	"fdim/pkg/model"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetIncrementalGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIncrementalGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalGroupMemberLogic {
	return &GetIncrementalGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetIncrementalGroupMemberLogic) GetIncrementalGroupMember(req *user.GetIncrementalGroupMemberReq) (*user.GetIncrementalGroupMemberResp, error) {
	resp := &user.GetIncrementalGroupMemberResp{}

	// ������֤
	if req.GroupID == "" {
		return nil, fmt.Errorf("groupID is empty")
	}

	// Ȩ�޼�飺��Ҫ��Ⱥ��Ա�����Ա
	if !authverify.IsAdmin(l.ctx) {
		opUserID := mcontext.GetOpUserID(l.ctx)
		_, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, fmt.Errorf("user not in group")
		}
	}

	// TODO: ʵ������Ⱥ��Աͬ���߼�
	// ������Ҫ���ݰ汾�Ż�ȡ������Ⱥ��Ա���
	// ��Ҫ�汾����ģ��֧�֣�
	// 1. ���ݰ汾�Ų�ѯ�汾��־
	// 2. ��ȡ���������¡�ɾ����Ⱥ��Ա�б�
	// 3. ���ذ汾��Ϣ�ͱ���б�

	// ��ʱʵ�֣�����汾��Ϊ0��δ�ṩ�����������б��Full=true��
	// ע�⣺����������ͬ����Ҫ�汾����ģ��֧��
	if req.Version == 0 && req.VersionID == "" {
		// ��ȡ����Ⱥ��Ա
		members, err := l.svcCtx.GroupDB.FindGroupMemberAll(l.ctx, req.GroupID)
		if err != nil {
			return nil, err
		}

		// 转换为 GroupMemberFullInfo
		resp.Insert = convert.ModelGroupMembersDB2Pb(members)

		// ���Ⱥ��Ա���û���Ϣ
		if len(resp.Insert) > 0 {
			// ��ȡ����Ⱥ��Ա���û�ID
			userIDs := make([]string, 0, len(members))
			for _, m := range members {
				userIDs = append(userIDs, m.UserID)
			}

			// ��ȡ�û���Ϣ
			users, err := l.svcCtx.UserDB.Find(l.ctx, userIDs)
			if err == nil {
				// �����û�ӳ��
				userMap := make(map[string]*model.User)
				for _, u := range users {
					userMap[u.UserID] = u
				}

				// ���Ⱥ��Ա���û���Ϣ
				for i, memberInfo := range resp.Insert {
					if i < len(members) {
						member := members[i]
						if user, ok := userMap[member.UserID]; ok {
							memberInfo.Nickname = user.Nickname
							memberInfo.FaceURL = user.FaceURL
						}
					}
				}
			}
		}
		resp.Full = true
	} else {
		// ����ͬ�������ؿ��б����Ҫ�汾����֧�֣�
		// TODO: ʵ�ְ汾�����߼�
		// 1. ���ݰ汾�Ų�ѯ�汾��־
		// 2. ��ȡ���������¡�ɾ����Ⱥ��Ա�б�
		// 3. ���ذ汾��Ϣ�ͱ���б�
		resp.Insert = []*sdkws.GroupMemberFullInfo{}
		resp.Update = []*sdkws.GroupMemberFullInfo{}
		resp.Delete = []string{}
		resp.Full = false
	}

	// TODO: �Ӱ汾����ģ���ȡ�汾��Ϣ
	resp.Version = 0
	resp.VersionID = ""

	return resp, nil
}
