package handler

import (
	"errors"
	"net/http"

	"fdim/pkg/errs"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 统一响应格式（与官方 open-im-server 一致）
type Response struct {
	ErrCode int         `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	ErrDlt  string      `json:"errDlt"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(w http.ResponseWriter, data interface{}) {
	httpx.WriteJson(w, http.StatusOK, Response{
		ErrCode: 0,
		ErrMsg:  "",
		ErrDlt:  "",
		Data:    data,
	})
}

// ApiError 从 error 提取 CodeError 信息，返回与官方一致的错误响应
func ApiError(w http.ResponseWriter, err error) {
	var codeErr errs.CodeError
	if errors.As(err, &codeErr) {
		httpx.WriteJson(w, http.StatusOK, Response{
			ErrCode: codeErr.Code(),
			ErrMsg:  codeErr.Msg(),
			ErrDlt:  codeErr.Detail(),
			Data:    nil,
		})
		return
	}
	httpx.WriteJson(w, http.StatusOK, Response{
		ErrCode: errs.ServerInternalError,
		ErrMsg:  err.Error(),
		ErrDlt:  "",
		Data:    nil,
	})
}

// ParamError 参数错误
func ParamError(w http.ResponseWriter, err error) {
	httpx.WriteJson(w, http.StatusOK, Response{
		ErrCode: errs.ArgsError,
		ErrMsg:  "ArgsError",
		ErrDlt:  err.Error(),
		Data:    nil,
	})
}

// ServerError 服务器错误（兼容旧调用，内部提取 CodeError）
func ServerError(w http.ResponseWriter, err error) {
	ApiError(w, err)
}
