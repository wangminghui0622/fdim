package logic

import (
	"context"
	"strconv"
	"strings"

	"fdim/pkg/errs"
	"fdim/protocol/chat"
	"fdim/protocol/wrapperspb"
	"fdim/rpc/chat/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *chat.UpdateUserInfoReq) (*chat.UpdateUserInfoResp, error) {
	// 1. 验证参数
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("user ID cannot be empty")
	}

	// 2. 检查用户是否存在
	_, err := l.svcCtx.ChatDB.GetUserAccountByUserID(l.ctx, req.UserID)
	if err != nil {
		return nil, errs.ErrArgs.WrapMsg("user not found")
	}

	// 3. 验证手机号和区号必须同时设置
	if (req.PhoneNumber != nil && req.PhoneNumber.Value != "") || (req.AreaCode != nil && req.AreaCode.Value != "") {
		if req.PhoneNumber == nil || req.AreaCode == nil || req.PhoneNumber.Value == "" || req.AreaCode.Value == "" {
			return nil, errs.ErrArgs.WrapMsg("area code and phone number must be set together")
		}
		// 验证手机号格式
		if !strings.HasPrefix(req.AreaCode.Value, "+") {
			areaCodeValue := "+" + req.AreaCode.Value
			req.AreaCode = &wrapperspb.StringValue{Value: areaCodeValue}
		}
		if _, err := strconv.ParseUint(req.AreaCode.Value[1:], 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("area code must be number")
		}
		if _, err := strconv.ParseUint(req.PhoneNumber.Value, 10, 64); err != nil {
			return nil, errs.ErrArgs.WrapMsg("phone number must be number")
		}
	}

	// 4. 构建更新字段
	update := make(map[string]interface{})
	if req.Account != nil && req.Account.Value != "" {
		update["account"] = req.Account.Value
	}
	if req.PhoneNumber != nil && req.PhoneNumber.Value != "" {
		update["phone_number"] = req.PhoneNumber.Value
	}
	if req.AreaCode != nil && req.AreaCode.Value != "" {
		update["area_code"] = req.AreaCode.Value
	}
	if req.Email != nil && req.Email.Value != "" {
		update["email"] = req.Email.Value
	}
	if req.Nickname != nil {
		update["nickname"] = req.Nickname.Value
	}
	if req.FaceURL != nil {
		update["face_url"] = req.FaceURL.Value
	}
	if req.Gender != nil {
		update["gender"] = req.Gender.Value
	}
	if req.Level != nil {
		update["level"] = req.Level.Value
	}
	if req.Birth != nil {
		update["birth"] = req.Birth.Value
	}
	if req.AllowAddFriend != nil {
		update["allow_add_friend"] = req.AllowAddFriend.Value
	}
	if req.AllowBeep != nil {
		update["allow_beep"] = req.AllowBeep.Value
	}
	if req.AllowVibration != nil {
		update["allow_vibration"] = req.AllowVibration.Value
	}
	if req.GlobalRecvMsgOpt != nil {
		update["global_recv_msg_opt"] = req.GlobalRecvMsgOpt.Value
	}
	if req.RegisterType != nil {
		update["register_type"] = req.RegisterType.Value
	}

	// 5. 更新用户信息
	if len(update) > 0 {
		if err := l.svcCtx.ChatDB.UpdateUserInfo(l.ctx, req.UserID, update); err != nil {
			l.Errorf("UpdateUserInfo failed: %v", err)
			return nil, errs.WrapMsg(err, "failed to update user info")
		}
	}

	// 6. 更新用户账户信息（如果需要）
	accountUpdate := make(map[string]interface{})
	if req.Account != nil && req.Account.Value != "" {
		accountUpdate["account"] = req.Account.Value
	}
	if req.PhoneNumber != nil && req.PhoneNumber.Value != "" {
		accountUpdate["phone_number"] = req.PhoneNumber.Value
	}
	if req.AreaCode != nil && req.AreaCode.Value != "" {
		accountUpdate["area_code"] = req.AreaCode.Value
	}
	if req.Email != nil && req.Email.Value != "" {
		accountUpdate["email"] = req.Email.Value
	}
	if len(accountUpdate) > 0 {
		if err := l.svcCtx.ChatDB.UpdateUserAccount(l.ctx, req.UserID, accountUpdate); err != nil {
			l.Errorw("UpdateUserAccount failed", logx.Field("error", err))
		}
	}

	// 7. 获取更新后的信息用于返回
	userInfos, err := l.svcCtx.ChatDB.FindUserFullInfo(l.ctx, []string{req.UserID})
	if err != nil || len(userInfos) == 0 {
		return &chat.UpdateUserInfoResp{}, nil
	}

	userInfo := userInfos[0]
	return &chat.UpdateUserInfoResp{
		FaceUrl:  userInfo.FaceURL,
		NickName: userInfo.Nickname,
	}, nil
}
