package logic

import (
	"context"
	"fmt"

	"fdim/pkg/convert"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDesignateUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDesignateUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDesignateUsersLogic {
	return &GetDesignateUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDesignateUsersLogic) GetDesignateUsers(req *user.GetDesignateUsersReq) (*user.GetDesignateUsersResp, error) {
	resp := &user.GetDesignateUsersResp{}

	// ������֤
	if len(req.UserIDs) == 0 {
		return nil, fmt.Errorf("userIDs is empty")
	}

	// ��ѯ�û���Ϣ
	users, err := l.svcCtx.UserDB.Find(l.ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}

	// 转换为 Protocol Buffer 格式
	resp.UsersInfo = convert.ModelUsersDB2Pb(users)
	return resp, nil
}
