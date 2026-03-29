package server

import (
	"context"

	"fdim/protocol/chat"
	"fdim/rpc/chat/internal/logic"
	"fdim/rpc/chat/internal/svc"
)

type ChatServer struct {
	chat.UnimplementedChatServer
	svcCtx *svc.ServiceContext
}

func NewChatServer(svcCtx *svc.ServiceContext) *ChatServer {
	return &ChatServer{
		svcCtx: svcCtx,
	}
}

// Edit personal information - called by the user or an administrator
func (s *ChatServer) UpdateUserInfo(ctx context.Context, req *chat.UpdateUserInfoReq) (*chat.UpdateUserInfoResp, error) {
	l := logic.NewUpdateUserInfoLogic(ctx, s.svcCtx)
	return l.UpdateUserInfo(req)
}

func (s *ChatServer) AddUserAccount(ctx context.Context, req *chat.AddUserAccountReq) (*chat.AddUserAccountResp, error) {
	l := logic.NewAddUserAccountLogic(ctx, s.svcCtx)
	return l.AddUserAccount(req)
}

// Get user's public information - called by strangers
func (s *ChatServer) SearchUserPublicInfo(ctx context.Context, req *chat.SearchUserPublicInfoReq) (*chat.SearchUserPublicInfoResp, error) {
	l := logic.NewSearchUserPublicInfoLogic(ctx, s.svcCtx)
	return l.SearchUserPublicInfo(req)
}

func (s *ChatServer) FindUserPublicInfo(ctx context.Context, req *chat.FindUserPublicInfoReq) (*chat.FindUserPublicInfoResp, error) {
	l := logic.NewFindUserPublicInfoLogic(ctx, s.svcCtx)
	return l.FindUserPublicInfo(req)
}

// Search user information - called by administrators, other users get public fields
func (s *ChatServer) SearchUserFullInfo(ctx context.Context, req *chat.SearchUserFullInfoReq) (*chat.SearchUserFullInfoResp, error) {
	l := logic.NewSearchUserFullInfoLogic(ctx, s.svcCtx)
	return l.SearchUserFullInfo(req)
}

func (s *ChatServer) FindUserFullInfo(ctx context.Context, req *chat.FindUserFullInfoReq) (*chat.FindUserFullInfoResp, error) {
	l := logic.NewFindUserFullInfoLogic(ctx, s.svcCtx)
	return l.FindUserFullInfo(req)
}

func (s *ChatServer) SendVerifyCode(ctx context.Context, req *chat.SendVerifyCodeReq) (*chat.SendVerifyCodeResp, error) {
	l := logic.NewSendVerifyCodeLogic(ctx, s.svcCtx)
	return l.SendVerifyCode(req)
}

func (s *ChatServer) VerifyCode(ctx context.Context, req *chat.VerifyCodeReq) (*chat.VerifyCodeResp, error) {
	l := logic.NewVerifyCodeLogic(ctx, s.svcCtx)
	return l.VerifyCode(req)
}

func (s *ChatServer) RegisterUser(ctx context.Context, req *chat.RegisterUserReq) (*chat.RegisterUserResp, error) {
	l := logic.NewRegisterUserLogic(ctx, s.svcCtx)
	return l.RegisterUser(req)
}

func (s *ChatServer) Login(ctx context.Context, req *chat.LoginReq) (*chat.LoginResp, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(req)
}

func (s *ChatServer) ResetPassword(ctx context.Context, req *chat.ResetPasswordReq) (*chat.ResetPasswordResp, error) {
	l := logic.NewResetPasswordLogic(ctx, s.svcCtx)
	return l.ResetPassword(req)
}

func (s *ChatServer) ChangePassword(ctx context.Context, req *chat.ChangePasswordReq) (*chat.ChangePasswordResp, error) {
	l := logic.NewChangePasswordLogic(ctx, s.svcCtx)
	return l.ChangePassword(req)
}

func (s *ChatServer) CheckUserExist(ctx context.Context, req *chat.CheckUserExistReq) (*chat.CheckUserExistResp, error) {
	l := logic.NewCheckUserExistLogic(ctx, s.svcCtx)
	return l.CheckUserExist(req)
}

func (s *ChatServer) DelUserAccount(ctx context.Context, req *chat.DelUserAccountReq) (*chat.DelUserAccountResp, error) {
	l := logic.NewDelUserAccountLogic(ctx, s.svcCtx)
	return l.DelUserAccount(req)
}

func (s *ChatServer) FindUserAccount(ctx context.Context, req *chat.FindUserAccountReq) (*chat.FindUserAccountResp, error) {
	l := logic.NewFindUserAccountLogic(ctx, s.svcCtx)
	return l.FindUserAccount(req)
}

func (s *ChatServer) FindAccountUser(ctx context.Context, req *chat.FindAccountUserReq) (*chat.FindAccountUserResp, error) {
	l := logic.NewFindAccountUserLogic(ctx, s.svcCtx)
	return l.FindAccountUser(req)
}

func (s *ChatServer) FDIMCallback(ctx context.Context, req *chat.FDIMCallbackReq) (*chat.FDIMCallbackResp, error) {
	l := logic.NewFDIMCallbackLogic(ctx, s.svcCtx)
	return l.FDIMCallback(req)
}

// Statistics
func (s *ChatServer) UserLoginCount(ctx context.Context, req *chat.UserLoginCountReq) (*chat.UserLoginCountResp, error) {
	l := logic.NewUserLoginCountLogic(ctx, s.svcCtx)
	return l.UserLoginCount(req)
}

func (s *ChatServer) SearchUserInfo(ctx context.Context, req *chat.SearchUserInfoReq) (*chat.SearchUserInfoResp, error) {
	l := logic.NewSearchUserInfoLogic(ctx, s.svcCtx)
	return l.SearchUserInfo(req)
}

// Audio/video call and video meeting
func (s *ChatServer) GetTokenForVideoMeeting(ctx context.Context, req *chat.GetTokenForVideoMeetingReq) (*chat.GetTokenForVideoMeetingResp, error) {
	l := logic.NewGetTokenForVideoMeetingLogic(ctx, s.svcCtx)
	return l.GetTokenForVideoMeeting(req)
}

func (s *ChatServer) SetAllowRegister(ctx context.Context, req *chat.SetAllowRegisterReq) (*chat.SetAllowRegisterResp, error) {
	l := logic.NewSetAllowRegisterLogic(ctx, s.svcCtx)
	return l.SetAllowRegister(req)
}

func (s *ChatServer) GetAllowRegister(ctx context.Context, req *chat.GetAllowRegisterReq) (*chat.GetAllowRegisterResp, error) {
	l := logic.NewGetAllowRegisterLogic(ctx, s.svcCtx)
	return l.GetAllowRegister(req)
}
