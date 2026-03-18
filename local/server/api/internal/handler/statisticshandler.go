package handler

import (
	"net/http"

	"fdim/api/internal/logic"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetActiveUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetActiveUserReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewGetActiveUserLogic(r.Context(), svcCtx)
		resp, err := l.GetActiveUser(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GroupCreateCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupCreateCountReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewGroupCreateCountLogic(r.Context(), svcCtx)
		resp, err := l.GroupCreateCount(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetActiveGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetActiveGroupReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewGetActiveGroupLogic(r.Context(), svcCtx)
		resp, err := l.GetActiveGroup(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}
