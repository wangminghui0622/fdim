package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetIncrementalFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIncrementalFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalFriendsLogic {
	return &GetIncrementalFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIncrementalFriendsLogic) GetIncrementalFriends(req *types.GetIncrementalFriendsReq) (resp *types.GetIncrementalFriendsResp, err error) {
	rpcReq := &user.GetIncrementalFriendsReq{
		UserID:    req.UserID,
		VersionID: req.VersionID,
		Version:   req.Version,
	}

	rpcResp, err := l.svcCtx.FriendClient.GetIncrementalFriends(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var insert []types.FriendInfo
	for _, fi := range rpcResp.Insert {
		friendUserID := ""
		if fi.FriendUser != nil {
			friendUserID = fi.FriendUser.UserID
		}
		insert = append(insert, types.FriendInfo{
			OwnerUserID:  fi.OwnerUserID,
			FriendUserID: friendUserID,
			Remark:       fi.Remark,
			CreateTime:   fi.CreateTime,
			AddSource:    fi.AddSource,
			Ex:           fi.Ex,
		})
	}

	var update []types.FriendInfo
	for _, fi := range rpcResp.Update {
		friendUserID := ""
		if fi.FriendUser != nil {
			friendUserID = fi.FriendUser.UserID
		}
		update = append(update, types.FriendInfo{
			OwnerUserID:  fi.OwnerUserID,
			FriendUserID: friendUserID,
			Remark:       fi.Remark,
			CreateTime:   fi.CreateTime,
			AddSource:    fi.AddSource,
			Ex:           fi.Ex,
		})
	}

	resp = &types.GetIncrementalFriendsResp{
		Insert:    insert,
		Update:    update,
		Delete:    rpcResp.Delete,
		Version:   rpcResp.Version,
		VersionID: rpcResp.VersionID,
		Full:      rpcResp.Full,
	}

	return resp, nil
}
