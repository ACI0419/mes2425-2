package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 响应状态码常量
const (
	SuccessCode       = 200
	ErrorCode         = 400
	UnauthorizedCode  = 401
	ForbiddenCode     = 403
	NotFoundCode      = 404
	InternalErrorCode = 500
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 请求错误
func BadRequest(c *gin.Context, message string) {
	Error(c, ErrorCode, message)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	Error(c, UnauthorizedCode, message)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	Error(c, ForbiddenCode, message)
}

// NotFound 未找到
func NotFound(c *gin.Context, message string) {
	Error(c, NotFoundCode, message)
}

// InternalError 内部错误
func InternalError(c *gin.Context, message string) {
	Error(c, InternalErrorCode, message)
}
