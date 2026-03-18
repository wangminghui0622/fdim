package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/constant"
	"fdim/pkg/convert"
	"fdim/pkg/model"
	"fdim/pkg/util"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaginationUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaginationUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaginationUsersLogic {
	return &GetPaginationUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPaginationUsersLogic) GetPaginationUsers(req *user.GetPaginationUsersReq) (*user.GetPaginationUsersResp, error) {
	resp := &user.GetPaginationUsersResp{}

	// Ȩ����֤����Ҫ����ԱȨ��
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	// �����ҳ����
	offset := util.CalculateOffset(req.Pagination.PageNumber, req.Pagination.ShowNumber)
	limit := util.CalculateLimit(req.Pagination.ShowNumber)

	var total int64
	var users []*model.User
	var err error

	// �����Ƿ��йؼ���ѡ��ͬ�Ĳ�ѯ����
	if req.UserID == "" && req.NickName == "" {
		// ��ѯ�����û�
		total, users, err = l.svcCtx.UserDB.PageFindUser(l.ctx, constant.IMOrdinaryUser, constant.AppOrdinaryUsers, offset, limit)
	} else {
		// ���ݹؼ��ʲ�ѯ
		total, users, err = l.svcCtx.UserDB.PageFindUserWithKeyword(l.ctx, constant.IMOrdinaryUser, constant.AppOrdinaryUsers, req.UserID, req.NickName, offset, limit)
	}

	if err != nil {
		return nil, err
	}

	// 转换为 Protocol Buffer 格式
	resp.Total = int32(total)
	resp.Users = convert.ModelUsersDB2Pb(users)

	return resp, nil
}
