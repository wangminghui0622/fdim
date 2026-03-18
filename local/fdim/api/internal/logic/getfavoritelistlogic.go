package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
)

type GetFavoriteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteListLogic {
	return &GetFavoriteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFavoriteListLogic) GetFavoriteList(req *types.GetFavoriteListReq, userID string) (resp *types.GetFavoriteListResp, err error) {
	favorites, total, err := l.svcCtx.FavoriteDB.GetByUserID(userID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}

	var result []types.FavoriteMsg
	for _, f := range favorites {
		result = append(result, types.FavoriteMsg{
			FavoriteID:     f.ID,
			ClientMsgID:    f.ClientMsgID,
			ServerMsgID:    f.ServerMsgID,
			ConversationID: f.ConversationID,
			SendID:         f.SendID,
			RecvID:         f.RecvID,
			GroupID:        f.GroupID,
			SenderNickname: f.SenderNickname,
			SenderFaceURL:  f.SenderFaceURL,
			SessionType:    f.SessionType,
			ContentType:    f.ContentType,
			Content:        f.Content,
			SendTime:       f.SendTime,
			CreateTime:     f.CreateTime.UnixMilli(),
			Ex:             f.Ex,
		})
	}

	resp = &types.GetFavoriteListResp{
		Total:     total,
		Favorites: result,
	}

	return resp, nil
}
