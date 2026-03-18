package handler

import (
	"net/http"

	"fdim/api/internal/logic"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FindUserFullInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindUserFullInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewFindUserFullInfoLogic(r.Context(), svcCtx)
		resp, err := l.FindUserFullInfo(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func FindUserPublicInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FindUserPublicInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewFindUserPublicInfoLogic(r.Context(), svcCtx)
		resp, err := l.FindUserPublicInfo(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchUserPublicInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchUserPublicInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewSearchUserPublicInfoLogic(r.Context(), svcCtx)
		resp, err := l.SearchUserPublicInfo(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetTokenForVideoMeetingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTokenForVideoMeetingReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetTokenForVideoMeetingLogic(r.Context(), svcCtx)
		resp, err := l.GetTokenForVideoMeeting(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func ChatUpdateUserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatUpdateUserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewChatUpdateUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.ChatUpdateUserInfo(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetClientConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetClientConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetClientConfigLogic(r.Context(), svcCtx)
		resp, err := l.GetClientConfig(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func FindAppletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewFindAppletLogic(r.Context(), svcCtx)
		resp, err := l.FindApplet()
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func LatestApplicationVersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LatestApplicationVersionReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewLatestApplicationVersionLogic(r.Context(), svcCtx)
		resp, err := l.LatestApplicationVersion(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func PageApplicationVersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageApplicationVersionReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewPageApplicationVersionLogic(r.Context(), svcCtx)
		resp, err := l.PageApplicationVersion(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func OpenIMCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenIMCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewOpenIMCallbackLogic(r.Context(), svcCtx)
		resp, err := l.OpenIMCallback(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}
