package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/convert"
	"fdim/pkg/model"
	"fdim/pkg/util"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSelfApplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSelfApplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelfApplyListLogic {
	return &GetSelfApplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSelfApplyListLogic) GetSelfApplyList(req *user.GetPaginationFriendsApplyFromReq) (*user.GetPaginationFriendsApplyFromResp, error) {
	resp := &user.GetPaginationFriendsApplyFromResp{}

	// Ȩ����֤
	if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
		return nil, err
	}

	// ת�� handleResults
	var handleResults []int32
	if len(req.HandleResults) > 0 {
		handleResults = req.HandleResults
	}

	// �����ҳ����
	offset := util.CalculateOffset(req.Pagination.PageNumber, req.Pagination.ShowNumber)
	limit := util.CalculateLimit(req.Pagination.ShowNumber)

	// ��ҳ��ѯ�Լ����͵ĺ�������
	total, requests, err := l.svcCtx.FriendDB.PageFriendRequestFromMe(
		l.ctx,
		req.UserID,
		handleResults,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}

	// 转换为 Protocol Buffer 格式
	resp.Total = int32(total)
	resp.FriendRequests = convert.ModelFriendRequestsDB2Pb(requests)

	// ����û���Ϣ
	if len(requests) > 0 {
		userIDs := make([]string, 0)
		for _, req := range requests {
			userIDs = append(userIDs, req.ToUserID)
		}

		users, err := l.svcCtx.UserDB.Find(l.ctx, userIDs)
		if err != nil {
			return nil, err
		}

		// �����û�ӳ��
		userMap := make(map[string]*model.User)
		for _, u := range users {
			userMap[u.UserID] = u
		}

		// ������������û���Ϣ
		for i, friendRequest := range resp.FriendRequests {
			if i < len(requests) {
				req := requests[i]
				if user, ok := userMap[req.ToUserID]; ok {
					friendRequest.ToNickname = user.Nickname
					friendRequest.ToFaceURL = user.FaceURL
				}
			}
		}
	}

	return resp, nil
}
