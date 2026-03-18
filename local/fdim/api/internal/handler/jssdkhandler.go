package handler

import (
	"net/http"

	"fdim/api/internal/logic"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func JSSdkGetConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JSSdkGetConversationsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewJSSdkGetConversationsLogic(r.Context(), svcCtx)
		resp, err := l.JSSdkGetConversations(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func JSSdkGetActiveConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JSSdkGetActiveConversationsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewJSSdkGetActiveConversationsLogic(r.Context(), svcCtx)
		resp, err := l.JSSdkGetActiveConversations(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}
