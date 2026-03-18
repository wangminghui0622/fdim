package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/wrapperspb"
)

type SetGroupMemberInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetGroupMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupMemberInfoLogic {
	return &SetGroupMemberInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetGroupMemberInfoLogic) SetGroupMemberInfo(req *types.SetGroupMemberInfoReq) (resp *types.SetGroupMemberInfoResp, err error) {
	memberInfo := &user.SetGroupMemberInfo{
		GroupID: req.GroupID,
		UserID:  req.UserID,
	}

	if req.Info != nil {
		if info, ok := req.Info.(map[string]interface{}); ok {
			if nickname, ok := info["nickname"].(string); ok {
				memberInfo.Nickname = wrapperspb.String(nickname)
			}
			if faceURL, ok := info["faceURL"].(string); ok {
				memberInfo.FaceURL = wrapperspb.String(faceURL)
			}
			if roleLevel, ok := info["roleLevel"].(float64); ok {
				val := int32(roleLevel)
				memberInfo.RoleLevel = wrapperspb.Int32(val)
			}
			if ex, ok := info["ex"].(string); ok {
				memberInfo.Ex = wrapperspb.String(ex)
			}
		}
	}

	rpcReq := &user.SetGroupMemberInfoReq{
		Members: []*user.SetGroupMemberInfo{memberInfo},
	}

	_, err = l.svcCtx.GroupClient.SetGroupMemberInfo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.SetGroupMemberInfoResp{
	}

	return resp, nil
}
