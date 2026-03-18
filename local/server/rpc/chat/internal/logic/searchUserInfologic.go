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

type SearchUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserInfoLogic {
	return &SearchUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserInfoLogic) SearchUserInfo(req *chat.SearchUserInfoReq) (*chat.SearchUserInfoResp, error) {
	// 1. 如果指定了userIDs，直接查找这些用户
	if len(req.UserIDs) > 0 {
		userInfos, err := l.svcCtx.ChatDB.FindUserFullInfo(l.ctx, req.UserIDs)
		if err != nil {
			l.Errorf("FindUserFullInfo failed: %v", err)
			return nil, errs.WrapMsg(err, "failed to find user info")
		}

		// 过滤性别
		var filteredInfos []*database.UserFullInfo
		if len(req.Genders) > 0 {
			genderMap := make(map[int32]bool)
			for _, g := range req.Genders {
				genderMap[g] = true
			}
			for _, info := range userInfos {
				if genderMap[info.Gender] {
					filteredInfos = append(filteredInfos, info)
				}
			}
			userInfos = filteredInfos
		}

		// 转换为protobuf格式
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

		return &chat.SearchUserInfoResp{
			Total: uint32(len(users)),
			Users: users,
		}, nil
	}

	// 2. 使用关键词搜索
	total, userInfos, err := l.svcCtx.ChatDB.SearchUserFullInfo(l.ctx, req.Keyword, req.Pagination)
	if err != nil {
		l.Errorf("SearchUserFullInfo failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to search user info")
	}

	// 3. 过滤性别
	var filteredInfos []*database.UserFullInfo
	if len(req.Genders) > 0 {
		genderMap := make(map[int32]bool)
		for _, g := range req.Genders {
			genderMap[g] = true
		}
		for _, info := range userInfos {
			if genderMap[info.Gender] {
				filteredInfos = append(filteredInfos, info)
			}
		}
		userInfos = filteredInfos
	}

	// 4. 转换为protobuf格式
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

	return &chat.SearchUserInfoResp{
		Total: uint32(total),
		Users: users,
	}, nil
}
