package handler

import (
	"net/http"

	"fdim/api/internal/logic"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Auth API Handlers

func GetAdminTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAdminTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetAdminTokenLogic(r.Context(), svcCtx)
		resp, err := l.GetAdminToken(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetUserTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetUserTokenLogic(r.Context(), svcCtx)
		resp, err := l.GetUserToken(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func ParseTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ParseTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewParseTokenLogic(r.Context(), svcCtx)
		resp, err := l.ParseToken(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func ForceLogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ForceLogoutReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewForceLogoutLogic(r.Context(), svcCtx)
		resp, err := l.ForceLogout(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}
