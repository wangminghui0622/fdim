package types

// Third 相关请求和响应类型
// 注意：这些类型需要与 open-im-server 的 API 接口保持一致

// FcmUpdateToken 相关类型
type FcmUpdateTokenReq struct {
	PlatformID int32  `json:"platformID"`
	FcmToken   string `json:"fcmToken"`
	Account    string `json:"account,optional"`
	ExpireTime int64  `json:"expireTime,optional"`
}

type FcmUpdateTokenResp struct {
	BaseResp
}

// SetAppBadge 相关类型
type SetAppBadgeReq struct {
	UserID string `json:"userID"`
	AppUnreadCount int32 `json:"appUnreadCount"`
}

type SetAppBadgeResp struct {
	BaseResp
}

// UploadLogs 相关类型
type UploadLogsReq struct {
	Logs []interface{} `json:"logs"` // 使用 LogInfo
}

type UploadLogsResp struct {
	BaseResp
}

// DeleteLogs 相关类型
type DeleteLogsReq struct {
	LogIDs []string `json:"logIDs"`
}

type DeleteLogsResp struct {
	BaseResp
}

// SearchLogs 相关类型
type SearchLogsReq struct {
	Keyword    string     `json:"keyword,optional"`
	StartTime  int64      `json:"startTime,optional"`
	EndTime     int64      `json:"endTime,optional"`
	Pagination Pagination `json:"pagination"`
}

type SearchLogsResp struct {
	BaseResp
	Logs  []interface{} `json:"logs"` // 使用 LogInfo
	Total int64          `json:"total"`
}

// PartLimit 相关类型
type PartLimitReq struct {
}

type PartLimitResp struct {
	BaseResp
	MinPartSize int64 `json:"minPartSize"`
	MaxPartSize int64 `json:"maxPartSize"`
	MaxNumSize  int32 `json:"maxNumSize"`
}

// PartSize 相关类型
type PartSizeReq struct {
	Size int64 `json:"size"`
}

type PartSizeResp struct {
	BaseResp
	Size int64 `json:"size"`
}

// InitiateMultipartUpload 相关类型
type InitiateMultipartUploadReq struct {
	Hash        string `json:"hash"`
	Size        int64  `json:"size"`
	PartSize    int64  `json:"partSize"`
	MaxParts    int32  `json:"maxParts"`
	Cause       string `json:"cause,optional"`
	UrlPrefix   string `json:"urlPrefix,optional"`
}

type InitiateMultipartUploadResp struct {
	BaseResp
	Url    string      `json:"url"`
	Upload interface{} `json:"upload"` // 使用 UploadInfo
}

// AuthSign 相关类型
type AuthSignReq struct {
	UploadID   string  `json:"uploadID"`
	PartNumbers []int32 `json:"partNumbers"`
}

type AuthSignResp struct {
	BaseResp
	Url    string            `json:"url"`
	Query  map[string]string `json:"query"`
	Header map[string]string `json:"header"`
}

// CompleteMultipartUpload 相关类型
type CompleteMultipartUploadReq struct {
	UploadID    string   `json:"uploadID"`
	Parts       []string `json:"parts"`
	Name        string   `json:"name"`
	ContentType string   `json:"contentType,optional"`
	Cause       string   `json:"cause,optional"`
	UrlPrefix   string   `json:"urlPrefix,optional"`
}

type CompleteMultipartUploadResp struct {
	BaseResp
	Url string `json:"url"`
}

// AccessURL 相关类型
type AccessURLReq struct {
	Name  string            `json:"name"`
	Query map[string]string `json:"query,optional"`
}

type AccessURLResp struct {
	BaseResp
	Url string `json:"url"`
}

// InitiateFormData 相关类型
type InitiateFormDataReq struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType,optional"`
	Group       string `json:"group,optional"`
	Millisecond int64  `json:"millisecond,optional"`
	Absolute    bool   `json:"absolute,optional"`
}

type InitiateFormDataResp struct {
	BaseResp
	Id          string            `json:"id"`
	Url         string            `json:"url"`
	File        string            `json:"file"`
	Header      []interface{}     `json:"header,optional"` // 使用 KeyValues
	FormData    map[string]string `json:"formData"`
	Expires     int64             `json:"expires,optional"`
	SuccessCodes []int32          `json:"successCodes,optional"`
}

// CompleteFormData 相关类型
type CompleteFormDataReq struct {
	Id        string `json:"id"`
	UrlPrefix string `json:"urlPrefix,optional"`
}

type CompleteFormDataResp struct {
	BaseResp
	Url string `json:"url"`
}

// GetPrometheus 相关类型（GET 请求，无请求体）
type GetPrometheusResp struct {
	BaseResp
	Url string `json:"url"`
}
