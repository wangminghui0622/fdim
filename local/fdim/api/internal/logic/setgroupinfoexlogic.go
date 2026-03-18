package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/wrapperspb"
)

type SetGroupInfoExLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetGroupInfoExLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupInfoExLogic {
	return &SetGroupInfoExLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetGroupInfoExLogic) SetGroupInfoEx(req *types.SetGroupInfoExReq) (resp *types.SetGroupInfoExResp, err error) {
	rpcReq := &user.SetGroupInfoExReq{
		GroupID: req.GroupID,
	}

	if req.GroupName != nil {
		rpcReq.GroupName = wrapperspb.String(*req.GroupName)
	}
	if req.Notification != nil {
		rpcReq.Notification = wrapperspb.String(*req.Notification)
	}
	if req.Introduction != nil {
		rpcReq.Introduction = wrapperspb.String(*req.Introduction)
	}
	if req.FaceURL != nil {
		rpcReq.FaceURL = wrapperspb.String(*req.FaceURL)
	}
	if req.Ex != nil {
		rpcReq.Ex = wrapperspb.String(*req.Ex)
	}
	if req.NeedVerification != nil {
		rpcReq.NeedVerification = wrapperspb.Int32(*req.NeedVerification)
	}
	if req.LookMemberInfo != nil {
		rpcReq.LookMemberInfo = wrapperspb.Int32(*req.LookMemberInfo)
	}
	if req.ApplyMemberFriend != nil {
		rpcReq.ApplyMemberFriend = wrapperspb.Int32(*req.ApplyMemberFriend)
	}

	_, err = l.svcCtx.GroupClient.SetGroupInfoEx(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.SetGroupInfoExResp{
	}

	return resp, nil
}
