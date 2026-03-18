package logic

import (
	"context"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/protocol/sdkws"
)

type SearchUserFullInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchUserFullInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserFullInfoLogic {
	return &SearchUserFullInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchUserFullInfoLogic) SearchUserFullInfo(req *types.SearchUserFullInfoReq) (*types.SearchUserFullInfoResp, error) {
	if l.svcCtx.ChatClient == nil {
		return nil, errs.ErrInternalServer.WrapMsg("chat service not available")
	}

	// 调用 Chat RPC SearchUserFullInfo（Chat层存储了手机号等注册信息）
	rpcResp, err := l.svcCtx.ChatClient.SearchUserFullInfo(l.ctx, &chat.SearchUserFullInfoReq{
		Keyword: req.Keyword,
		Pagination: &sdkws.RequestPagination{
			PageNumber: req.Pagination.PageNumber,
			ShowNumber: req.Pagination.ShowNumber,
		},
	})
	if err != nil {
		return nil, err
	}

	// 转换响应
	var result []types.UserFullInfo
	for _, u := range rpcResp.Users {
		result = append(result, types.UserFullInfo{
			UserID:           u.UserID,
			Nickname:         u.Nickname,
			FaceURL:          u.FaceURL,
			PhoneNumber:      u.PhoneNumber,
			Gender:           u.Gender,
			GlobalRecvMsgOpt: u.GlobalRecvMsgOpt,
		})
	}

	return &types.SearchUserFullInfoResp{
		Total: int32(rpcResp.Total),
		Users: result,
	}, nil
}
