package logic

import (
	"context"
	"strings"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/sdkws"
	"fdim/protocol/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchFriendLogic {
	return &SearchFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchFriendLogic) SearchFriend(req *types.SearchFriendReq) (*types.SearchFriendResp, error) {
	if l.svcCtx.FriendClient == nil {
		return nil, errs.ErrInternalServer.WrapMsg("friend service not available")
	}

	// 获取好友列表
	friendResp, err := l.svcCtx.FriendClient.GetPaginationFriends(l.ctx, &user.GetPaginationFriendsReq{
		UserID: req.UserID,
		Pagination: &sdkws.RequestPagination{
			PageNumber: 1,
			ShowNumber: 1000,
		},
	})
	if err != nil {
		l.Errorf("GetPaginationFriends failed: %v", err)
		return nil, err
	}

	// 本地过滤匹配 keyword
	keyword := strings.ToLower(req.Keyword)
	var matched []interface{}
	if friendResp != nil && friendResp.FriendsInfo != nil {
		for _, f := range friendResp.FriendsInfo {
			if f == nil {
				continue
			}
			remark := strings.ToLower(f.Remark)
			var nick, uid string
			if f.FriendUser != nil {
				nick = strings.ToLower(f.FriendUser.Nickname)
				uid = strings.ToLower(f.FriendUser.UserID)
			}
			if strings.Contains(nick, keyword) ||
				strings.Contains(remark, keyword) ||
				strings.Contains(uid, keyword) {
				matched = append(matched, f)
			}
		}
	}

	return &types.SearchFriendResp{Friends: matched}, nil
}
