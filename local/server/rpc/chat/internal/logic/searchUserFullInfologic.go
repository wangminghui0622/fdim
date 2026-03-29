package logic

import (
	"context"

	"fdim/pkg/database"
	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/protocol/sdkws"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserFullInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserFullInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserFullInfoLogic {
	return &SearchUserFullInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserFullInfoLogic) SearchUserFullInfo(req *chat.SearchUserFullInfoReq) (*chat.SearchUserFullInfoResp, error) {
	// 1. ûϢ
	total, userInfos, err := l.svcCtx.ChatDB.SearchUserFullInfo(l.ctx, req.Keyword, req.Pagination)
	if err != nil {
		l.Errorf("SearchUserFullInfo failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to search user full info")
	}

	// 2. Աָ
	var filteredInfos []*database.UserFullInfo
	if req.Genders != 0 {
		for _, info := range userInfos {
			if info.Gender == req.Genders {
				filteredInfos = append(filteredInfos, info)
			}
		}
		userInfos = filteredInfos
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

	return &chat.SearchUserFullInfoResp{
		Total: uint32(total),
		Users: users,
	}, nil
}
