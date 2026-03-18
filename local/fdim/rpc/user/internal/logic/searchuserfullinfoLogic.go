package logic

import (
	"context"

	"fdim/protocol/sdkws"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
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

// SearchUserFullInfo 搜索用户完整信息
func (l *SearchUserFullInfoLogic) SearchUserFullInfo(req *user.SearchUserFullInfoReq) (*user.SearchUserFullInfoResp, error) {
	resp := &user.SearchUserFullInfoResp{}

	// 计算分页参数
	offset := int32(0)
	limit := int32(20)
	if req.Pagination != nil && req.Pagination.PageNumber > 0 && req.Pagination.ShowNumber > 0 {
		offset = (req.Pagination.PageNumber - 1) * req.Pagination.ShowNumber
		limit = req.Pagination.ShowNumber
	}

	// 使用数据库模糊搜索用户
	total, users, err := l.svcCtx.UserDB.SearchUsers(l.ctx, req.Keyword, offset, limit)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var result []*sdkws.UserInfo
	for _, u := range users {
		result = append(result, &sdkws.UserInfo{
			UserID:           u.UserID,
			Nickname:         u.Nickname,
			FaceURL:          u.FaceURL,
			Ex:               u.Ex,
			CreateTime:       u.CreateTime.UnixMilli(),
			AppMangerLevel:   u.AppMangerLevel,
			GlobalRecvMsgOpt: u.GlobalRecvMsgOpt,
		})
	}

	resp.Total = uint32(total)
	resp.Users = result

	return resp, nil
}
