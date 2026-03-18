package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetUsersInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersInfoLogic {
	return &GetUsersInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersInfoLogic) GetUsersInfo(req *types.GetUsersInfoReq) (resp *types.GetUsersInfoResp, err error) {
	// 转换请求参数
	rpcReq := &user.GetDesignateUsersReq{
		UserIDs: req.UserIDs,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.GetDesignateUsers(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var usersData []types.UserInfo
	for _, u := range rpcResp.UsersInfo {
		usersData = append(usersData, types.UserInfo{
			UserID:   u.UserID,
			Nickname: u.Nickname,
			FaceURL:  u.FaceURL,
			Ex:       u.Ex,
		})
	}

	resp = &types.GetUsersInfoResp{
		UsersData: usersData,
	}

	return resp, nil
}
