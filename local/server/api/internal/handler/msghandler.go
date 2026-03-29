package handler

import (
	"net/http"

	"fdim/api/internal/logic"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Msg API Handlers

func GetSeqHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSeqReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetSeqLogic(r.Context(), svcCtx)
		resp, err := l.GetSeq(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchMsgReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewSearchMsgLogic(r.Context(), svcCtx)
		resp, err := l.SearchMsg(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SendMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendMsgReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewSendMessageLogic(r.Context(), svcCtx)
		resp, err := l.SendMessage(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SendBusinessNotificationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendBusinessNotificationReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewSendBusinessNotificationLogic(r.Context(), svcCtx)
		resp, err := l.SendBusinessNotification(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func PullMsgBySeqsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PullMsgBySeqsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewPullMsgBySeqsLogic(r.Context(), svcCtx)
		resp, err := l.PullMsgBySeqs(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func RevokeMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RevokeMsgReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewRevokeMsgLogic(r.Context(), svcCtx)
		resp, err := l.RevokeMsg(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func MarkMsgsAsReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MarkMsgsAsReadReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewMarkMsgsAsReadLogic(r.Context(), svcCtx)
		resp, err := l.MarkMsgsAsRead(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func MarkConversationAsReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MarkConversationAsReadReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewMarkConversationAsReadLogic(r.Context(), svcCtx)
		resp, err := l.MarkConversationAsRead(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetConversationsHasReadAndMaxSeqHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetConversationsHasReadAndMaxSeqReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetConversationsHasReadAndMaxSeqLogic(r.Context(), svcCtx)
		resp, err := l.GetConversationsHasReadAndMaxSeq(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SetConversationHasReadSeqHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetConversationHasReadSeqReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewSetConversationHasReadSeqLogic(r.Context(), svcCtx)
		resp, err := l.SetConversationHasReadSeq(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func ClearConversationsMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ClearConversationsMsgReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewClearConversationsMsgLogic(r.Context(), svcCtx)
		resp, err := l.ClearConversationsMsg(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func UserClearAllMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserClearAllMsgReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewUserClearAllMsgLogic(r.Context(), svcCtx)
		resp, err := l.UserClearAllMsg(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DeleteMsgsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteMsgsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewDeleteMsgsLogic(r.Context(), svcCtx)
		resp, err := l.DeleteMsgs(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DeleteMsgPhysicalBySeqHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteMsgPhysicalBySeqReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewDeleteMsgPhysicalBySeqLogic(r.Context(), svcCtx)
		resp, err := l.DeleteMsgPhysicalBySeq(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DeleteMsgPhysicalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteMsgPhysicalReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewDeleteMsgPhysicalLogic(r.Context(), svcCtx)
		resp, err := l.DeleteMsgPhysical(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func BatchSendMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchSendMsgReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewBatchSendMsgLogic(r.Context(), svcCtx)
		resp, err := l.BatchSendMsg(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SendSimpleMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendSimpleMsgReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewSendSimpleMessageLogic(r.Context(), svcCtx)
		resp, err := l.SendSimpleMessage(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func CheckMsgIsSendSuccessHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckMsgIsSendSuccessReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewCheckMsgIsSendSuccessLogic(r.Context(), svcCtx)
		resp, err := l.CheckMsgIsSendSuccess(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetServerTimeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetServerTimeReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetServerTimeLogic(r.Context(), svcCtx)
		resp, err := l.GetServerTime(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}
