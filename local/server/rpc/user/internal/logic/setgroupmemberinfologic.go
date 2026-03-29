package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/constant"
	"fdim/pkg/mcontext"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetGroupMemberInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetGroupMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupMemberInfoLogic {
	return &SetGroupMemberInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetGroupMemberInfoLogic) SetGroupMemberInfo(req *user.SetGroupMemberInfoReq) (*user.SetGroupMemberInfoResp, error) {
	resp := &user.SetGroupMemberInfoResp{}

	// ???????
	if len(req.Members) == 0 {
		return nil, fmt.Errorf("members is empty")
	}

	opUserID := mcontext.GetOpUserID(l.ctx)
	if opUserID == "" {
		return nil, fmt.Errorf("no operator user id")
	}

	// ??????????????
	for _, memberInfo := range req.Members {
		if memberInfo.GroupID == "" {
			return nil, fmt.Errorf("groupID is empty")
		}
		if memberInfo.UserID == "" {
			return nil, fmt.Errorf("userID is empty")
		}

		// ?????????????
		member, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, memberInfo.GroupID, memberInfo.UserID)
		if err != nil {
			return nil, fmt.Errorf("member not found: %s in group %s", memberInfo.UserID, memberInfo.GroupID)
		}

		// ?????
		isAdmin := authverify.IsAdmin(l.ctx)
		if !isAdmin {
			// ?????????????????
			opMember, err := l.svcCtx.GroupDB.TakeGroupMember(l.ctx, memberInfo.GroupID, opUserID)
			if err != nil {
				return nil, fmt.Errorf("operator not in group")
			}

			// ??????????
			if memberInfo.RoleLevel != nil {
				switch memberInfo.RoleLevel.Value {
				case constant.GroupOwner:
					return nil, fmt.Errorf("cannot set to group owner")
				case constant.GroupAdmin, constant.GroupOrdinaryUsers:
					// ????
				default:
					return nil, fmt.Errorf("invalid role level")
				}

				// ?????????????????
				switch opMember.RoleLevel {
				case constant.GroupOwner:
					// ????????????????????????????????
				case constant.GroupAdmin:
					// ??????????????????????????
					if member.RoleLevel == constant.GroupOwner {
						return nil, fmt.Errorf("admin cannot change group owner")
					}
					if member.RoleLevel == constant.GroupAdmin && member.UserID != opUserID {
						return nil, fmt.Errorf("admin cannot change other admin")
					}
				case constant.GroupOrdinaryUsers:
					// ??????????????????????????????????
					if member.UserID != opUserID {
						return nil, fmt.Errorf("ordinary member can only change own info")
					}
					if memberInfo.RoleLevel != nil {
						return nil, fmt.Errorf("ordinary member cannot change role level")
					}
				}

				// ??????????????????
				if member.UserID == opUserID && memberInfo.RoleLevel.Value > member.RoleLevel {
					return nil, fmt.Errorf("cannot improve own role level")
				}

				// ?????????????????
				if member.UserID == opUserID && member.RoleLevel == constant.GroupOwner {
					return nil, fmt.Errorf("group owner cannot change own role level")
				}
			} else {
				// ??????????????????????
				if member.UserID != opUserID {
					// ??????????????????????????
					if opMember.RoleLevel != constant.GroupOwner && opMember.RoleLevel != constant.GroupAdmin {
						return nil, fmt.Errorf("only owner or admin can change other member info")
					}
				}
			}
		}

		// ??? Webhook ???????
		var nickname, faceURL, ex *string
		if memberInfo.Nickname != nil {
			nickname = &memberInfo.Nickname.Value
		}
		if memberInfo.FaceURL != nil {
			faceURL = &memberInfo.FaceURL.Value
		}
		if memberInfo.Ex != nil {
			ex = &memberInfo.Ex.Value
		}

		// Webhook BeforeSetGroupMemberInfo ???
		if l.svcCtx.WebhookClient != nil {
			cbReq := &webhook.CallbackBeforeSetGroupMemberInfoReq{
				CallbackCommand: webhook.CallbackBeforeSetGroupMemberInfoCommand,
				GroupID:         memberInfo.GroupID,
				UserID:          memberInfo.UserID,
				Nickname:        nickname,
				FaceURL:         faceURL,
				Ex:              ex,
			}
			cbResp := &webhook.CallbackBeforeSetGroupMemberInfoResp{}
			if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
				if err != webhook.ErrCallbackContinue {
					return nil, err
				}
				// ErrCallbackContinue ??????????
			}
			// ??? webhook ???????????????????
			if cbResp.Nickname != nil && memberInfo.Nickname != nil {
				memberInfo.Nickname.Value = *cbResp.Nickname
			}
			if cbResp.FaceURL != nil && memberInfo.FaceURL != nil {
				memberInfo.FaceURL.Value = *cbResp.FaceURL
			}
			if cbResp.Ex != nil && memberInfo.Ex != nil {
				memberInfo.Ex.Value = *cbResp.Ex
			}
		}

		// ????????????
		data := make(map[string]interface{})
		if memberInfo.Nickname != nil {
			data["nickname"] = memberInfo.Nickname.Value
		}
		if memberInfo.FaceURL != nil {
			data["face_url"] = memberInfo.FaceURL.Value
		}
		if memberInfo.Ex != nil {
			data["ex"] = memberInfo.Ex.Value
		}
		if memberInfo.RoleLevel != nil {
			data["role_level"] = memberInfo.RoleLevel.Value
		}

		if len(data) > 0 {
			if err := l.svcCtx.GroupDB.UpdateGroupMemberMap(l.ctx, memberInfo.GroupID, memberInfo.UserID, data); err != nil {
				return nil, err
			}
		}

		// ??????
		if memberInfo.RoleLevel != nil {
			switch memberInfo.RoleLevel.Value {
			case constant.GroupAdmin:
				// ????????????????????
				if l.svcCtx.GroupNotification != nil {
					l.svcCtx.GroupNotification.GroupMemberSetToAdminNotification(l.ctx, memberInfo.GroupID, opUserID, []string{memberInfo.UserID})
				}
			case constant.GroupOrdinaryUsers:
				// ?????????????????????
				if l.svcCtx.GroupNotification != nil {
					l.svcCtx.GroupNotification.GroupMemberSetToOrdinaryUserNotification(l.ctx, memberInfo.GroupID, opUserID, []string{memberInfo.UserID})
				}
			}
		}
		if memberInfo.Nickname != nil || memberInfo.FaceURL != nil || memberInfo.Ex != nil {
			// ?????????????????
			if l.svcCtx.GroupNotification != nil {
				l.svcCtx.GroupNotification.GroupMemberInfoSetNotification(l.ctx, memberInfo.GroupID, opUserID, memberInfo.UserID)
			}
		}

		// Webhook AfterSetGroupMemberInfo ???
		if l.svcCtx.WebhookClient != nil {
			nicknameStr := ""
			faceURLStr := ""
			exStr := ""
			if memberInfo.Nickname != nil {
				nicknameStr = memberInfo.Nickname.Value
			}
			if memberInfo.FaceURL != nil {
				faceURLStr = memberInfo.FaceURL.Value
			}
			if memberInfo.Ex != nil {
				exStr = memberInfo.Ex.Value
			}
			cbReq := &webhook.CallbackAfterSetGroupMemberInfoReq{
				CallbackCommand: webhook.CallbackAfterSetGroupMemberInfoCommand,
				GroupID:         memberInfo.GroupID,
				UserID:          memberInfo.UserID,
				Nickname:        nicknameStr,
				FaceURL:         faceURLStr,
				Ex:              exStr,
			}
			cbResp := &webhook.CallbackAfterSetGroupMemberInfoResp{}
			l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
		}
	}

	return resp, nil
}
