package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/model"
	"fdim/pkg/util"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaginationBlacksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaginationBlacksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaginationBlacksLogic {
	return &GetPaginationBlacksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPaginationBlacksLogic) GetPaginationBlacks(req *user.GetPaginationBlacksReq) (*user.GetPaginationBlacksResp, error) {
	resp := &user.GetPaginationBlacksResp{}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// �����ҳ����
	offset := util.CalculateOffset(req.Pagination.PageNumber, req.Pagination.ShowNumber)
	limit := util.CalculateLimit(req.Pagination.ShowNumber)

	// ��ҳ��ѯ������
	total, blacks, err := l.svcCtx.BlackDB.PageOwnerBlacks(l.ctx, req.UserID, offset, limit)
	if err != nil {
		return nil, err
	}

	// 转换为 Protocol Buffer 格式
	resp.Total = int32(total)
	resp.Blacks = convert.ModelBlacksDB2Pb(blacks)

	// ��ȡ�������û����û���Ϣ
	if len(blacks) > 0 {
		blackUserIDs := make([]string, 0, len(blacks))
		for _, b := range blacks {
			blackUserIDs = append(blackUserIDs, b.BlockUserID)
		}

		users, err := l.svcCtx.UserDB.Find(l.ctx, blackUserIDs)
		if err != nil {
			return nil, err
		}

		// �����û�ӳ��
		userMap := make(map[string]*model.User)
		for _, u := range users {
			userMap[u.UserID] = u
		}

		// ���������û����û���Ϣ
		// ���� convert.BlacksDB2Pb ���ص��б�˳��������� blacks �б�˳��һ�£�����ֱ��ͨ������ƥ��
		for i, blackInfo := range resp.Blacks {
			if i < len(blacks) {
				black := blacks[i]
				if user, ok := userMap[black.BlockUserID]; ok {
					blackInfo.BlackUserInfo = &sdkws.PublicUserInfo{
						UserID:   user.UserID,
						Nickname: user.Nickname,
						FaceURL:  user.FaceURL,
						Ex:       user.Ex,
					}
				}
			}
		}
	}

	return resp, nil
}
