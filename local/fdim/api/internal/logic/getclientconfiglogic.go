package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetClientConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetClientConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClientConfigLogic {
	return &GetClientConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetClientConfigLogic) GetClientConfig(_ *types.GetClientConfigReq) (*types.GetClientConfigResp, error) {
	// 返回客户端初始化配置（可根据需求扩展）
	config := map[string]interface{}{
		"needInvitationCodeRegister": false,
		"allowSendMsgNotFriend":     true,
		"bossUserID":               "",
	}
	return &types.GetClientConfigResp{Config: config}, nil
}
