package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseCode string

const (
	CodeSuccess      ResponseCode = "SUCCESS"
	CodeBadRequest   ResponseCode = "BAD_REQUEST"
	CodeUnauthorized ResponseCode = "UNAUTHORIZED"
	CodeForbidden    ResponseCode = "FORBIDDEN"
	CodeNotFound     ResponseCode = "NOT_FOUND"
	CodeServerError  ResponseCode = "SERVER_ERROR"
)

type APIResponse struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
}

type LegacyResponse struct {
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code: CodeSuccess,
		Data: data,
	})
}

func SuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Code:    CodeBadRequest,
		Message: message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Code:    CodeUnauthorized,
		Message: message,
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, APIResponse{
		Code:    CodeForbidden,
		Message: message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Code:    CodeNotFound,
		Message: message,
	})
}

func ServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Code:    CodeServerError,
		Message: message,
	})
}
