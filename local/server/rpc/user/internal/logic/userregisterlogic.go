package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"fdim/pkg/model"
	"fdim/pkg/webhook"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRegisterLogic) UserRegister(req *user.UserRegisterReq) (*user.UserRegisterResp, error) {
	resp := &user.UserRegisterResp{}

	if len(req.Users) == 0 {
		return nil, fmt.Errorf("users is empty")
	}

	// 注：移除管理员权限检查，允许用户自注?
	// if err := authverify.CheckAdmin(l.ctx); err != nil {
	// 	return nil, err
	// }

	// ûIDǷظ
	userIDs := make([]string, 0, len(req.Users))
	userIDMap := make(map[string]bool)
	for _, u := range req.Users {
		if u.UserID == "" {
			return nil, fmt.Errorf("userID is empty")
		}
		if strings.Contains(u.UserID, ":") {
			return nil, fmt.Errorf("userID contains ':' is invalid userID")
		}
		if userIDMap[u.UserID] {
			return nil, fmt.Errorf("userID repeated: %s", u.UserID)
		}
		userIDMap[u.UserID] = true
		userIDs = append(userIDs, u.UserID)
	}

	// ûǷѴ?
	exist, err := l.svcCtx.UserDB.IsExist(l.ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("userID registered already")
	}

	// Webhook BeforeUserRegister ص
	usersData := make([]map[string]interface{}, 0, len(req.Users))
	for _, u := range req.Users {
		userData := map[string]interface{}{
			"userID":           u.UserID,
			"nickname":         u.Nickname,
			"faceURL":          u.FaceURL,
			"ex":               u.Ex,
			"appMangerLevel":   u.AppMangerLevel,
			"globalRecvMsgOpt": u.GlobalRecvMsgOpt,
		}
		usersData = append(usersData, userData)
	}

	if l.svcCtx.WebhookClient != nil {
		cbReq := &webhook.CallbackBeforeUserRegisterReq{
			CallbackCommand: webhook.CallbackBeforeUserRegisterCommand,
			Users:           usersData,
		}
		cbResp := &webhook.CallbackBeforeUserRegisterResp{}
		if err := l.svcCtx.WebhookClient.SyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30); err != nil {
			if err != webhook.ErrCallbackContinue {
				return nil, err
			}
			// ErrCallbackContinue ʾִ
		}
		// ?webhook ޸ĺûбʹ
		if len(cbResp.Users) > 0 {
			usersData = cbResp.Users
		}
	}

	// û
	now := time.Now()
	users := make([]*model.User, 0, len(usersData))
	for _, userData := range usersData {
		userID, _ := userData["userID"].(string)
		nickname, _ := userData["nickname"].(string)
		faceURL, _ := userData["faceURL"].(string)
		ex, _ := userData["ex"].(string)
		var appMangerLevel int32
		var globalRecvMsgOpt int32
		if val, ok := userData["appMangerLevel"].(float64); ok {
			appMangerLevel = int32(val)
		} else if val, ok := userData["appMangerLevel"].(int32); ok {
			appMangerLevel = val
		}
		if val, ok := userData["globalRecvMsgOpt"].(float64); ok {
			globalRecvMsgOpt = int32(val)
		} else if val, ok := userData["globalRecvMsgOpt"].(int32); ok {
			globalRecvMsgOpt = val
		}

		users = append(users, &model.User{
			UserID:           userID,
			Nickname:         nickname,
			FaceURL:          faceURL,
			Ex:               ex,
			CreateTime:       now,
			AppMangerLevel:   appMangerLevel,
			GlobalRecvMsgOpt: globalRecvMsgOpt,
		})
	}

	if err := l.svcCtx.UserDB.Create(l.ctx, users); err != nil {
		return nil, err
	}

	// Webhook AfterUserRegister ص
	if l.svcCtx.WebhookClient != nil {
		// ¹û After ص
		afterUsersData := make([]map[string]interface{}, 0, len(users))
		for _, u := range users {
			userData := map[string]interface{}{
				"userID":           u.UserID,
				"nickname":         u.Nickname,
				"faceURL":          u.FaceURL,
				"ex":               u.Ex,
				"appMangerLevel":   u.AppMangerLevel,
				"globalRecvMsgOpt": u.GlobalRecvMsgOpt,
			}
			afterUsersData = append(afterUsersData, userData)
		}
		cbReq := &webhook.CallbackAfterUserRegisterReq{
			CallbackCommand: webhook.CallbackAfterUserRegisterCommand,
			Users:           afterUsersData,
		}
		cbResp := &webhook.CallbackAfterUserRegisterResp{}
		l.svcCtx.WebhookClient.AsyncPost(l.ctx, cbReq.GetCallbackCommand(), cbReq, cbResp, 30)
	}

	return resp, nil
}
