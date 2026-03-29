package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/util"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllUserIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllUserIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllUserIDLogic {
	return &GetAllUserIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllUserIDLogic) GetAllUserID(req *user.GetAllUserIDReq) (*user.GetAllUserIDResp, error) {
	resp := &user.GetAllUserIDResp{}

	// Ȩ֤ҪԱȨ
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	// ҳ?
	offset := util.CalculateOffset(req.Pagination.PageNumber, req.Pagination.ShowNumber)
	limit := util.CalculateLimit(req.Pagination.ShowNumber)

	// ҳȡûID
	total, userIDs, err := l.svcCtx.UserDB.GetAllUserID(l.ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	resp.Total = int32(total)
	resp.UserIDs = userIDs

	return resp, nil
}
