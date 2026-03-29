package logic

import (
	"context"
	"fmt"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/model"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSpecifiedBlacksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSpecifiedBlacksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSpecifiedBlacksLogic {
	return &GetSpecifiedBlacksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSpecifiedBlacksLogic) GetSpecifiedBlacks(req *user.GetSpecifiedBlacksReq) (*user.GetSpecifiedBlacksResp, error) {
	resp := &user.GetSpecifiedBlacksResp{}

	// ???????
	if len(req.UserIDList) == 0 {
		return nil, fmt.Errorf("userIDList is empty")
	}

	// ??????
	if err := authverify.CheckAccess(l.ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	// ?????????
	blacks, err := l.svcCtx.BlackDB.FindBlackInfos(l.ctx, req.OwnerUserID, req.UserIDList)
	if err != nil {
		return nil, err
	}

	// ????????????????????
	users, err := l.svcCtx.UserDB.Find(l.ctx, req.UserIDList)
	if err != nil {
		return nil, err
	}

	// ??????????
	userMap := make(map[string]*model.User)
	for _, u := range users {
		userMap[u.UserID] = u
	}

	// ?????????????
	blackMap := make(map[string]*model.Black)
	for _, b := range blacks {
		blackMap[b.BlockUserID] = b
	}

	// תΪ Protocol Buffer ʽûϢ
	result := make([]*sdkws.BlackInfo, 0, len(req.UserIDList))
	for _, userID := range req.UserIDList {
		if black, ok := blackMap[userID]; ok {
			blackInfo := convert.ModelBlackDB2Pb(black, nil)
			result = append(result, blackInfo)
		}
	}

	resp.Blacks = result
	resp.Total = int32(len(result))

	return resp, nil
}
