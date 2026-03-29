package types

// Auth 相关请求和响应类?

// GetAdminToken 相关类型
type GetAdminTokenReq struct {
	Secret string `json:"secret"`
	UserID string `json:"userID"`
}

type GetAdminTokenResp struct {
	BaseResp
	Token             string `json:"token"`
	ExpireTimeSeconds int64  `json:"expireTimeSeconds"`
}

// GetUserToken 相关类型
type GetUserTokenReq struct {
	PlatformID int32  `json:"platformID"`
	UserID     string `json:"userID"`
	Password   string `json:"password"` // 添加密码字段
}

type GetUserTokenResp struct {
	BaseResp
	Token             string `json:"token"`
	ExpireTimeSeconds int64  `json:"expireTimeSeconds"`
}

// ParseToken 相关类型
type ParseTokenReq struct {
	Token string `json:"token"`
}

type ParseTokenResp struct {
	BaseResp
	UserID            string `json:"userID"`
	PlatformID        int32  `json:"platformID"`
	ExpireTimeSeconds int64  `json:"expireTimeSeconds"`
}

// ForceLogout 相关类型
type ForceLogoutReq struct {
	PlatformID int32  `json:"platformID"`
	UserID     string `json:"userID"`
}

type ForceLogoutResp struct {
	BaseResp
}
