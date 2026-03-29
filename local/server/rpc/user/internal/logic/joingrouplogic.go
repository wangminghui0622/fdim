package logic

import (
	"context"
	"fmt"
	"time"

	"fdim/pkg/constant"
	"fdim/pkg/mcontext"
	"fdim/pkg/model"
	"fdim/pkg/webhook"
	"fdim/protocol/conversation"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinGroupLogic) JoinGroup(req *user.JoinGroupReq) (*user.JoinGroupResp, error) {
	resp := &user.JoinGroupResp{}

	// ????????????
	groupInfo, err := l.svcCtx.GroupDB.TakeGroup(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// ????????
	if groupInfo.Status == constant.GroupStatusDismissed {
		return nil, fmt.Errorf("group is dismissed")
	}

	// ???????
	opUserID := mcontext.GetOpUserID(l.ctx)
	userID := req.InviterUserID
	if userID == "" {
		if opUserID == "" {
			return nil, fmt.Errorf("userID is empty")
		}
		userID = opUserID
	}

	// ???????????????
	_, err = l.svcCtx.GroupDB.TakeGroupMember(l.ctx, req.GroupID, userID)
	if err == nil {
		return nil, fmt.Errorf("user already in group")
	}

	// ?????????????????????????????????
	// NeedVerification: 0=??????, 1=??????
	const Directly = 0
	const AllNeedVerification = 1

	if groupInfo.NeedVerification == Directly {
		// ??????
		if opUserID == "" {
			opUserID = userID
		}

		groupMember := &model.GroupMember{
			GroupID:        req.GroupID,
			UserID:         userID,
			RoleLevel:      constant.GroupOrdinaryUsers,
			JoinTime:       time.Now(),
			OperatorUserID: opUserID,
			JoinSource:     int32(req.JoinSource),
			InviterUserID:  req.InviterUserID,
			MuteEndTime:    time.UnixMilli(0),
		}

		// Webhook BeforeMembersJoinGroup ???
		if l.svcCtx.WebhookClient != nil {
			cbReq := &webhook.CallbackBeforeMembersJoinGroupReq{
				CallbackCommand: webhook.CallbackBeforeMembersJoinGroupCommand,
				GroupID:         req.GroupID,
				GroupType:       int32(groupInfo.GroupType),
				MemberUserIDs:   []string{userID},
				Ex:              groupInfo.Ex,
			}
			cbResp := &webhook.CallbackBeforeMembersJoinGroupResp{}
			if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
				if err != webhook.ErrCallbackContinue {
					return nil, err
				}
				// ErrCallbackContinue ??????????
			}
			// ??? webhook ??????????????????????1??? groupMember
			if len(cbResp.MemberUserIDs) > 0 {
				// ??????????????????????????????1???
			}
		}

		if err := l.svcCtx.GroupDB.CreateGroup(l.ctx, nil, []*model.GroupMember{groupMember}); err != nil {
			return nil, err
		}

		// Ϊ³ԱȺĻỰٷһ£
		if l.svcCtx.ConversationClient != nil {
			go func() {
				ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, err := l.svcCtx.ConversationClient.CreateGroupChatConversations(ctx2,
					&conversation.CreateGroupChatConversationsReq{
						GroupID: req.GroupID,
						UserIDs: []string{userID},
					})
				if err != nil {
					logx.Errorf("CreateGroupChatConversations failed on join: groupID=%s, error=%v", req.GroupID, err)
				}
			}()
		}

		if l.svcCtx.GroupNotification != nil {
			l.svcCtx.GroupNotification.MemberEnterNotification(l.ctx, req.GroupID, req.InviterUserID, []string{userID})
		}

		// Webhook AfterJoinGroup ???
		if l.svcCtx.WebhookClient != nil {
			cbReq := &webhook.CallbackAfterJoinGroupReq{
				CallbackCommand: webhook.CallbackAfterJoinGroupCommand,
				GroupID:         req.GroupID,
				UserID:          userID,
				ReqMessage:      req.ReqMessage,
				JoinSource:      req.JoinSource,
				InviterUserID:   req.InviterUserID,
			}
			cbResp := &webhook.CallbackAfterJoinGroupResp{}
			l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
		}
	} else {
		// ????????
		groupRequest := &model.GroupRequest{
			UserID:        userID,
			GroupID:       req.GroupID,
			ReqMsg:        req.ReqMessage,
			JoinSource:    int32(req.JoinSource),
			InviterUserID: req.InviterUserID,
			ReqTime:       time.Now(),
			HandledTime:   time.Unix(0, 0),
			Ex:            req.Ex,
			HandleResult:  constant.GroupRequestUnhandled, // ????
		}

		if err := l.svcCtx.GroupDB.CreateGroupRequest(l.ctx, []*model.GroupRequest{groupRequest}); err != nil {
			return nil, err
		}

		// ????????????????
		if l.svcCtx.GroupNotification != nil {
			l.svcCtx.GroupNotification.JoinGroupApplicationNotification(l.ctx, req.GroupID, userID)
		}
	}

	return resp, nil
}
