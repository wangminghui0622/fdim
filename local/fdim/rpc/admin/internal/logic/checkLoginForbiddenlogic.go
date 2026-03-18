package logic

import (
	"context"
	"fmt"

	"fdim/pkg/errs"
	"fdim/protocol/admin"
	"fdim/rpc/admin/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/mongo"
)

type CheckLoginForbiddenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckLoginForbiddenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckLoginForbiddenLogic {
	return &CheckLoginForbiddenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckLoginForbiddenLogic) CheckLoginForbidden(req *admin.CheckLoginForbiddenReq) (*admin.CheckLoginForbiddenResp, error) {
	// 1. 验证参数
	if req.Ip == "" {
		return nil, errs.ErrArgs.WrapMsg("ip cannot be empty")
	}
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID cannot be empty")
	}

	// 2. 检查IP是否被禁止登录
	forbiddens, err := l.svcCtx.AdminDB.FindIPForbidden(l.ctx, []string{req.Ip})
	if err != nil {
		l.Errorf("FindIPForbidden failed: %v", err)
		return nil, errs.WrapMsg(err, "failed to check IP forbidden")
	}
	for _, forbidden := range forbiddens {
		if forbidden.LimitLogin {
			return nil, fmt.Errorf("ip %s is forbidden for login", req.Ip)
		}
	}

	// 3. 检查用户IP登录限制
	_, err = l.svcCtx.AdminDB.GetLimitUserLoginIP(l.ctx, req.UserID, req.Ip)
	if err != nil {
		// 如果未找到此IP的限制记录，检查用户是否有其他IP限制
		if err == mongo.ErrNoDocuments || fmt.Sprintf("%v", err) == "limit not found" {
			// 检查用户是否有任何IP限制
			count, err := l.svcCtx.AdminDB.CountLimitUserLoginIP(l.ctx, req.UserID)
			if err != nil {
				l.Errorf("CountLimitUserLoginIP failed: %v", err)
				return nil, errs.WrapMsg(err, "failed to count user IP limits")
			}
			if count > 0 {
				// 用户有IP限制，但当前IP不在允许列表中
				return nil, fmt.Errorf("user %s is restricted to specific IPs, current IP %s is not allowed", req.UserID, req.Ip)
			}
		} else {
			l.Errorf("GetLimitUserLoginIP failed: %v", err)
			return nil, errs.WrapMsg(err, "failed to check user IP limit")
		}
	}

	// 4. 检查用户是否被封禁
	blockInfo, err := l.svcCtx.AdminDB.GetBlockInfo(l.ctx, req.UserID)
	if err == nil && blockInfo != nil {
		return nil, fmt.Errorf("user %s is blocked: %s", req.UserID, blockInfo.Reason)
	}

	return &admin.CheckLoginForbiddenResp{}, nil
}
