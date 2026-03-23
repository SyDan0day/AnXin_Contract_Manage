package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 统一的错误响应结构
type ErrorResponse struct {
	Error     string `json:"error"`
	Code      string `json:"code,omitempty"`
	Details   string `json:"details,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// ErrorHandlerMiddleware 全局错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误堆栈信息
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

				// 返回统一的错误响应
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Internal server error",
					Code:    "INTERNAL_ERROR",
					Details: fmt.Sprintf("%v", err),
				})
				c.Abort()
			}
		}()
		c.Next()

		// 处理路由未找到的情况
		if c.Writer.Status() == http.StatusNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Resource not found",
				Code:  "NOT_FOUND",
			})
			c.Abort()
		}
	}
}

// AppError 应用错误类型
type AppError struct {
	StatusCode int
	Code       string
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError 创建应用错误
func NewAppError(statusCode int, code, message string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// HandleError 处理错误并返回统一的JSON响应
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		c.JSON(appErr.StatusCode, ErrorResponse{
			Error: appErr.Message,
			Code:  appErr.Code,
		})
	} else {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}
}
