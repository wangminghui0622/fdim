package com

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"fdim/pkg/constant"
)

// GetOpUserID 从context 中获取操作用户ID
func GetOpUserID(ctx context.Context) string {
	userID, _ := ctx.Value(constant.CtxOpUserIDKey).(string)
	return userID
}

// WithOpUserID 将操作用户ID添加到context
func WithOpUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, constant.CtxOpUserIDKey, userID)
}

// GetOperationID 从context 中获取操作ID
func GetOperationID(ctx context.Context) string {
	operationID, _ := ctx.Value(constant.CtxOperationIDKey).(string)
	if operationID == "" {
		// 如果没有操作ID，生成一?
		operationID = generateOperationID()
	}
	return operationID
}

// WithOperationID 将操作ID添加到context
func WithOperationID(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, constant.CtxOperationIDKey, operationID)
}

// generateOperationID 生成操作ID
func generateOperationID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int63())
}

// GenGroupID 生成群组ID（使用MD5 哈希?
func GenGroupID(ctx context.Context, checkExists func(string) (bool, error)) (string, error) {
	operationID := GetOperationID(ctx)
	for i := 0; i < 10; i++ {
		// 使用 MD5 哈希生成ID
		data := strings.Join([]string{
			operationID,
			strconv.FormatInt(time.Now().UnixNano(), 10),
			strconv.Itoa(rand.Int()),
		}, ",;,")

		hash := md5.Sum([]byte(data))
		hashStr := fmt.Sprintf("%x", hash)

		// 取前8个字符转换为大整数，然后转为字符??
		bi := big.NewInt(0)
		bi.SetString(hashStr[0:8], 16)
		groupID := bi.String()

		// 检查ID是否已存??
		exists, err := checkExists(groupID)
		if err != nil {
			// 如果检查出错，继续尝试下一个ID
			continue
		}
		if !exists {
			return groupID, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique group ID after 10 attempts")
}
