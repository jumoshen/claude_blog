package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

func ErrorWithStatus(c *gin.Context, status int, code int, message string) {
	c.JSON(status, Response{
		Code:    code,
		Message: message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusUnauthorized, 401, message)
}

func Forbidden(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusForbidden, 403, message)
}

func NotFound(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusNotFound, 404, message)
}

func InternalError(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusInternalServerError, 500, message)
}

func BadRequest(c *gin.Context, message string) {
	ErrorWithStatus(c, http.StatusBadRequest, 400, message)
}
