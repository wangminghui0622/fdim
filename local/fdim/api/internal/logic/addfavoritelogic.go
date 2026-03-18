package logic

import (
	"context"

	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/model"
)

type AddFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFavoriteLogic {
	return &AddFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFavoriteLogic) AddFavorite(req *types.AddFavoriteReq, userID string) (resp *types.AddFavoriteResp, err error) {
	favorite := &model.FavoriteMsg{
		UserID:         userID,
		ClientMsgID:    req.ClientMsgID,
		ServerMsgID:    req.ServerMsgID,
		ConversationID: req.ConversationID,
		SendID:         req.SendID,
		RecvID:         req.RecvID,
		GroupID:        req.GroupID,
		SenderNickname: req.SenderNickname,
		SenderFaceURL:  req.SenderFaceURL,
		SessionType:    req.SessionType,
		ContentType:    req.ContentType,
		Content:        req.Content,
		SendTime:       req.SendTime,
		Ex:             req.Ex,
	}

	err = l.svcCtx.FavoriteDB.Create(favorite)
	if err != nil {
		return nil, err
	}

	resp = &types.AddFavoriteResp{
		FavoriteID: favorite.ID,
	}

	return resp, nil
}
