package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/model"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetIncrementalFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIncrementalFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalFriendsLogic {
	return &GetIncrementalFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetIncrementalFriendsLogic) GetIncrementalFriends(req *user.GetIncrementalFriendsReq) (*user.GetIncrementalFriendsResp, error) {
	resp := &user.GetIncrementalFriendsResp{}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// ������֤
	if req.UserID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	// TODO: ʵ������ͬ���߼�
	// ������Ҫ���ݰ汾�Ż�ȡ�����ĺ��ѱ��
	// ��Ҫ�汾����ģ��֧�֣�
	// 1. ���ݰ汾�Ų�ѯ�汾��־
	// 2. ��ȡ���������¡�ɾ���ĺ����б�
	// 3. ���ذ汾��Ϣ�ͱ���б�

	// ��ʱʵ�֣�����汾��Ϊ0��δ�ṩ�����������б��Full=true��
	// ע�⣺����������ͬ����Ҫ�汾����ģ��֧��
	if req.Version == 0 && req.VersionID == "" {
		// ��ȡ���к���ID
		friendIDs, err := l.svcCtx.FriendDB.FindFriendUserIDs(l.ctx, req.UserID)
		if err != nil {
			return nil, err
		}

		// ��ȡ������ϸ��Ϣ
		if len(friendIDs) > 0 {
			friends, err := l.svcCtx.FriendDB.FindFriends(l.ctx, req.UserID, friendIDs)
			if err != nil {
				return nil, err
			}

			// 转换为 sdkws.FriendInfo
			relationFriends := convert.ModelFriendsDB2Pb(friends)
			resp.Insert = make([]*sdkws.FriendInfo, 0, len(relationFriends))
			for _, rf := range relationFriends {
				sdkwsFriend := &sdkws.FriendInfo{
					OwnerUserID:    rf.OwnerUserID,
					Remark:         rf.Remark,
					CreateTime:     rf.CreateTime,
					AddSource:      rf.AddSource,
					OperatorUserID: rf.OperatorUserID,
					Ex:             rf.Ex,
					IsPinned:       rf.IsPinned,
				}
				resp.Insert = append(resp.Insert, sdkwsFriend)
			}

			// �����ѵ��û���Ϣ
			if len(resp.Insert) > 0 {
				// ��ȡ���к��ѵ��û�ID
				friendUserIDs := make([]string, 0, len(friends))
				for _, f := range friends {
					friendUserIDs = append(friendUserIDs, f.FriendUserID)
				}

				// ��ȡ�û���Ϣ
				users, err := l.svcCtx.UserDB.Find(l.ctx, friendUserIDs)
				if err == nil {
					// �����û�ӳ��
					userMap := make(map[string]*model.User)
					for _, u := range users {
						userMap[u.UserID] = u
					}

					// �����ѵ��û���Ϣ
					for i, friendInfo := range resp.Insert {
						if i < len(friends) {
							friend := friends[i]
							if user, ok := userMap[friend.FriendUserID]; ok {
								friendInfo.FriendUser = &sdkws.UserInfo{
									UserID:   user.UserID,
									Nickname: user.Nickname,
									FaceURL:  user.FaceURL,
									Ex:       user.Ex,
								}
							}
						}
					}
				}
			}
		} else {
			resp.Insert = []*sdkws.FriendInfo{}
		}
		resp.Full = true
	} else {
		// ����ͬ�������ؿ��б����Ҫ�汾����֧�֣�
		// TODO: ʵ�ְ汾�����߼�
		// 1. ���ݰ汾�Ų�ѯ�汾��־
		// 2. ��ȡ���������¡�ɾ���ĺ����б�
		// 3. ���ذ汾��Ϣ�ͱ���б�
		resp.Insert = []*sdkws.FriendInfo{}
		resp.Update = []*sdkws.FriendInfo{}
		resp.Delete = []string{}
		resp.Full = false
	}

	// TODO: �Ӱ汾����ģ���ȡ�汾��Ϣ
	resp.Version = 0
	resp.VersionID = ""
	resp.SortVersion = 0

	return resp, nil
}
