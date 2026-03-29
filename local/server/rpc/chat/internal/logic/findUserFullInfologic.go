package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/protocol/sdkws"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserFullInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserFullInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserFullInfoLogic {
	return &FindUserFullInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserFullInfoLogic) FindUserFullInfo(req *chat.FindUserFullInfoReq) (*chat.FindUserFullInfoResp, error) {
	// 1. ֤
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("user IDs cannot be empty")
	}

	// 2. ûϢ
	userInfos, err := l.svcCtx.ChatDB.FindUserFullInfo(l.ctx, req.UserIDs)
	if err != nil {
		l.Errorf("FindUserFullInfo failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to find user full info")
	}

	// 3. תΪprotobufʽ
	users := make([]*sdkws.ChatUserFullInfo, 0, len(userInfos))
	for _, info := range userInfos {
		users = append(users, &sdkws.ChatUserFullInfo{
			UserID:           info.UserID,
			Account:          info.Account,
			PhoneNumber:      info.PhoneNumber,
			AreaCode:         info.AreaCode,
			Email:            info.Email,
			Nickname:         info.Nickname,
			FaceURL:          info.FaceURL,
			Gender:           info.Gender,
			Level:            info.Level,
			Birth:            info.Birth,
			AllowAddFriend:   info.AllowAddFriend,
			AllowBeep:        info.AllowBeep,
			AllowVibration:   info.AllowVibration,
			GlobalRecvMsgOpt: info.GlobalRecvMsgOpt,
			RegisterType:     info.RegisterType,
		})
	}

	return &chat.FindUserFullInfoResp{
		Users: users,
	}, nil
}
