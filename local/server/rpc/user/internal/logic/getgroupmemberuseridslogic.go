package logic

import (
	"context"
	"fmt"

	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberUserIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberUserIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberUserIDsLogic {
	return &GetGroupMemberUserIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMemberUserIDsLogic) GetGroupMemberUserIDs(req *user.GetGroupMemberUserIDsReq) (*user.GetGroupMemberUserIDsResp, error) {
	resp := &user.GetGroupMemberUserIDsResp{}

	// 参数校验
	if req.GroupID == "" {
		return nil, fmt.Errorf("groupID is empty")
	}

	// 注意：此接口被内部服务（如 push、msgtransfer）调用时，context 中可能没有 opUserID
	// 因此移除权限检查，改为信任内部服务调用
	// 如果需要权限控制，应该在 API 层进行，而不是 RPC 层

	// 获取群成员列表
	members, err := l.svcCtx.GroupDB.FindGroupMemberAll(l.ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	// 提取用户ID列表
	userIDs := make([]string, 0, len(members))
	for _, m := range members {
		userIDs = append(userIDs, m.UserID)
	}

	resp.UserIDs = userIDs
	return resp, nil
}
