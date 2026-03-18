package handler

import (
	"net/http"

	"fdim/api/internal/logic"
	"fdim/api/internal/svc"
	"fdim/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdminLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAdminLoginLogic(r.Context(), svcCtx)
		resp, err := l.AdminLogin(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AdminUpdateInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminUpdateInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAdminUpdateInfoLogic(r.Context(), svcCtx)
		resp, err := l.AdminUpdateInfo(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AdminInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAdminInfoLogic(r.Context(), svcCtx)
		resp, err := l.AdminInfo()
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func ChangeAdminPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangeAdminPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewChangeAdminPasswordLogic(r.Context(), svcCtx)
		resp, err := l.ChangeAdminPassword(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AddAdminAccountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddAdminAccountReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAddAdminAccountLogic(r.Context(), svcCtx)
		resp, err := l.AddAdminAccount(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DelAdminAccountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelAdminAccountReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewDelAdminAccountLogic(r.Context(), svcCtx)
		resp, err := l.DelAdminAccount(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchAdminAccountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchAdminAccountReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewSearchAdminAccountLogic(r.Context(), svcCtx)
		resp, err := l.SearchAdminAccount(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AddUserAccountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddUserAccountReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAddUserAccountLogic(r.Context(), svcCtx)
		resp, err := l.AddUserAccount(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func ImportUserByJsonHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImportUserByJsonReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewImportUserByJsonLogic(r.Context(), svcCtx)
		resp, err := l.ImportUserByJson(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GetAllowRegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetAllowRegisterLogic(r.Context(), svcCtx)
		resp, err := l.GetAllowRegister()
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SetAllowRegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetAllowRegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewSetAllowRegisterLogic(r.Context(), svcCtx)
		resp, err := l.SetAllowRegister(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AddDefaultFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddDefaultFriendReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAddDefaultFriendLogic(r.Context(), svcCtx)
		resp, err := l.AddDefaultFriend(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DelDefaultFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelDefaultFriendReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewDelDefaultFriendLogic(r.Context(), svcCtx)
		resp, err := l.DelDefaultFriend(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func FindDefaultFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewFindDefaultFriendLogic(r.Context(), svcCtx)
		resp, err := l.FindDefaultFriend()
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchDefaultFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchDefaultFriendReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewSearchDefaultFriendLogic(r.Context(), svcCtx)
		resp, err := l.SearchDefaultFriend(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AddDefaultGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddDefaultGroupReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAddDefaultGroupLogic(r.Context(), svcCtx)
		resp, err := l.AddDefaultGroup(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DelDefaultGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelDefaultGroupReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewDelDefaultGroupLogic(r.Context(), svcCtx)
		resp, err := l.DelDefaultGroup(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func FindDefaultGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewFindDefaultGroupLogic(r.Context(), svcCtx)
		resp, err := l.FindDefaultGroup()
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchDefaultGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchDefaultGroupReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewSearchDefaultGroupLogic(r.Context(), svcCtx)
		resp, err := l.SearchDefaultGroup(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AddInvitationCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddInvitationCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAddInvitationCodeLogic(r.Context(), svcCtx)
		resp, err := l.AddInvitationCode(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func GenInvitationCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenInvitationCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewGenInvitationCodeLogic(r.Context(), svcCtx)
		resp, err := l.GenInvitationCode(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DelInvitationCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelInvitationCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewDelInvitationCodeLogic(r.Context(), svcCtx)
		resp, err := l.DelInvitationCode(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchInvitationCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchInvitationCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewSearchInvitationCodeLogic(r.Context(), svcCtx)
		resp, err := l.SearchInvitationCode(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AddIPForbiddenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddIPForbiddenReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAddIPForbiddenLogic(r.Context(), svcCtx)
		resp, err := l.AddIPForbidden(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DelIPForbiddenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelIPForbiddenReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewDelIPForbiddenLogic(r.Context(), svcCtx)
		resp, err := l.DelIPForbidden(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchIPForbiddenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchIPForbiddenReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewSearchIPForbiddenLogic(r.Context(), svcCtx)
		resp, err := l.SearchIPForbidden(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AddUserIPLimitLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddUserIPLimitLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAddUserIPLimitLoginLogic(r.Context(), svcCtx)
		resp, err := l.AddUserIPLimitLogin(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DelUserIPLimitLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelUserIPLimitLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewDelUserIPLimitLoginLogic(r.Context(), svcCtx)
		resp, err := l.DelUserIPLimitLogin(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchUserIPLimitLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchUserIPLimitLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewSearchUserIPLimitLoginLogic(r.Context(), svcCtx)
		resp, err := l.SearchUserIPLimitLogin(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AddAppletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddAppletReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAddAppletLogic(r.Context(), svcCtx)
		resp, err := l.AddApplet(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DelAppletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelAppletReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewDelAppletLogic(r.Context(), svcCtx)
		resp, err := l.DelApplet(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func UpdateAppletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateAppletReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewUpdateAppletLogic(r.Context(), svcCtx)
		resp, err := l.UpdateApplet(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchAppletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchAppletReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewSearchAppletLogic(r.Context(), svcCtx)
		resp, err := l.SearchApplet(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func BlockUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BlockUserReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewBlockUserLogic(r.Context(), svcCtx)
		resp, err := l.BlockUser(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func UnblockUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnblockUserReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewUnblockUserLogic(r.Context(), svcCtx)
		resp, err := l.UnblockUser(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func SearchBlockUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchBlockUserReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewSearchBlockUserLogic(r.Context(), svcCtx)
		resp, err := l.SearchBlockUser(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AdminResetUserPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminResetUserPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAdminResetUserPasswordLogic(r.Context(), svcCtx)
		resp, err := l.AdminResetUserPassword(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AdminSetClientConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetClientConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAdminSetClientConfigLogic(r.Context(), svcCtx)
		resp, err := l.SetClientConfig(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AdminDelClientConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelClientConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAdminDelClientConfigLogic(r.Context(), svcCtx)
		resp, err := l.DelClientConfig(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func NewUserCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NewUserCountReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewNewUserCountLogic(r.Context(), svcCtx)
		resp, err := l.NewUserCount(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func LoginUserCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginUserCountReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewLoginUserCountLogic(r.Context(), svcCtx)
		resp, err := l.LoginUserCount(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func AddApplicationVersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddApplicationVersionReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewAddApplicationVersionLogic(r.Context(), svcCtx)
		resp, err := l.AddApplicationVersion(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func UpdateApplicationVersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateApplicationVersionReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewUpdateApplicationVersionLogic(r.Context(), svcCtx)
		resp, err := l.UpdateApplicationVersion(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}

func DeleteApplicationVersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteApplicationVersionReq
		if err := httpx.Parse(r, &req); err != nil {
			ParamError(w, err)
			return
		}
		l := logic.NewDeleteApplicationVersionLogic(r.Context(), svcCtx)
		resp, err := l.DeleteApplicationVersion(&req)
		if err != nil {
			ServerError(w, err)
		} else {
			Success(w, resp)
		}
	}
}
