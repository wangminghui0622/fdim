package handler

import (
	"net/http"

	"fdim/api/internal/logic"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"fdim/pkg/util/ctxutil"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddFavoriteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddFavoriteReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		userID := ctxutil.GetUserID(r.Context())
		l := logic.NewAddFavoriteLogic(r.Context(), svcCtx)
		resp, err := l.AddFavorite(&req, userID)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DeleteFavoriteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteFavoriteReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		userID := ctxutil.GetUserID(r.Context())
		l := logic.NewDeleteFavoriteLogic(r.Context(), svcCtx)
		resp, err := l.DeleteFavorite(&req, userID)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetFavoriteListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFavoriteListReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		userID := ctxutil.GetUserID(r.Context())
		l := logic.NewGetFavoriteListLogic(r.Context(), svcCtx)
		resp, err := l.GetFavoriteList(&req, userID)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}
