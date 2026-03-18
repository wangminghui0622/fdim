package authverify

import (
	"context"
	"fmt"

	"fdim/pkg/constant"
	"fdim/pkg/mcontext"
	pct "fdim/protocol/constant"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

// CheckAdmin 检查是否是管理�?
func CheckAdmin(ctx context.Context) error {
	if IsAdmin(ctx) {
		return nil
	}
	opUserID := mcontext.GetOpUserID(ctx)
	return fmt.Errorf("user %s is not admin userID", opUserID)
}

// IsAdmin 检查是否是管理员（临时管理员或系统管理员）
func IsAdmin(ctx context.Context) bool {
	return IsTempAdmin(ctx) || IsSystemAdmin(ctx)
}

// IsTempAdmin 检查是否是临时管理�?
func IsTempAdmin(ctx context.Context) bool {
	flag, _ := ctx.Value(constant.CtxTempAdminKey).(bool)
	return flag
}

// WithTempAdmin 在 context 中标记临时管理员
func WithTempAdmin(ctx context.Context) context.Context {
	return context.WithValue(ctx, constant.CtxTempAdminKey, true)
}

// IsSystemAdmin 检查是否是系统管理�?
func IsSystemAdmin(ctx context.Context) bool {
	opUserID := mcontext.GetOpUserID(ctx)
	adminUserIDs := GetIMAdminUserIDs(ctx)
	for _, adminID := range adminUserIDs {
		if opUserID == adminID {
			return true
		}
	}
	return false
}

// CheckAccess 检查访问权限（本人或管理员�?
func CheckAccess(ctx context.Context, ownerUserID string) error {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == ownerUserID {
		return nil
	}
	if IsAdmin(ctx) {
		return nil
	}
	return fmt.Errorf("no permission: ownerUserID=%s, opUserID=%s", ownerUserID, opUserID)
}

// CheckAccessIn 检查是否在指定的用户ID列表中或管理�?
func CheckAccessIn(ctx context.Context, ownerUserIDs ...string) error {
	opUserID := mcontext.GetOpUserID(ctx)
	for _, userID := range ownerUserIDs {
		if opUserID == userID {
			return nil
		}
	}
	if IsAdmin(ctx) {
		return nil
	}
	return fmt.Errorf("no permission: opUserID not in ownerUserIDs")
}

// GetIMAdminUserIDs �?context 中获取管理员用户ID列表
func GetIMAdminUserIDs(ctx context.Context) []string {
	adminUserIDs, _ := ctx.Value(constant.CtxAdminUserIDsKey).([]string)
	if adminUserIDs == nil {
		return []string{}
	}
	return adminUserIDs
}

// WithIMAdminUserIDs 将管理员用户ID列表添加�?context
func WithIMAdminUserIDs(ctx context.Context, imAdminUserID []string) context.Context {
	return context.WithValue(ctx, constant.CtxAdminUserIDsKey, imAdminUserID)
}

// CheckUserIsAdmin 检查指定用户是否是管理�?
func CheckUserIsAdmin(ctx context.Context, userID string) bool {
	adminUserIDs := GetIMAdminUserIDs(ctx)
	for _, adminID := range adminUserIDs {
		if userID == adminID {
			return true
		}
	}
	return false
}

// CheckSystemAccount 检查是否是系统账号
func CheckSystemAccount(ctx context.Context, level int32) bool {
	return level >= pct.AppAdmin
}

// LogAccessDenied 记录访问被拒绝的日志
func LogAccessDenied(ctx context.Context, reason string) {
	opUserID := mcontext.GetOpUserID(ctx)
	logx.WithContext(ctx).Errorf("Access denied: userID=%s, reason=%s", opUserID, reason)
}

// Secret returns a jwt.Keyfunc that returns the secret key for token validation
func Secret(secret string) func(*jwt.Token) (any, error) {
	return func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}
}
