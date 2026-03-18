package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
)

type CreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGroupLogic) CreateGroup(req *types.CreateGroupReq) (resp *types.CreateGroupResp, err error) {
	rpcReq := &user.CreateGroupReq{
		MemberUserIDs: req.MemberUserIDs,
		AdminUserIDs:  req.AdminUserIDs,
		OwnerUserID:   req.OwnerUserID,
	}

	if req.SendMessage != nil {
		rpcReq.SendMessage = req.SendMessage
	}

	if req.GroupInfo != nil {
		if gi, ok := req.GroupInfo.(map[string]interface{}); ok {
			groupInfo := &sdkws.GroupInfo{}
			if groupID, ok := gi["groupID"].(string); ok {
				groupInfo.GroupID = groupID
			}
			if groupName, ok := gi["groupName"].(string); ok {
				groupInfo.GroupName = groupName
			}
			if faceURL, ok := gi["faceURL"].(string); ok {
				groupInfo.FaceURL = faceURL
			}
			if introduction, ok := gi["introduction"].(string); ok {
				groupInfo.Introduction = introduction
			}
			if notification, ok := gi["notification"].(string); ok {
				groupInfo.Notification = notification
			}
			if ex, ok := gi["ex"].(string); ok {
				groupInfo.Ex = ex
			}
			rpcReq.GroupInfo = groupInfo
		}
	}

	rpcResp, err := l.svcCtx.GroupClient.CreateGroup(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.CreateGroupResp{
		GroupInfo: rpcResp.GroupInfo,
	}

	return resp, nil
}
