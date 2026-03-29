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

	// ??????
	opUserID := mcontext.GetOpUserID(l.ctx)
	if req.UserID == "" {
		req.UserID = opUserID
	} else {
		if err := authverify.CheckAccess(l.ctx, req.UserID); err != nil {
			return nil, err
		}
	}

	// TODO: ????????????
	// maxVersion, err := l.svcCtx.FriendCache.FindMaxFriendVersionCache(l.ctx, req.UserID)
	// if err != nil {
	// 	return nil, err
	// }

	// ???????ID??
	friendIDs, err := l.svcCtx.FriendDB.FindFriendUserIDs(l.ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// TODO: ???????ID??????
	// idHash := hashutil.IdHash(friendIDs)
	// if req.IdHash == idHash {
	// 	friendIDs = nil // ?????????????????????????????
	// }

	resp.UserIDs = friendIDs
	// TODO: ???e???
	// resp.Version = uint64(maxVersion.Version)
	// resp.VersionID = maxVersion.ID.Hex()
	// resp.Equal = req.IdHash == idHash

	return resp, nil
}
