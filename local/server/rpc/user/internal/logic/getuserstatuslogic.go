package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserStatusLogic {
	return &GetUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserStatusLogic) GetUserStatus(req *user.GetUserStatusReq) (*user.GetUserStatusResp, error) {
	resp := &user.GetUserStatusResp{}

	// ������֤
	if len(req.UserIDs) == 0 {
		return nil, fmt.Errorf("userIDs is empty")
	}

	// �ӻ����ȡ�û�����״̬
	statusList := make([]*user.OnlineStatus, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		platformIDs, err := l.svcCtx.UserCache.GetUserOnline(l.ctx, userID)
		if err != nil {
			// �����ȡʧ�ܣ���������״̬
			statusList = append(statusList, &user.OnlineStatus{
				UserID:      userID,
				Status:      0, // Offline
				PlatformIDs: []int32{},
			})
			continue
		}

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
	return resp, nil
}
