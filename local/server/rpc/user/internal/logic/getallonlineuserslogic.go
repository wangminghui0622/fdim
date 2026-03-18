package logic

import (
	"context"
	"strconv"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllOnlineUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllOnlineUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllOnlineUsersLogic {
	return &GetAllOnlineUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllOnlineUsersLogic) GetAllOnlineUsers(req *user.GetAllOnlineUsersReq) (*user.GetAllOnlineUsersResp, error) {
	resp := &user.GetAllOnlineUsersResp{}

	// �� uint64 cursor ת��Ϊ string
	cursorStr := ""
	if req.Cursor > 0 {
		cursorStr = strconv.FormatUint(req.Cursor, 10)
	}

	// �� Redis ��ȡ���������û�
	userMap, nextCursorStr, err := l.svcCtx.UserCache.GetAllOnlineUsers(l.ctx, cursorStr)
	if err != nil {
		return nil, err
	}

	// ת��Ϊ OnlineStatus �б�
	statusList := make([]*user.OnlineStatus, 0, len(userMap))
	for userID, platformIDs := range userMap {
		status := int32(0) // Offline
		if len(platformIDs) > 0 {
			status = 1 // Online
		}

		statusList = append(statusList, &user.OnlineStatus{
			UserID:      userID,
			Status:      status,
			PlatformIDs: platformIDs,
		})
	}

	resp.StatusList = statusList
	// �� string cursor ת��Ϊ uint64
	if nextCursorStr != "" {
		nextCursor, err := strconv.ParseUint(nextCursorStr, 10, 64)
		if err == nil {
			resp.NextCursor = nextCursor
		}
	}
	return resp, nil
}
