package logic

import (
	"context"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchDefaultFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchDefaultFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchDefaultFriendLogic {
	return &SearchDefaultFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchDefaultFriendLogic) SearchDefaultFriend(req *admin.SearchDefaultFriendReq) (*admin.SearchDefaultFriendResp, error) {
	// 1. 验证分页参数
	if req.Pagination == nil || req.Pagination.PageNumber <= 0 || req.Pagination.ShowNumber <= 0 {
		return nil, errs.ErrArgs.WrapMsg("invalid pagination parameters")
	}

	// 2. 搜索默认好友
	total, defaultFriends, err := l.svcCtx.AdminDB.SearchDefaultFriend(l.ctx, req.Keyword, req.Pagination)
	if err != nil {
		l.Errorf("SearchDefaultFriend failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to search default friends")
	}

	// 3. 转换为响应格式
	attributes := make([]*admin.DefaultFriendAttribute, 0, len(defaultFriends))
	for _, df := range defaultFriends {
		attributes = append(attributes, &admin.DefaultFriendAttribute{
			UserID:     df.UserID,
			CreateTime: df.CreateTime.UnixMilli(),
		})
	}

	return &admin.SearchDefaultFriendResp{
		Total: uint32(total),
		Users: attributes,
	}, nil
}
