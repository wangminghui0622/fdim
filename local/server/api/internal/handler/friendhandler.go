package handler

import (
	"net/http"

	"fdim/api/internal/logic"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Friend API Handlers

func ApplyToAddFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApplyToAddFriendReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewApplyToAddFriendLogic(r.Context(), svcCtx)
		resp, err := l.ApplyToAddFriend(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func RespondFriendApplyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RespondFriendApplyReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewRespondFriendApplyLogic(r.Context(), svcCtx)
		resp, err := l.RespondFriendApply(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DeleteFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteFriendReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewDeleteFriendLogic(r.Context(), svcCtx)
		resp, err := l.DeleteFriend(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetFriendApplyListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFriendApplyListReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetFriendApplyListLogic(r.Context(), svcCtx)
		resp, err := l.GetFriendApplyList(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetDesignatedFriendsApplyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDesignatedFriendsApplyReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetDesignatedFriendsApplyLogic(r.Context(), svcCtx)
		resp, err := l.GetDesignatedFriendsApply(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetSelfApplyListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSelfApplyListReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetSelfApplyListLogic(r.Context(), svcCtx)
		resp, err := l.GetSelfApplyList(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetFriendListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFriendListReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetFriendListLogic(r.Context(), svcCtx)
		resp, err := l.GetFriendList(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetDesignatedFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDesignatedFriendsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetDesignatedFriendsLogic(r.Context(), svcCtx)
		resp, err := l.GetDesignatedFriends(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SetFriendRemarkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetFriendRemarkReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewSetFriendRemarkLogic(r.Context(), svcCtx)
		resp, err := l.SetFriendRemark(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AddBlackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddBlackReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewAddBlackLogic(r.Context(), svcCtx)
		resp, err := l.AddBlack(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetPaginationBlacksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPaginationBlacksReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetPaginationBlacksLogic(r.Context(), svcCtx)
		resp, err := l.GetPaginationBlacks(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetSpecifiedBlacksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSpecifiedBlacksReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetSpecifiedBlacksLogic(r.Context(), svcCtx)
		resp, err := l.GetSpecifiedBlacks(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func RemoveBlackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveBlackReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewRemoveBlackLogic(r.Context(), svcCtx)
		resp, err := l.RemoveBlack(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetIncrementalBlacksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetIncrementalBlacksReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetIncrementalBlacksLogic(r.Context(), svcCtx)
		resp, err := l.GetIncrementalBlacks(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func ImportFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImportFriendsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewImportFriendsLogic(r.Context(), svcCtx)
		resp, err := l.ImportFriends(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func IsFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IsFriendReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewIsFriendLogic(r.Context(), svcCtx)
		resp, err := l.IsFriend(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetFriendIDsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFriendIDsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetFriendIDsLogic(r.Context(), svcCtx)
		resp, err := l.GetFriendIDs(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetSpecifiedFriendsInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSpecifiedFriendsInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetSpecifiedFriendsInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetSpecifiedFriendsInfo(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func UpdateFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateFriendsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewUpdateFriendsLogic(r.Context(), svcCtx)
		resp, err := l.UpdateFriends(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetIncrementalFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetIncrementalFriendsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetIncrementalFriendsLogic(r.Context(), svcCtx)
		resp, err := l.GetIncrementalFriends(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetFullFriendUserIDsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFullFriendUserIDsReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetFullFriendUserIDsLogic(r.Context(), svcCtx)
		resp, err := l.GetFullFriendUserIDs(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetSelfUnhandledApplyCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSelfUnhandledApplyCountReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewGetSelfUnhandledApplyCountLogic(r.Context(), svcCtx)
		resp, err := l.GetSelfUnhandledApplyCount(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchFriendReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}

		l := logic.NewSearchFriendLogic(r.Context(), svcCtx)
		resp, err := l.SearchFriend(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}
