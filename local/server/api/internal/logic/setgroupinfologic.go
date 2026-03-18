package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
	"fdim/protocol/sdkws"
	"fdim/protocol/wrapperspb"
)

type SetGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupInfoLogic {
	return &SetGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetGroupInfoLogic) SetGroupInfo(req *types.SetGroupInfoReq) (resp *types.SetGroupInfoResp, err error) {
	rpcReq := &user.SetGroupInfoReq{}

	if req.GroupInfoForSet != nil {
		if gifs, ok := req.GroupInfoForSet.(map[string]interface{}); ok {
			groupInfoForSet := &sdkws.GroupInfoForSet{}
			if groupID, ok := gifs["groupID"].(string); ok {
				groupInfoForSet.GroupID = groupID
			}
			if groupName, ok := gifs["groupName"].(string); ok {
				groupInfoForSet.GroupName = groupName
			}
			if faceURL, ok := gifs["faceURL"].(string); ok {
				groupInfoForSet.FaceURL = faceURL
			}
			if introduction, ok := gifs["introduction"].(string); ok {
				groupInfoForSet.Introduction = introduction
			}
			if notification, ok := gifs["notification"].(string); ok {
				groupInfoForSet.Notification = notification
			}
			if ex, ok := gifs["ex"].(string); ok {
				groupInfoForSet.Ex = wrapperspb.String(ex)
			}
			if needVerification, ok := gifs["needVerification"].(float64); ok {
				val := int32(needVerification)
				groupInfoForSet.NeedVerification = wrapperspb.Int32(val)
			}
			if lookMemberInfo, ok := gifs["lookMemberInfo"].(float64); ok {
				val := int32(lookMemberInfo)
				groupInfoForSet.LookMemberInfo = wrapperspb.Int32(val)
			}
			if applyMemberFriend, ok := gifs["applyMemberFriend"].(float64); ok {
				val := int32(applyMemberFriend)
				groupInfoForSet.ApplyMemberFriend = wrapperspb.Int32(val)
			}
			rpcReq.GroupInfoForSet = groupInfoForSet
		}
	}

	_, err = l.svcCtx.GroupClient.SetGroupInfo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &types.SetGroupInfoResp{
	}

	return resp, nil
}
