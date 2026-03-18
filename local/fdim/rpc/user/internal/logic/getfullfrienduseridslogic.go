package logic

import (
	"context"

	"fdim/pkg/authverify"
	"fdim/pkg/mcontext"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFullFriendUserIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFullFriendUserIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFullFriendUserIDsLogic {
	return &GetFullFriendUserIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFullFriendUserIDsLogic) GetFullFriendUserIDs(req *user.GetFullFriendUserIDsReq) (*user.GetFullFriendUserIDsResp, error) {
	resp := &user.GetFullFriendUserIDsResp{}

	// Ȩ����֤
	opUserID := mcontext.GetOpUserID(l.ctx)
	if req.UserID == "" {
		req.UserID = opUserID
	} else {
		if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
			return nil, err
		}
	}

	// TODO: �ӻ����ȡ���汾��
	// maxVersion, err := l.svcCtx.FriendCache.FindMaxFriendVersionCache(l.ctx, req.UserID)
	// if err != nil {
	// 	return nil, err
	// }

	// ��ȡ����ID�б�
	friendIDs, err := l.svcCtx.FriendDB.FindFriendUserIDs(l.ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// TODO: �������ID�б�Ĺ�ϣֵ
	// idHash := hashutil.IdHash(friendIDs)
	// if req.IdHash == idHash {
	// 	friendIDs = nil // �����ϣֵ��ͬ��˵������δ�仯�����ؿ��б�
	// }

	resp.UserIDs = friendIDs
	// TODO: ���ð汾��Ϣ
	// resp.Version = uint64(maxVersion.Version)
	// resp.VersionID = maxVersion.ID.Hex()
	// resp.Equal = req.IdHash == idHash

	return resp, nil
}
