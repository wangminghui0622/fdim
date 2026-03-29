package logic

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"fdim/pkg/authverify"
	"fdim/pkg/constant"
	"fdim/pkg/model"
	"fdim/protocol/user"
	"fdim/rpc/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddNotificationAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddNotificationAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddNotificationAccountLogic {
	return &AddNotificationAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddNotificationAccountLogic) AddNotificationAccount(req *user.AddNotificationAccountReq) (*user.AddNotificationAccountResp, error) {
	// 权限验证：需要管理员权限
	if err := authverify.CheckAdmin(l.ctx); err != nil {
		return nil, err
	}

	// 检?appMangerLevel
	if req.AppMangerLevel < constant.AppNotificationAdmin {
		return nil, fmt.Errorf("app level not supported")
	}

	// 如果 UserID 为空，生成一?
	if req.UserID == "" {
		for i := 0; i < 20; i++ {
			userID := l.genUserID()
			_, err := l.svcCtx.UserDB.FindWithError(l.ctx, []string{userID})
			if err != nil {
				// 用户不存在，可以使用
				req.UserID = userID
				break
			}
		}
		if req.UserID == "" {
			return nil, fmt.Errorf("gen user id failed")
		}
	} else {
		// 检?UserID 是否已被使用
		_, err := l.svcCtx.UserDB.FindWithError(l.ctx, []string{req.UserID})
		if err == nil {
			return nil, fmt.Errorf("userID is used")
		}
	}

	// 创建通知账户用户
	now := time.Now()
	users := []*model.User{
		{
			UserID:         req.UserID,
			Nickname:       req.NickName,
			FaceURL:        req.FaceURL,
			CreateTime:     now,
			AppMangerLevel: req.AppMangerLevel,
		},
	}

	if err := l.svcCtx.UserDB.Create(l.ctx, users); err != nil {
		return nil, err
	}

	return &user.AddNotificationAccountResp{
		UserID:         req.UserID,
		NickName:       req.NickName,
		FaceURL:        req.FaceURL,
		AppMangerLevel: req.AppMangerLevel,
	}, nil
}

func (l *AddNotificationAccountLogic) genUserID() string {
	const length = 10
	data := make([]byte, length)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		if i == 0 {
			data[i] = chars[1:][data[i]%9]
		} else {
			data[i] = chars[data[i]%10]
		}
	}
	return string(data)
}
