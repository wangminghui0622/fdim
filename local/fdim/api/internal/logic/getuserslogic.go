package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/sdkws"
	"fdim/protocol/user"
)

type GetUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersLogic {
	return &GetUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersLogic) GetUsers(req *types.GetUsersReq) (resp *types.GetUsersResp, err error) {
	// 转换请求参数
	rpcReq := &user.GetPaginationUsersReq{
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
		UserID:   req.UserID,
		NickName: req.Nickname,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.UserClient.GetPaginationUsers(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var users []types.UserInfo
	for _, u := range rpcResp.Users {
		users = append(users, types.UserInfo{
			UserID:   u.UserID,
			Nickname: u.Nickname,
			FaceURL:  u.FaceURL,
			Ex:       u.Ex,
		})
	}

	resp = &types.GetUsersResp{
		Total: rpcResp.Total,
		Users: users,
	}

	return resp, nil
}
