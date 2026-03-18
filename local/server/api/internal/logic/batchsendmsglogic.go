package logic

import (
	"context"
	"encoding/json"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/msg"
	"fdim/protocol/sdkws"
)

type BatchSendMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchSendMsgLogic {
	return &BatchSendMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchSendMsgLogic) BatchSendMsg(req *types.BatchSendMsgReq) (resp *types.BatchSendMsgResp, err error) {
	// 验证 SendMsg 字段
	if req.SendMsg == nil {
		return nil, errs.ErrArgs.WrapMsg("field \"sendMsg\" is required")
	}

	// 验证必需字段
	if req.SendMsg.ClientMsgID == "" {
		return nil, errs.ErrArgs.WrapMsg("clientMsgID is required in sendMsg")
	}
	if req.SendMsg.SendID == "" {
		return nil, errs.ErrArgs.WrapMsg("sendID is required in sendMsg")
	}

	recvIDs := req.RecvIDs
	if req.IsSendAll {
		recvIDs = []string{}
	}

	var results []interface{}
	var failedIDs []string

	for _, recvID := range recvIDs {
		msgData := &sdkws.MsgData{
			SendID:           req.SendMsg.SendID,
			RecvID:           recvID,
			GroupID:          req.SendMsg.GroupID,
			ClientMsgID:      req.SendMsg.ClientMsgID,
			ServerMsgID:      req.SendMsg.ServerMsgID,
			SenderPlatformID: req.SendMsg.SenderPlatformID,
			SenderNickname:   req.SendMsg.SenderNickname,
			SenderFaceURL:    req.SendMsg.SenderFaceURL,
			SessionType:      req.SendMsg.SessionType,
			MsgFrom:          req.SendMsg.MsgFrom,
			ContentType:      req.SendMsg.ContentType,
			Content:          []byte(req.SendMsg.Content),
			Seq:              req.SendMsg.Seq,
			SendTime:         req.SendMsg.SendTime,
			CreateTime:       req.SendMsg.CreateTime,
			Status:           req.SendMsg.Status,
			IsRead:           req.SendMsg.IsRead,
			Options:          req.SendMsg.Options,
			AtUserIDList:     req.SendMsg.AtUserIDList,
			AttachedInfo:     req.SendMsg.AttachedInfo,
			Ex:               req.SendMsg.Ex,
		}

		// 处理 OfflinePushInfo
		if req.SendMsg.OfflinePushInfo != nil {
			jsonData, err := json.Marshal(req.SendMsg.OfflinePushInfo)
			if err == nil {
				var offlinePush sdkws.OfflinePushInfo
				if err := json.Unmarshal(jsonData, &offlinePush); err == nil {
					msgData.OfflinePushInfo = &offlinePush
				}
			}
		}

		rpcReq := &msg.SendMsgReq{
			MsgData: msgData,
		}

		rpcResp, err := l.svcCtx.MsgClient.SendMsg(l.ctx, rpcReq)
		if err != nil {
			failedIDs = append(failedIDs, recvID)
			continue
		}

		results = append(results, map[string]interface{}{
			"serverMsgID": rpcResp.ServerMsgID,
			"clientMsgID": rpcResp.ClientMsgID,
			"sendTime":    rpcResp.SendTime,
			"recvID":      recvID,
		})
	}

	// 转换响应
	resp = &types.BatchSendMsgResp{
		Results:   results,
		FailedIDs: failedIDs,
	}

	return resp, nil
}
