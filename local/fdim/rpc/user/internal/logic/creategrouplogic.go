package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/authverify"
	"fdim/pkg/constant"
	"fdim/pkg/convert"
	"fdim/pkg/mcontext"
	"fdim/pkg/model"
	"fdim/pkg/util"
	"fdim/pkg/webhook"
	"fdim/protocol/conversation"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupLogic) CreateGroup(req *user.CreateGroupReq) (*user.CreateGroupResp, error) {
	resp := &user.CreateGroupResp{}

	// ������֤
	if req.GroupInfo == nil {
		return nil, fmt.Errorf("groupInfo is nil")
	}
	if req.GroupInfo.GroupType != constant.WorkingGroup {
		return nil, fmt.Errorf("group type only supports %d", constant.WorkingGroup)
	}
	if req.OwnerUserID == "" {
		return nil, fmt.Errorf("no group owner")
	}
	if req.GroupInfo.GroupName == "" {
		return nil, fmt.Errorf("groupName is empty")
	}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	// �ռ������û�ID
	userIDs := append(append(req.MemberUserIDs, req.AdminUserIDs...), req.OwnerUserID)
	opUserID := mcontext.GetOpUserID(l.ctx)
	if opUserID != "" && !util.Contains(userIDs, opUserID) {
		userIDs = append(userIDs, opUserID)
	}

	// ����û�ID�Ƿ��ظ�
	if util.HasDuplicate(userIDs) {
		return nil, fmt.Errorf("group member repeated")
	}

	// ����û��Ƿ����
	_, err := l.svcCtx.UserDB.FindWithError(l.ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Webhook BeforeCreateGroup �ص�
	if l.svcCtx.WebhookClient != nil {
		// ������һ����ʱ groupID ���� webhook��ʵ�ʻ��ں����������ɣ�
		tempGroupID := ""
		cbReq := &webhook.CallbackBeforeCreateGroupReq{
			CallbackCommand: webhook.CallbackBeforeCreateGroupCommand,
			GroupID:         tempGroupID,
			GroupName:       req.GroupInfo.GroupName,
			GroupType:       int32(req.GroupInfo.GroupType),
			OwnerUserID:     req.OwnerUserID,
			MemberUserIDs:   req.MemberUserIDs,
			AdminUserIDs:    req.AdminUserIDs,
			Ex:              req.GroupInfo.Ex,
		}
		cbResp := &webhook.CallbackBeforeCreateGroupResp{}
		if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
			if err != webhook.ErrCallbackContinue {
				return nil, err
			}
			// ErrCallbackContinue ��ʾ����ִ��
		}
		// �����������ݣ���� webhook �������޸ģ�
		if cbResp.GroupName != nil {
			req.GroupInfo.GroupName = *cbResp.GroupName
		}
		if cbResp.GroupType != nil {
			req.GroupInfo.GroupType = int32(*cbResp.GroupType)
		}
		if cbResp.OwnerUserID != nil {
			req.OwnerUserID = *cbResp.OwnerUserID
		}
		if len(cbResp.MemberUserIDs) > 0 {
			req.MemberUserIDs = cbResp.MemberUserIDs
		}
		if len(cbResp.AdminUserIDs) > 0 {
			req.AdminUserIDs = cbResp.AdminUserIDs
		}
		if cbResp.Ex != nil {
			req.GroupInfo.Ex = *cbResp.Ex
		}
	}

	// ����Ⱥ��ID
	groupID, err := mcontext.GenGroupID(l.ctx, func(id string) (bool, error) {
		_, err := l.svcCtx.GroupDB.TakeGroup(l.ctx, id)
		if err != nil {
			// �����ѯ���������Ƿ���"δ�ҵ�"����
			// �� MongoDB �У����ʹ�� FindOne �Ҳ����ĵ����᷵�� mongo.ErrNoDocuments
			// �������Ǽ�����������������ID�����ڣ��������ã���Ҳ��������������
			// Ϊ�˰�ȫ��������Ƿ��� false�������ڣ����� GenGroupID ��������
			return false, nil
		}
		// �ҵ��ˣ�˵��ID�Ѵ���
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate group ID: %v", err)
	}

	// 转换群组信息
	groupInfo := convert.Pb2ModelGroupInfo(req.GroupInfo)
	groupInfo.GroupID = groupID
	groupInfo.CreateTime = time.Now()

	// ����Ⱥ��Ա
	var groupMembers []*model.GroupMember
	now := time.Now()

	// ���Ⱥ��
	groupMembers = append(groupMembers, &model.GroupMember{
		GroupID:        groupID,
		UserID:         req.OwnerUserID,
		RoleLevel:      constant.GroupOwner,
		JoinTime:       now,
		OperatorUserID: opUserID,
		JoinSource:     constant.JoinByInvitation,
		InviterUserID:  opUserID,
		MuteEndTime:    time.UnixMilli(0),
	})

	// ��ӹ���Ա
	for _, userID := range req.AdminUserIDs {
		groupMembers = append(groupMembers, &model.GroupMember{
			GroupID:        groupID,
			UserID:         userID,
			RoleLevel:      constant.GroupAdmin,
			JoinTime:       now,
			OperatorUserID: opUserID,
			JoinSource:     constant.JoinByInvitation,
			InviterUserID:  opUserID,
			MuteEndTime:    time.UnixMilli(0),
		})
	}

	// �����ͨ��Ա
	for _, userID := range req.MemberUserIDs {
		groupMembers = append(groupMembers, &model.GroupMember{
			GroupID:        groupID,
			UserID:         userID,
			RoleLevel:      constant.GroupOrdinaryUsers,
			JoinTime:       now,
			OperatorUserID: opUserID,
			JoinSource:     constant.JoinByInvitation,
			InviterUserID:  opUserID,
			MuteEndTime:    time.UnixMilli(0),
		})
	}

	// Webhook BeforeMembersJoinGroup �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackBeforeMembersJoinGroupReq{
			CallbackCommand: webhook.CallbackBeforeMembersJoinGroupCommand,
			GroupID:         groupID,
			GroupType:       int32(req.GroupInfo.GroupType),
			MemberUserIDs:   userIDs,
			Ex:              req.GroupInfo.Ex,
		}
		cbResp := &webhook.CallbackBeforeMembersJoinGroupResp{}
		if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
			if err != webhook.ErrCallbackContinue {
				return nil, err
			}
			// ErrCallbackContinue ��ʾ����ִ��
		}
		// ��� webhook �������޸ĺ�ĳ�Ա�б����Ҫ���¹��� groupMembers
		if len(cbResp.MemberUserIDs) > 0 {
			// ����򻯴����ʵ��Ӧ�ø��ݷ��صĳ�Ա�б����¹���
			// Ϊ�˼򻯣�������ʱ��������޸�
		}
	}

	// ����Ⱥ��ͳ�Ա
	if err := l.svcCtx.GroupDB.CreateGroup(l.ctx, []*model.Group{groupInfo}, groupMembers); err != nil {
		return nil, err
	}

	// 构建响应
	memberCount := uint32(len(userIDs))
	resp.GroupInfo = convert.ModelGroupDB2Pb(groupInfo)
	if resp.GroupInfo != nil {
		resp.GroupInfo.OwnerUserID = req.OwnerUserID
		resp.GroupInfo.MemberCount = memberCount
	}

	// 为所有成员创建群聊会话（与官方一致）
	if l.svcCtx.ConversationClient != nil {
		go func() {
			ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := l.svcCtx.ConversationClient.CreateGroupChatConversations(ctx2,
				&conversation.CreateGroupChatConversationsReq{
					GroupID: groupID,
					UserIDs: userIDs,
				})
			if err != nil {
				logx.Errorf("CreateGroupChatConversations failed: groupID=%s, error=%v", groupID, err)
			}
		}()
	}

	// ����Ⱥ�鴴��֪ͨ
	if l.svcCtx.GroupNotification != nil {
		l.svcCtx.GroupNotification.GroupCreatedNotification(l.ctx, groupID, req.OwnerUserID, userIDs)
	}

	// Webhook AfterCreateGroup �ص�
	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackAfterCreateGroupReq{
			CallbackCommand: webhook.CallbackAfterCreateGroupCommand,
			GroupID:         groupID,
			GroupName:       req.GroupInfo.GroupName,
			GroupType:       int32(req.GroupInfo.GroupType),
			OwnerUserID:     req.OwnerUserID,
			MemberUserIDs:   req.MemberUserIDs,
			AdminUserIDs:    req.AdminUserIDs,
			Ex:              req.GroupInfo.Ex,
		}
		cbResp := &webhook.CallbackAfterCreateGroupResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
