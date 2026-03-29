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

	// ???????
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

	// ??????
	if err := authverify.CheckAccess(l.ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	// ??????????ID
	userIDs := append(append(req.MemberUserIDs, req.AdminUserIDs...), req.OwnerUserID)
	opUserID := mcontext.GetOpUserID(l.ctx)
	if opUserID != "" && !util.Contains(userIDs, opUserID) {
		userIDs = append(userIDs, opUserID)
	}

	// ??????ID??????
	if util.HasDuplicate(userIDs) {
		return nil, fmt.Errorf("group member repeated")
	}

	// ????????????
	_, err := l.svcCtx.UserDB.FindWithError(l.ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Webhook BeforeCreateGroup ???
	if l.svcCtx.WebhookClient != nil {
		// ???????????? groupID ???? webhook????????????????????
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
			// ErrCallbackContinue ??????????
		}
		// ???????????????? webhook ??????????
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

	// ???????ID
	groupID, err := mcontext.GenGroupID(l.ctx, func(id string) (bool, error) {
		_, err := l.svcCtx.GroupDB.TakeGroup(l.ctx, id)
		if err != nil {
			// ???????????????????"???"????
			// ?? MongoDB ???????? FindOne ?????????????? mongo.ErrNoDocuments
			// ?????????????????????????ID?????????????????????????????????
			// ????????????????? false????????????? GenGroupID ????????
			return false, nil
		}
		// ?????????ID?????
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate group ID: %v", err)
	}

	// תȺϢ
	groupInfo := convert.Pb2ModelGroupInfo(req.GroupInfo)
	groupInfo.GroupID = groupID
	groupInfo.CreateTime = time.Now()

	// ????????
	var groupMembers []*model.GroupMember
	now := time.Now()

	// ??????
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

	// ???????
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

	// ?????????
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

	// Webhook BeforeMembersJoinGroup ???
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
			// ErrCallbackContinue ??????????
		}
		// ??? webhook ??????????????????????1??? groupMembers
		if len(cbResp.MemberUserIDs) > 0 {
			// ??????????????????????????????1???
			// ???????????????????????
		}
	}

	// ??????????
	if err := l.svcCtx.GroupDB.CreateGroup(l.ctx, []*model.Group{groupInfo}, groupMembers); err != nil {
		return nil, err
	}

	// Ӧ
	memberCount := uint32(len(userIDs))
	resp.GroupInfo = convert.ModelGroupDB2Pb(groupInfo)
	if resp.GroupInfo != nil {
		resp.GroupInfo.OwnerUserID = req.OwnerUserID
		resp.GroupInfo.MemberCount = memberCount
	}

	// ΪгԱȺĻỰٷһ£
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

	// ??????????
	if l.svcCtx.GroupNotification != nil {
		l.svcCtx.GroupNotification.GroupCreatedNotification(l.ctx, groupID, req.OwnerUserID, userIDs)
	}

	// Webhook AfterCreateGroup ???
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
