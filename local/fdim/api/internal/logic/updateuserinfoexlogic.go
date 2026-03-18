package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/sdkws"
	"fdim/protocol/user"
	"fdim/protocol/wrapperspb"
)

type UpdateUserInfoExLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoExLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoExLogic {
	return &UpdateUserInfoExLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoExLogic) UpdateUserInfoEx(req *types.UpdateUserInfoExReq) (resp *types.UpdateUserInfoExResp, err error) {
	// 转换请求参数
	rpcReq := &user.UpdateUserInfoExReq{
		UserInfo: &sdkws.UserInfoWithEx{
			UserID: req.UserInfo.UserID,
		},
	}

	if req.UserInfo.Nickname != nil {
		rpcReq.UserInfo.Nickname = wrapperspb.String(*req.UserInfo.Nickname)
	}
	if req.UserInfo.FaceURL != nil {
		rpcReq.UserInfo.FaceURL = wrapperspb.String(*req.UserInfo.FaceURL)
	}
	if req.UserInfo.Ex != nil {
		rpcReq.UserInfo.Ex = wrapperspb.String(*req.UserInfo.Ex)
	}

	// 调用 RPC 服务
	_, err = l.svcCtx.UserClient.UpdateUserInfoEx(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.UpdateUserInfoExResp{
	}

	return resp, nil
}
