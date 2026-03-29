package handler

import (
	"net/http"

	"fdim/api/internal/logic"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Third API Handlers

func GetPrometheusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetPrometheusLogic(r.Context(), svcCtx)
		resp, err := l.GetPrometheus()
		if err != nil {
			ServerError(w, err)
		} else {
			// GET 请求，重定向?Grafana URL
			http.Redirect(w, r, resp.Url, http.StatusFound)
		}
	}
}

func FcmUpdateTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FcmUpdateTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewFcmUpdateTokenLogic(r.Context(), svcCtx)
		resp, err := l.FcmUpdateToken(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SetAppBadgeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetAppBadgeReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewSetAppBadgeLogic(r.Context(), svcCtx)
		resp, err := l.SetAppBadge(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func UploadLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadLogsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewUploadLogsLogic(r.Context(), svcCtx)
		resp, err := l.UploadLogs(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DeleteLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteLogsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewDeleteLogsLogic(r.Context(), svcCtx)
		resp, err := l.DeleteLogs(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchLogsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewSearchLogsLogic(r.Context(), svcCtx)
		resp, err := l.SearchLogs(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func PartLimitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PartLimitReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewPartLimitLogic(r.Context(), svcCtx)
		resp, err := l.PartLimit(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func PartSizeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PartSizeReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewPartSizeLogic(r.Context(), svcCtx)
		resp, err := l.PartSize(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func InitiateMultipartUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InitiateMultipartUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewInitiateMultipartUploadLogic(r.Context(), svcCtx)
		resp, err := l.InitiateMultipartUpload(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AuthSignHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthSignReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewAuthSignLogic(r.Context(), svcCtx)
		resp, err := l.AuthSign(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func CompleteMultipartUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CompleteMultipartUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewCompleteMultipartUploadLogic(r.Context(), svcCtx)
		resp, err := l.CompleteMultipartUpload(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AccessURLHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AccessURLReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewAccessURLLogic(r.Context(), svcCtx)
		resp, err := l.AccessURL(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func InitiateFormDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InitiateFormDataReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewInitiateFormDataLogic(r.Context(), svcCtx)
		resp, err := l.InitiateFormData(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func CompleteFormDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CompleteFormDataReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewCompleteFormDataLogic(r.Context(), svcCtx)
		resp, err := l.CompleteFormData(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func ObjectRedirectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从路径参数获?name
		name := r.URL.Path
		if name != "" && name[0] == '/' {
			name = name[1:]
		}

		// 从查询参数构?query map
		query := make(map[string]string)
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				query[key] = values[0]
			}
		}

		req := &types.AccessURLReq{
			Name:  name,
			Query: query,
		}

		l := logic.NewObjectRedirectLogic(r.Context(), svcCtx)
		resp, err := l.ObjectRedirect(req)
		if err != nil {
			ServerError(w, err)
		} else {
			http.Redirect(w, r, resp.Url, http.StatusFound)
		}
	}
}
