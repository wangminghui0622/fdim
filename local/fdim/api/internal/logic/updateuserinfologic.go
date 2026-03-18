package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/sdkws"
	"fdim/protocol/user"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (resp *types.UpdateUserInfoResp, err error) {
	// 转换请求参数
	rpcReq := &user.UpdateUserInfoReq{
		UserInfo: &sdkws.UserInfo{
			UserID:   req.UserInfo.UserID,
			Nickname: req.UserInfo.Nickname,
			FaceURL:  req.UserInfo.FaceURL,
			Ex:       req.UserInfo.Ex,
		},
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.UserClient.UpdateUserInfo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.UpdateUserInfoResp{
	}

	return resp, nil
}
