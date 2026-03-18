package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type SetFriendRemarkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetFriendRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetFriendRemarkLogic {
	return &SetFriendRemarkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetFriendRemarkLogic) SetFriendRemark(req *types.SetFriendRemarkReq) (resp *types.SetFriendRemarkResp, err error) {
	rpcReq := &user.SetFriendRemarkReq{
		OwnerUserID:  req.OwnerUserID,
		FriendUserID: req.FriendUserID,
		Remark:       req.Remark,
	}

	_, err = l.svcCtx.FriendClient.SetFriendRemark(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.SetFriendRemarkResp{
	}

	return resp, nil
}
