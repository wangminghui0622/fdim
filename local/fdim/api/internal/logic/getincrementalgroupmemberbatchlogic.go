package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/protocol/user"
)

type GetIncrementalGroupMemberBatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIncrementalGroupMemberBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIncrementalGroupMemberBatchLogic {
	return &GetIncrementalGroupMemberBatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIncrementalGroupMemberBatchLogic) GetIncrementalGroupMemberBatch(req *types.GetIncrementalGroupMemberBatchReq) (resp *types.GetIncrementalGroupMemberBatchResp, err error) {
	var reqList []*user.GetIncrementalGroupMemberReq
	for _, r := range req.ReqList {
		reqList = append(reqList, &user.GetIncrementalGroupMemberReq{
			GroupID:   r.GroupID,
			VersionID: r.VersionID,
			Version:   r.Version,
		})
	}

	rpcReq := &user.BatchGetIncrementalGroupMemberReq{
		ReqList: reqList,
	}

	rpcResp, err := l.svcCtx.GroupClient.BatchGetIncrementalGroupMember(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	respList := make(map[string]types.GetIncrementalGroupMemberResp)
	for k, v := range rpcResp.RespList {
		var insert []interface{}
		for _, m := range v.Insert {
			insert = append(insert, m)
		}
		var update []interface{}
		for _, m := range v.Update {
			update = append(update, m)
		}
		respList[k] = types.GetIncrementalGroupMemberResp{
			Version:   v.Version,
			VersionID: v.VersionID,
			Full:      v.Full,
			Delete:    v.Delete,
			Insert:    insert,
			Update:    update,
		}
	}

	resp = &types.GetIncrementalGroupMemberBatchResp{
		RespList: respList,
	}

	return resp, nil
}
