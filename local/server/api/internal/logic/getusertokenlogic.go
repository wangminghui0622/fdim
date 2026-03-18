package logic

import (
	"context"
	"errors"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/auth"
)

type GetUserTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserTokenLogic {
	return &GetUserTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserTokenLogic) GetUserToken(req *types.GetUserTokenReq) (resp *types.GetUserTokenResp, err error) {
	// 验证密码是否为空
	if req.Password == "" {
		return nil, errors.New("密码不能为空")
	}

	// 从数据库查询用户信息
	if l.svcCtx.UserDB == nil {
		return nil, errors.New("user database not initialized")
	}

	user, err := l.svcCtx.UserDB.GetUserByID(l.ctx, req.UserID)
	if err != nil {
		return nil, errors.New("用户不存在或密码错误")
	}

	// 验证密码（密码存储在 Ex 字段中）
	if !VerifyPassword(req.Password, user.Ex) {
		return nil, errors.New("用户不存在或密码错误")
	}

	// 转换请求参数
	rpcReq := &auth.GetUserTokenReq{
		PlatformID: req.PlatformID,
		UserID:     req.UserID,
	}

	// 调用 RPC 服务
	rpcResp, err := l.svcCtx.AuthClient.GetUserToken(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.GetUserTokenResp{
		Token:             rpcResp.Token,
		ExpireTimeSeconds: rpcResp.ExpireTimeSeconds,
	}

	return resp, nil
}
