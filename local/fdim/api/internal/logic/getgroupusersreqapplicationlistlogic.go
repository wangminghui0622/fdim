package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetGroupUsersReqApplicationListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupUsersReqApplicationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupUsersReqApplicationListLogic {
	return &GetGroupUsersReqApplicationListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupUsersReqApplicationListLogic) GetGroupUsersReqApplicationList(req *types.GetGroupUsersReqApplicationListReq) (resp *types.GetGroupUsersReqApplicationListResp, err error) {
	rpcReq := &user.GetGroupUsersReqApplicationListReq{
		GroupID: req.GroupID,
		UserIDs: req.UserIDs,
	}

	rpcResp, err := l.svcCtx.GroupClient.GetGroupUsersReqApplicationList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	var groupRequests []interface{}
	for _, gr := range rpcResp.GroupRequests {
		groupRequests = append(groupRequests, gr)
	}

	resp = &types.GetGroupUsersReqApplicationListResp{
		Total:         uint32(rpcResp.Total),
		GroupRequests: groupRequests,
	}

	return resp, nil
}
