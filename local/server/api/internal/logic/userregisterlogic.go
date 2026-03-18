package logic

import (
	"context"
	"errors"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/sdkws"
	"fdim/protocol/user"
)

type UserRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterReq) (resp *types.UserRegisterResp, err error) {
	// 验证密码
	if req.Password == "" {
		return nil, errors.New("密码不能为空")
	}
	
	// 验证密码强度
	if err := ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// 加密密码
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 调用 RPC 服务
	rpcReq := &user.UserRegisterReq{
		Users: []*sdkws.UserInfo{
			{
				UserID:   req.UserID,
				Nickname: req.Nickname,
				FaceURL:  req.FaceURL,
				Ex:       hashedPassword, // 将加密后的密码存储在 Ex 字段?
			},
		},
	}

	_, err = l.svcCtx.UserClient.UserRegister(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.UserRegisterResp{
	}

	return resp, nil
}
