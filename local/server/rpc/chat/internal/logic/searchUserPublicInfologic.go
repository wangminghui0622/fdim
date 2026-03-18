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

type SearchUserPublicInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserPublicInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserPublicInfoLogic {
	return &SearchUserPublicInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserPublicInfoLogic) SearchUserPublicInfo(req *chat.SearchUserPublicInfoReq) (*chat.SearchUserPublicInfoResp, error) {
	// 1. 搜索用户公开信息
	total, userInfos, err := l.svcCtx.ChatDB.SearchUserPublicInfo(l.ctx, req.Keyword, req.Pagination)
	if err != nil {
		l.Errorf("SearchUserPublicInfo failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to search user public info")
	}

	// 2. 过滤性别（如果指定）
	var filteredInfos []*database.UserPublicInfo
	if req.Genders != 0 {
		for _, info := range userInfos {
			if info.Gender == req.Genders {
				filteredInfos = append(filteredInfos, info)
			}
		}
		userInfos = filteredInfos
	}

	// 3. 转换为protobuf格式
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

	return &chat.SearchUserPublicInfoResp{
		Total: uint32(total),
		Users: users,
	}, nil
}
