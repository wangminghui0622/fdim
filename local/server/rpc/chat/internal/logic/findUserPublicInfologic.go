package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/protocol/sdkws"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserPublicInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserPublicInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserPublicInfoLogic {
	return &FindUserPublicInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserPublicInfoLogic) FindUserPublicInfo(req *chat.FindUserPublicInfoReq) (*chat.FindUserPublicInfoResp, error) {
	// 1. ֤
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("user IDs cannot be empty")
	}

	// 2. ûϢ
	userInfos, err := l.svcCtx.ChatDB.FindUserPublicInfo(l.ctx, req.UserIDs)
	if err != nil {
		l.Errorf("FindUserPublicInfo failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to find user public info")
	}

	// 3. תΪprotobufʽ
	users := make([]*sdkws.ChatUserPublicInfo, 0, len(userInfos))
	for _, info := range userInfos {
		users = append(users, &sdkws.ChatUserPublicInfo{
			UserID:   info.UserID,
			Account:  info.Account,
			Nickname: info.Nickname,
			FaceURL:  info.FaceURL,
			Gender:   info.Gender,
			Level:    info.Level,
		})
	}

	return &chat.FindUserPublicInfoResp{
		Users: users,
	}, nil
}
