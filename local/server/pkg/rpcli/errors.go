package rpcli

import "errors"

var (
	// ErrUserNotFound 用户未找�?
	ErrUserNotFound = errors.New("user not found")
	// ErrGroupNotFound 群组未找�?
	ErrGroupNotFound = errors.New("group not found")
)
